package api

import (
	"database/sql"
	"github.com/go-chi/render"
	"github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"yss-go-official/logger"
	conn "yss-go-official/orm"
)


// accounts
/** handleAccounts -> handleSingleAccount -> handleSingleAccountSync -> getAccountDataFromMssql -> getBillListFromMssql */
func handleAccounts(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(string)

	// 获取已绑定的户号
	var accounts []string
	stmt, err := conn.MysqlDB.Prepare(mysqlQueryAccountCmd)
	if err != nil {
		logger.GetLogEntry(r).Infof("Mysql prepare statement error, stmt: %s, err: %s ", mysqlQueryAccountCmd, err)
	} else {
		if rows, err := stmt.Query(id); err != nil {
			logger.GetLogEntry(r).Infof("Mysql statement exec error, params: %s, err: %s", id, err)
		} else {
			for rows.Next() {
				// 先生产存放的struct
				var account string
				// 把结果逐个赋值
				err := rows.Scan(&account)
				if err != nil {
					logger.GetLogEntry(r).Info("Mysql rows scan error, err: ", err)
					break
				} else {
					accounts = append(accounts, account)
				}
			}
			if err = rows.Err(); err != nil {
				logger.GetLogEntry(r).Info("Mysql rows error, err: ", err)
			}
		}
	}

	var responseContent interface{}

	if len(accounts) > 0 {
		// 获取默认账号
		var defaultAccount sql.NullString
		stmt, err := conn.MysqlDB.Prepare(mysqlQueryDefaultAccountCmd)
		if err != nil {
			logger.GetLogEntry(r).Infof("Mysql prepare statement error, stmt: %s, err: %s ", mysqlQueryDefaultAccountCmd, err)
		} else {
			err = stmt.QueryRow(id).Scan(&defaultAccount)
			if err != nil {
				logger.GetLogEntry(r).Infof("Mysql statement exec error, params: %s, err: %s", id, err)
			} else {
				// handle single
				outputCh := make(chan *accountStation, 10)

				var res [] *accountsResponse
				var wg sync.WaitGroup

				// 分发任务
				for _, account := range accounts {
					wg.Add(1)
					go handleSingleAccount(&wg, account, defaultAccount.String, outputCh, r)
				}
				wg.Wait()
				outputCh <- nil // 关闭通道，不然下面会一直阻塞

				// 收集结果
				for as := range outputCh {
					if as == nil {
						break
					}
					if len(as.Account) > 0 {
						res = append(res, newAccountsResponse(as))
					}
				}
				close(outputCh)

				responseContent = res
			}
		}
	}

	err = render.Render(w, r, NewResponseOK(responseContent))
	if err != nil {
		logger.GetLogEntry(r).Info(err)
	}
}

// 处理单个account的主逻辑 go
func handleSingleAccount(wg *sync.WaitGroup, account string, defaultAccount string, outputCh chan *accountStation, r *http.Request) {
	as := handleSingleAccountSync(account, defaultAccount, r)
	outputCh <- &as
	wg.Done()
}

// 处理单个account的主逻辑 同步
func handleSingleAccountSync(account string, defaultAccount string, r *http.Request) accountStation {
	var as accountStation
	stmt, err := conn.MysqlDB.Prepare(mysqlQueryAccountDataCmd)
	if err != nil {
		logger.GetLogEntry(r).Infof("Mysql prepare statement error, stmt: %s, err: %s ", mysqlQueryAccountCmd, err)
	} else {
		var _id sql.NullInt64

		err = stmt.QueryRow(account).Scan(&_id, &as.Account, &as.Address, &as.Name, &as.Phone, &as.Charge, &as.CurrentMeter, &as.Meter, &as.Paid, &as.UnpaidPeriodCount)
		if err != nil {
			// 遇到缺失的，在这里更新
			as = getAccountDataFromMssql(account, r)
			// 另起go去插入新数据，主逻辑先响应请求
			go insertToMysql(as, r)
		}

		as.Account = strings.Trim(as.Account, " ")

		// 是否为默认账户
		if defaultAccount == account {
			as.Default_ = true
		} else {
			as.Default_ = false
		}
	}

	return as
}

// 从用户档案表获取meta数据
func getAccountDataFromMssql(account string, r *http.Request) accountStation {
	var as accountStation

	var (
		// account 这个account是可信的，必定有的
		_address sql.NullString
		_name    sql.NullString
		_phone   sql.NullString
	)
	// 两个驱动的使用方法不一致，mssql的不支持stmt
	err := conn.MssqlDB.QueryRow(mssqlQueryAccountCmd, sql.Named("account", account)).Scan(&as.Account, &_address, &_name, &_phone)
	if err != nil {
		logger.GetLogEntry(r).Infof("Cant get account data from mssql,  err: %s ", err)
	} else {
		as.Address = _address.String
		as.Name = _name.String
		as.Phone = _phone.String

		as.Paid = true
		getBillListFromMssql(&as, r)
	}
	return as
}

// 调用存储过程Mob_GetHisBill
func getBillListFromMssql(as *accountStation, r *http.Request) {
	rows, err := conn.MssqlDB.Query(billList,
		sql.Named("account", as.Account),
	)
	if err != nil {
		logger.GetLogEntry(r).Info(err)
		return
		// log.Fatal("in billList ", err) // 在出现错误的时候，就直接跳出了
	} else {
		var (
			uData            string
			tmpCharge        string
			tmpCurrentMeter  sql.NullString
			tmpPreviousMeter sql.NullString
			tmpIsPaid        string
			uYszbh           sql.NullString
		)
		unpaidCount := 0
		unpaidCharge := 0.0
		var unpaidMeter int64 = 0

		for rows.Next() {
			err := rows.Scan(&uData, &tmpCharge, &tmpCurrentMeter, &tmpPreviousMeter, &tmpIsPaid, &uYszbh)
			if err != nil {
				log.Fatal("inner rows", err)
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

		err = rows.Err()
		if err != nil {
			logger.GetLogEntry(r).Info(err)
		}
	}
}

// 写入新的account data到mysql
func insertToMysql(as accountStation, r *http.Request) {
	stmt, err := conn.MysqlDB.Prepare(mysqlInsertAccountCmd)

	if err != nil {
		logger.GetLogEntry(r).Infof("Mysql prepare statement error, stmt: %s, err: %s ", mysqlInsertAccountCmd, err)
	} else {
		// 插入失败的话就更新
		_, err = stmt.Exec(as.Account, as.Address, as.Name, as.Phone, as.Charge, as.CurrentMeter, as.Meter, as.Paid, as.UnpaidPeriodCount)
		if err != nil {
			if driverErr, ok := err.(*mysql.MySQLError); ok {
				if driverErr.Number == 1062 {
					stmt, err = conn.MysqlDB.Prepare(mysqlUpdateAccountCmd)
					if err != nil {
						logger.GetLogEntry(r).Infof("Mysql prepare statement error, stmt: %s, err: %s ", mysqlUpdateAccountCmd, err)
					} else {
						_, err = stmt.Exec(as.Phone, as.Charge, as.CurrentMeter, as.Meter, as.Paid, as.UnpaidPeriodCount, as.Account)
						if err != nil {
							logger.GetLogEntry(r).Infof("Mysql update account data error, err: %s ", err)
						}
					}
				} else {
					logger.GetLogEntry(r).Infof("Mysql insert account data error, err: %s ", err)
				}
			}
		}
	}
}
