package api

import (
	"database/sql"
	"encoding/json"
	"github.com/go-chi/render"
	"github.com/go-sql-driver/mysql"
	"net/http"
	"yss-go-official/logger"
	conn "yss-go-official/orm"
)

type BindData struct {
	Name    string `json:"account_name"`
	Account string `json:"account"`
	Phone   string `json:"phone"`
}

func (b *BindData) CheckValidate() bool {
	return len(b.Name) > 0 && len(b.Account) > 0 && len(b.Phone) > 0
}

func handleBind(w http.ResponseWriter, r *http.Request) {
	//dump, err := httputil.DumpRequest(r, true)
	//if err != nil{
	//	panic(err)
	//}
	//fmt.Println(string(dump))
	//body, err := ioutil.ReadAll(r.Body)
	//if err != nil{
	//	panic(err)
	//}
	//println(body)
	decoder := json.NewDecoder(r.Body)
	var bindData BindData
	err := decoder.Decode(&bindData)
	if err != nil {
		logger.GetLogEntry(r).Info(err)
	}
	if bindData.CheckValidate() {
		idn := r.Context().Value("id")
		errFlag := false // 内部错误标记符
		res := NewResponseOK("")
		if checkAccountValidate(bindData) {
			stmt, err := conn.MysqlDB.Prepare(mysqlQueryBindingCmd)
			if err != nil {
				logger.GetLogEntry(r).Infof("Mysql prepare statement error, stmt: %s, err: %s ", mysqlQueryBindingCmd, err)
				errFlag = true
			} else {
				var bindingRecord wechatProfile
				err = stmt.QueryRow(idn).Scan(&bindingRecord.id, &bindingRecord.idn, &bindingRecord.defaultAccount)
				if err != nil {
					if err.Error() == "sql: no rows in result set" {
						stmt, err := conn.MysqlDB.Prepare(mysqlCreateWechatCmd)
						if err != nil {
							logger.GetLogEntry(r).Infof("Mysql prepare statement error, stmt: %s, err: %s ", mysqlQueryBindingCmd, err)
							errFlag = true
						} else {
							_, err = stmt.Exec(idn, bindData.Account)
							if err != nil {
								logger.GetLogEntry(r).Infof("Mysql exec error, stmt: %s, err: %s ", mysqlCreateWechatCmd, err)
								errFlag = true
							}
						}
					} else {
						logger.GetLogEntry(r).Infof("Mysql statement exec error, params: %s, err: %s", err)
						errFlag = true
					}
				}
				stmt, err = conn.MysqlDB.Prepare(mysqlCreateUserDataCmd)
				if err != nil {
					logger.GetLogEntry(r).Infof("Mysql prepare statement error, stmt: %s, err: %s ", mysqlCreateUserDataCmd, err)
					errFlag = true
				} else {
					_, err = stmt.Exec(idn, bindData.Account, bindData.Name, bindData.Phone)
					if driverErr, ok := err.(*mysql.MySQLError); ok {
						if driverErr.Number == 1062 {
							res = NewErrAlreadyExists()
						} else {
							logger.GetLogEntry(r).Infof("Mysql exec error, stmt: %s, err: %s ", mysqlCreateUserDataCmd, err)
							errFlag = true
						}
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
		} else {
			err = render.Render(w, r, NewErrInvalidAccountData())
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

func checkAccountValidate(data BindData) bool {
	var _useless string
	err := conn.MssqlDB.QueryRow(mssqlQueryCheckAccountCmd,
		sql.Named("account", data.Account),
		sql.Named("name", data.Name)).Scan(&_useless)
	if err != nil {
		return false
	}
	return true
}
