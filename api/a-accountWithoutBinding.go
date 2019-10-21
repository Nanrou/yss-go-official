package api

import (
	"database/sql"
	"github.com/go-chi/render"
	"github.com/go-sql-driver/mysql"
	"net/http"
	"yss-go-official/logger"
	conn "yss-go-official/orm"
)


// account_without_binding
func handleAccountWithoutBinding(w http.ResponseWriter, r *http.Request) {
	account := r.URL.Query().Get("account")
	name := r.URL.Query().Get("name")
	if len(account) ==0 || len(name) == 0 {
		err := render.Render(w, r, NewErrRequiredFields())
		if err != nil {
			logger.GetLogEntry(r).Info(err)
		}
		return
	}

	var _useless string
	err := conn.MssqlDB.QueryRow(mssqlQueryCheckAccountCmd,
		sql.Named("account", account),
		sql.Named("name", name)).Scan(&_useless)
	if err != nil {
		logger.GetLogEntry(r).Infof("Mssql check account error, err: %s", err)
		// 查不到则返回错误的预留信息
		err = render.Render(w, r, NewErrInvalidAccountData())
		if err != nil {
			logger.GetLogEntry(r).Info(err)
		}
	} else {
		var awbr accountWithoutBindingResponse
		// 先去mysql查询是否有这条记录
		stmt, err := conn.MysqlDB.Prepare(mysqlQueryFeeDetailCmd)
		if err != nil {
			logger.GetLogEntry(r).Infof("Mysql prepare statement error, stmt: %s, err: %s ", mysqlQueryAccountCmd, err)
		} else {
			var of otherFee
			err = stmt.QueryRow(account).Scan(&_useless, &_useless, &awbr.Account, &awbr.Name, &awbr.CurrentPeriod, &awbr.Charge,
				&awbr.CurrentMeter, &awbr.PreviousMeter, &awbr.Paid, &of.Wsf, &of.Xfft, &of.Ljf, &of.Ecjydf, &of.Szyf, &of.Cjhys, &of.Wyj, &of.Wswyj, &awbr.WaterCharge)
			if err != nil {
				// 遇到缺失的，在这里更新
				if err == sql.ErrNoRows {
					awbr = handleSingleAccountWithoutBinding(r, account)
				} else {
					logger.GetLogEntry(r).Infof("Mysql query detail error, err: %s ", err)
				}
			} else {
				awbr.OtherFee = of
			}
		}
		err = render.Render(w, r, NewResponseOK(awbr))
		if err != nil {
			logger.GetLogEntry(r).Info(err)
		}
	}
}

// 这里是处理缓存库中没有的，需要采用最新的
func handleSingleAccountWithoutBinding(r *http.Request, account string) accountWithoutBindingResponse {
	var awbr accountWithoutBindingResponse

	rows, err := conn.MssqlDB.Query(billList,
		sql.Named("account", account),
	)
	if err != nil {
		logger.GetLogEntry(r).Info("cant get bill list from mssql, err: ", err)
	} else {
		var (
			uDate            string
			tmpCharge        string
			tmpCurrentMeter  sql.NullString
			tmpPreviousMeter sql.NullString
			tmpIsPaid        string
			uYszbh           sql.NullString
		)

		for rows.Next() {
			err := rows.Scan(&uDate, &tmpCharge, &tmpCurrentMeter, &tmpPreviousMeter, &tmpIsPaid, &uYszbh)
			if err != nil {
				logger.GetLogEntry(r).Info("inner error in mssql, err: ", err)
				continue
			}

			// 只对有应收账编号的记录进行处理
			if uYszbh.Valid {
				awbr = *newAccountWithoutBindingResponse(getFeeDetail(r, uYszbh.String, uDate), account, uDate)
				// 另起go去写入缓存
				go insertFeeDetailToMysql(r, uYszbh.String, &awbr)
				break
			}
		}
	}

	return awbr
}

// 获取fee detail
func getFeeDetail(r *http.Request, yszbh string, date string) *feeDetail {
	var fd feeDetail
	rows, err := conn.MssqlDB.Query(billDetail,
		sql.Named("yszbh", yszbh),
		sql.Named("date", date),
	)
	if err != nil {
		logger.GetLogEntry(r).Info("cant get bill detail from mssql, err: ", err)
	} else {
		var tmpPaid string
		for rows.Next() {
			err = rows.Scan(&fd.Address, &fd.Name, &fd.Charge, &fd.CurrentMeter, &fd.MeterReadingDate,
				&tmpPaid, &fd.PreviousMeter, &fd.WaterCharge, &fd.WaterProperty, &fd.Wsf, &fd.Xfft,
				&fd.Ljf, &fd.Ecjydf, &fd.Szyf, &fd.Cjhys, &fd.Wyj, &fd.Wswyj)
			if err != nil {
				logger.GetLogEntry(r).Info("inner error in mssql, err: ", err)
				// 只会有一行
			}
		}
		if tmpPaid == "true" {
			fd.Paid = true
		} else {
			fd.Paid = false
		}
	}
	return &fd
}

// 将fee detail 写入mysql
func insertFeeDetailToMysql(r *http.Request, yszbh string, awbr *accountWithoutBindingResponse) {
	stmt, err := conn.MysqlDB.Prepare(mysqlInsertFeeDetailCmd)

	if err != nil {
		logger.GetLogEntry(r).Infof("Mysql prepare statement error, stmt: %s, err: %s ", mysqlInsertFeeDetailCmd, err)
	} else {
		// 插入失败的话就更新
		_, err = stmt.Exec(yszbh, awbr.Account, awbr.Name, awbr.CurrentPeriod, awbr.Charge, awbr.CurrentMeter, awbr.PreviousMeter,
			awbr.Paid, awbr.OtherFee.Wsf, awbr.OtherFee.Xfft, awbr.OtherFee.Ljf, awbr.OtherFee.Ecjydf, awbr.OtherFee.Szyf, awbr.OtherFee.Cjhys, awbr.OtherFee.Wyj, awbr.OtherFee.Wswyj, awbr.WaterCharge)
		if err != nil {
			if driverErr, ok := err.(*mysql.MySQLError); ok {
				if driverErr.Number == 1062 {
					stmt, err = conn.MysqlDB.Prepare(mysqlUpdateFeeDetailCmd)
					if err != nil {
						logger.GetLogEntry(r).Infof("Mysql prepare statement error, stmt: %s, err: %s ", mysqlUpdateFeeDetailCmd, err)
					}
					// 只更新缴费状态
					_, err = stmt.Exec(awbr.Paid, yszbh)
					if err != nil {
						logger.GetLogEntry(r).Infof("Mysql update fee detail data error, err: %s ", err)
					}
				} else {
					logger.GetLogEntry(r).Infof("Mysql insert fee detail data error, err: %s ", err)
				}
			}
		}
	}
}

