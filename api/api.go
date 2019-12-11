package api

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"net/http"
	"time"
	"yss-go-official/logger"
	"yss-go-official/meta"
	conn "yss-go-official/orm"
)

var YssRouter chi.Router

// right_content
func handleRightContent(w http.ResponseWriter, r *http.Request) {
	err := render.Render(w, r, NewResponseOK(newRightContentResponse(meta.Content)))
	if err != nil {
		logger.GetLogEntry(r).Info(err)
	}
}

// sites
func handleSites(w http.ResponseWriter, r *http.Request) {
	err := render.Render(w, r, NewResponseOK(meta.Sites))
	if err != nil {
		logger.GetLogEntry(r).Info(err)
	}
}

// 检查绑定
func checkBinding(r *http.Request) (string, bool) {
	id := r.Context().Value("id")
	account := r.URL.Query().Get("account")

	stmt, err := conn.MysqlDB.Prepare(mysqlQueryAccountCheckCmd)
	if err != nil {
		logger.GetLogEntry(r).Infof("Mysql prepare statement error, stmt: %s, err: %s ", mysqlQueryAccountCmd, err)
		return "", false
	}
	var useless string
	err = stmt.QueryRow(id, account).Scan(&useless)
	if err != nil {
		return "", false
	} else {
		return account, true
	}
}

// Get
func redisGet(r *http.Request, key string) (map[string]interface{}, bool) {
	var res map[string]interface{}
	var exist bool
	if ok, err := conn.RedisConn.Exists(key).Result(); err == nil {
		if ok == 1 {
			if s, err := conn.RedisConn.Get(key).Result(); err == nil {
				if err = json.Unmarshal([]byte(s), &res); err == nil {
					exist = true
				} else {
					logger.GetLogEntry(r).Info("json unmarshal error, err: ", err)
				}
			} else {
				logger.GetLogEntry(r).Infof("redis get error, key: %s, err: %s", key, err)
			}
		}
	} else {
		logger.GetLogEntry(r).Infof("redis exist error, key: %s, err: %s", key, err)
	}
	return res, exist
}

// Set
func redisSet(r *http.Request, key string, value interface{}) {
	if content, err := json.Marshal(value); err == nil {
		_, err := conn.RedisConn.Set(key, content, time.Hour*1).Result()
		if err != nil {
			logger.GetLogEntry(r).Info("redis set error, err: ", err)
		}
	} else {
		logger.GetLogEntry(r).Info("json marshal error, err: ", err)
	}
}

func init() {
	YssRouter = chi.NewRouter()
	YssRouter.Use(middleware.RealIP)
	YssRouter.Use(logger.Middleware)
	YssRouter.Use(render.SetContentType(render.ContentTypeJSON))
	YssRouter.Use(WhiteListMiddleware)
	YssRouter.Use(InjectIdentifierMiddleware)
	YssRouter.Route("/yss", func(r chi.Router) {
		r.Get("/right_content", handleRightContent)
		r.Get("/sites", handleSites)
		r.Get("/accounts", handleAccounts)
		r.Get("/default_account", handleDefaultAccount)
		r.Get("/account_without_binding", handleAccountWithoutBinding)
		r.Get("/fee_history", handleFeeHistory)
		r.Get("/fee_history_without_paid", handleFeeHistoryWithoutPaid)
		r.Get("/fee_detail", handleFeeDetail)
		r.Post("/bind", handleBind)
		r.Post("/cancel_bind", handleCancelBind)
		// r.Get("/redis_in", handleRedisIn)
		// r.Get("/redis_out", handleRedisOut)
	})
}
