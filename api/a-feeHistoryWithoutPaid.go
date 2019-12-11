package api

import (
	"database/sql"
	"fmt"
	"github.com/go-chi/render"
	"net/http"
	"strings"
	"sync"
	"yss-go-official/logger"
	conn "yss-go-official/orm"
)

var keyFeeHistoryWithoutPaid = "%s:feeHistoryWithoutPaid"

// fee_history_without_paid
func handleFeeHistoryWithoutPaid(w http.ResponseWriter, r *http.Request) {
	var err error
	if account, ok := checkBinding(r); ok {
		_key := fmt.Sprintf(keyFeeHistoryWithoutPaid, account)
		if resp, exist := redisGet(r, _key); exist {
			err = render.Render(w, r, NewResponseOK(resp))
			if err != nil {
				logger.GetLogEntry(r).Info(err)
			}
		} else {
			resp := handleFeeHistoryWithoutPaidInner(r, account)
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

func handleFeeHistoryWithoutPaidInner(r *http.Request, account string) *feeHistoryWithoutPaidResponse {
	var as accountStation
	var uph []unpaidHistory // 存放未缴费的数据

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
			as.Name = _name.String
			var unpaidPeriod [][]string

			for rows.Next() {
				err := rows.Scan(&uDate, &tmpCharge, &tmpCurrentMeter, &tmpPreviousMeter, &tmpIsPaid, &uYszbh)
				if err != nil {
					logger.GetLogEntry(r).Info("inner error in mssql, err: ", err)
					break
				}

				// 只对有应收账编号的记录进行处理
				if uYszbh.Valid {
					// 处理未缴费的
					if tmpIsPaid != "true" {
						unpaidPeriod = append(unpaidPeriod, []string{uYszbh.String, uDate})
					}
				}
			}

			outputCh := make(chan *unpaidHistory, 10)
			var wg sync.WaitGroup

			for _, period := range unpaidPeriod {
				go handleSingleUnpaid(r, &wg, outputCh, period)
				wg.Add(1)
			}

			wg.Wait()
			outputCh <- nil // 关闭通道，不然下面会一直阻塞

			// 收集结果
			for item := range outputCh {
				if item == nil {
					break
				}
				if len(item.Charge) > 0 {
					uph = append(uph, *item)
				}
			}
			close(outputCh)

			as.Account = strings.Trim(as.Account, " ")
		}
	}

	return newFeeHistoryWithoutPaidResponse(&as, uph)
}

func handleSingleUnpaid(r *http.Request, wg *sync.WaitGroup, outputCh chan *unpaidHistory, unpaidMeta []string) {
	res := newUnpaidHistory(getFeeDetail(r, unpaidMeta[0], unpaidMeta[1]), unpaidMeta[1])
	outputCh <- res
	wg.Done()
}
