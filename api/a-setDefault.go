package api

import (
	"encoding/json"
	"github.com/go-chi/render"
	"net/http"
	"yss-go-official/logger"
	conn "yss-go-official/orm"
)

type SetDefaultData struct {
	Account string `json:"account"`
}

func (b *SetDefaultData) CheckValidate() bool {
	return len(b.Account) > 0
}

func handleSetDefault(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var setDefaultData SetDefaultData
	err := decoder.Decode(&setDefaultData)
	if err != nil {
		logger.GetLogEntry(r).Info(err)
	}
	if setDefaultData.CheckValidate() {

		idn := r.Context().Value("id")
		errFlag := false // 内部错误标记符

		// 检查绑定关系
		stmt, err := conn.MysqlDB.Prepare(mysqlQueryAccountCheckCmd)
		if err != nil {
			logger.GetLogEntry(r).Infof("Mysql prepare statement error, stmt: %s, err: %s ", mysqlQueryAccountCmd, err)
		}
		var useless string
		err = stmt.QueryRow(idn, setDefaultData.Account).Scan(&useless)
		if err != nil {
			err = render.Render(w, r, NewErrInvalidBinding())
			if err != nil {
				logger.GetLogEntry(r).Info(err)
			}
		} else {
			// 正常逻辑
			res := NewResponseOK("")
			stmt, err = conn.MysqlDB.Prepare(mysqlSetDefaultAccountCmd)
			if err != nil {
				logger.GetLogEntry(r).Infof("Mysql prepare statement error, stmt: %s, err: %s ", mysqlDeleteUserDataCmd, err)
				errFlag = true
			} else {
				// 事务
				_, err := stmt.Exec(setDefaultData.Account, idn)
				if err != nil {
					logger.GetLogEntry(r).Infof("Mysql exec error, stmt: %s, err: %s ", mysqlDeleteUserDataCmd, err)
					errFlag = true
				}
				if errFlag {
					res = NewErrInnerException()
				}
				err = render.Render(w, r, res)
				if err != nil {
					logger.GetLogEntry(r).Info(err)
				}
			}
		}
	} else {
		err = render.Render(w, r, NewErrRequiredFields())
		if err != nil {
			logger.GetLogEntry(r).Info(err)
		}
	}
}

