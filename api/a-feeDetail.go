package api

import (
	"fmt"
	"github.com/go-chi/render"
	"net/http"
	"yss-go-official/logger"
)

var keyFeeDetail = "%s:%s:keyFeeDetail"

// fee_detail
func handleFeeDetail(w http.ResponseWriter, r *http.Request) {
	var err error
	// 检查输入参数
	billId := r.URL.Query().Get("bill_id")
	date := r.URL.Query().Get("date")
	if len(date) == 0 || len(billId) == 0 {
		err = render.Render(w, r, NewErrRequiredFields())
		if err != nil {
			logger.GetLogEntry(r).Info(err)
		}
		return
	}
	if account, ok := checkBinding(r); ok {
		_key := fmt.Sprintf(keyFeeDetail, account, billId)
		if resp, exist := redisGet(r, _key); exist {
			err = render.Render(w, r, NewResponseOK(resp))
			if err != nil {
				logger.GetLogEntry(r).Info(err)
			}
		} else {
			resp := handleFeeDetailInner(r, account, billId, date)
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

func handleFeeDetailInner(r *http.Request, account string, billId string, date string) *feeDetailResponse {
	fd := getFeeDetail(r, billId, date)
	return newFeeDetailResponse(fd, account, date)
}
