package api

import (
	"database/sql"
	"fmt"
	"github.com/go-chi/render"
	"net/http"
	"strconv"
	"strings"
	"yss-go-official/logger"
	conn "yss-go-official/orm"
)

var keyFeeHistory = "%s:feeHistory"
// 只通过redis做缓存

// fee_history
func handleFeeHistory(w http.ResponseWriter, r *http.Request) {
	var err error
	if account, ok := checkBinding(r); ok {
		_key := fmt.Sprintf(keyFeeHistory, account)
		if resp, exist := redisGet(r, _key); exist {
			err = render.Render(w, r, NewResponseOK(resp))
			if err != nil {
				logger.GetLogEntry(r).Info(err)
			}
		} else {
			resp := handleFeeHistoryInner(r, account)
			go redisSet(r, _key, resp)
			err = render.Render(w, r, NewResponseOK(resp))
			if err != nil {
				logger.GetLogEntry(r).Info(err)
			}
		}
	} else {
		err = render.Render(w, r, NewErrInvalidBinding())
		if err != nil {
			logger.GetLogEntry(r).Info(err)
		}
	}
}

func handleFeeHistoryInner(r *http.Request, account string) *feeHistoryResponse {
	var as accountStation
	var ph []paidHistory // 存放已缴费的数据

	rows, err := conn.MssqlDB.Query(billList,
		sql.Named("account", account),
	)
	if err != nil {
		logger.GetLogEntry(r).Info("cant get bill list from mssql, err: ", err)
	} else {
		/** 去mssql拿东西 */
		var (
			uDate            string
			tmpCharge        string
			tmpCurrentMeter  sql.NullString
			tmpPreviousMeter sql.NullString
			tmpIsPaid        string
			uYszbh           sql.NullString
		)
		var (
			// account 这个account是可信的，必定有的
			_address sql.NullString
			_name    sql.NullString
			_phone   sql.NullString
		)
		// 两个驱动的使用方法不一致，mssql的不支持stmt
		err = conn.MssqlDB.QueryRow(mssqlQueryAccountCmd, sql.Named("account", account)).Scan(&as.Account, &_address, &_name, &_phone)
		if err != nil {
			logger.GetLogEntry(r).Info("cant get account from mssql, err: ", err)
		} else {
			as.Address = _address.String
			as.Name = _name.String
			as.Phone = _phone.String

			as.Paid = true

			unpaidCount := 0
			unpaidCharge := 0.0
			var unpaidMeter int64 = 0

			for rows.Next() {
				err := rows.Scan(&uDate, &tmpCharge, &tmpCurrentMeter, &tmpPreviousMeter, &tmpIsPaid, &uYszbh)
				if err != nil {
					logger.GetLogEntry(r).Info("inner error in mssql, err: ", err)
					break
				}

				// 只对有应收账编号的记录进行处理
				if uYszbh.Valid {
					// 初始化码数
					if len(as.Meter) == 0 {
						if tmpCurrentMeter.Valid {
							as.CurrentMeter = tmpCurrentMeter.String // 记录当前表码
							if cr, err := strconv.ParseInt(tmpCurrentMeter.String, 10, 32); err == nil {
								if pr, err := strconv.ParseInt(tmpPreviousMeter.String, 10, 32); err == nil {
									as.Meter = strconv.FormatInt(cr-pr, 10) // 记录用水量
								}
							}
						}
						if len(as.Meter) == 0 {
							as.Meter = "0"
						}
					}

					// 处理未缴费的
					if tmpIsPaid != "true" {
						if c, err := strconv.ParseFloat(tmpCharge, 10); err == nil {
							unpaidCharge += c
						}
						as.Paid = false
						unpaidCount += 1
						if tmpCurrentMeter.Valid {
							if cr, err := strconv.ParseInt(tmpCurrentMeter.String, 10, 32); err == nil {
								if pr, err := strconv.ParseInt(tmpPreviousMeter.String, 10, 32); err == nil {
									unpaidMeter += cr - pr
								}
							}
						}
					} else {
						ph = append(ph, paidHistory{
							Date:   uDate,
							Charge: tmpCharge,
							BillId: uYszbh.String,
						})
					}
				}
			}

			// 根据是否欠费来更新字段
			if as.Paid {
				as.Charge = "0"
				as.UnpaidPeriodCount = 0
			} else {
				as.Charge = strconv.FormatFloat(unpaidCharge, 'f', -1, 32)
				as.UnpaidPeriodCount = unpaidCount
				as.Meter = strconv.FormatInt(unpaidMeter, 10)
			}

			as.Account = strings.Trim(as.Account, " ")
		}
	}

	return newFeeHistoryResponse(&as, ph)
}
