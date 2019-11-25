package api

import (
	"database/sql"
	"github.com/go-chi/render"
	"net/http"
	"yss-go-official/logger"
	conn "yss-go-official/orm"
)

// default_account
func handleDefaultAccount(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(string)

	// 获取默认户号
	var defaultAccount sql.NullString
	var as accountStation

	stmt, err := conn.MysqlDB.Prepare(mysqlQueryDefaultAccountCmd)
	if err != nil {
		logger.GetLogEntry(r).Infof("Mysql prepare statement error, stmt: %s, err: %s ", mysqlQueryDefaultAccountCmd, err)
	} else {
		err = stmt.QueryRow(id).Scan(&defaultAccount)
		// 无默认账号的直接返回空的as
		if err == nil {
			if defaultAccount.Valid {
				as = handleSingleAccountSync(defaultAccount.String, defaultAccount.String, r)
			}
		} else {
			if err.Error() != "sql: no rows in result set" {
				logger.GetLogEntry(r).Infof("Mysql statement exec error, params: %s, err: %s", id, err)
			}
		}
	}
	err = render.Render(w, r, NewResponseOK(newDefaultAccountResponse(&as)))
	if err != nil {
		logger.GetLogEntry(r).Info(err)
	}
}
