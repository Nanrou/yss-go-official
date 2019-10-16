package api

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"net/http"
	"yss-go-official/logger"
	"yss-go-official/meta"
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

func init() {
	YssRouter = chi.NewRouter()
	YssRouter.Use(middleware.RealIP)
	YssRouter.Use(render.SetContentType(render.ContentTypeJSON))
	YssRouter.Use(WhiteListMiddleware)
	YssRouter.Use(InjectIdentifierMiddleware)
	YssRouter.Route("/yss", func(r chi.Router) {
		r.Get("/right_content", handleRightContent)
		r.Get("/sites", handleSites)
		r.Get("/accounts", handleAccounts)
		r.Get("/default_account", handleDefaultAccount)
		r.Get("/account_without_binding", handleAccountWithoutBinding)
		//r.Get("/fee_history", handleFeeHistory)
		//r.Get("/fee_history_without_paid", handleFeeHistoryWithoutPaid)
		// r.Get("/redis_in", handleRedisIn)
		// r.Get("/redis_out", handleRedisOut)
	})
}
