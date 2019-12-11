package api

import (
	"encoding/json"
	"github.com/go-chi/render"
	"net/http"
	"yss-go-official/logger"
	conn "yss-go-official/orm"
)

type CancelData struct {
	Account string `json:"account"`
}

func (b *CancelData) CheckValidate() bool {
	return len(b.Account) > 0
}

func handleCancelBind(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var cancelData CancelData
	err := decoder.Decode(&cancelData)
	if err != nil {
		logger.GetLogEntry(r).Info(err)
	}
	if cancelData.CheckValidate() {
		idn := r.Context().Value("id")
		errFlag := false // 内部错误标记符
		res := NewResponseOK("")
		stmt, err := conn.MysqlDB.Prepare(mysqlDeleteUserDataCmd)
		if err != nil {
			logger.GetLogEntry(r).Infof("Mysql prepare statement error, stmt: %s, err: %s ", mysqlDeleteUserDataCmd, err)
			errFlag = true
		} else {
			// 事务
			_, err := stmt.Exec(idn, cancelData.Account)
			if err != nil {
				logger.GetLogEntry(r).Infof("Mysql exec error, stmt: %s, err: %s ", mysqlDeleteUserDataCmd, err)
				errFlag = true
			} else {
				// 前面删除成功才来这里删除默认账户
				stmt, err = conn.MysqlDB.Prepare(mysqlUnsetDefaultCmd)
				if err != nil {
					logger.GetLogEntry(r).Infof("Mysql prepare statement error, stmt: %s, err: %s ", mysqlUnsetDefaultCmd, err)
					errFlag = true
				} else {
					_, err := stmt.Exec(idn, cancelData.Account)
					if err != nil {
						logger.GetLogEntry(r).Infof("Mysql exec error, stmt: %s, err: %s ", mysqlUnsetDefaultCmd, err)
						errFlag = true
					}
				}
			}
			if errFlag {
				res = NewErrInnerException()
			}
			err = render.Render(w, r, res)
			if err != nil {
				logger.GetLogEntry(r).Info(err)
			}
		}
	} else {
		err = render.Render(w, r, NewErrRequiredFields())
		if err != nil {
			logger.GetLogEntry(r).Info(err)
		}
	}
}
