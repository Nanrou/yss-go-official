package api

import (
	"context"
	"github.com/go-chi/render"
	"net"
	"net/http"
	"strings"

	"yss-go-official/logger"
	"yss-go-official/orm"
)

var trustNetwork []*net.IPNet
var trustIp []string

func init() {
	config := orm.GetConfig()
	for _, network := range config.GetTrustNetwork() {
		_, n, err := net.ParseCIDR(network)
		if err == nil {
			trustNetwork = append(trustNetwork, n)
		}
	}
	trustIp = config.GetTrustIps()
}

/** 白名单校验 */
func WhiteListMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		trusted := false

		// RemoteAddr包含端口号，要去掉
		var _rip string
		if _index := strings.Index(r.RemoteAddr, ":"); _index > 0 {
			_rip = r.RemoteAddr[:_index]
		} else {
			_rip = r.RemoteAddr
		}

		// 检查是否在白名单
		for _, _ip := range trustIp {
			if _ip == _rip {
				trusted = true
				break
			}
		}

		// 检查是否在可信子网
		if !trusted {
			_ip := net.ParseIP(_rip)
			for _, _ipNet := range trustNetwork {
				if _ipNet.Contains(_ip) {
					trusted = true
					break
				}
			}
		}

		if trusted {
			next.ServeHTTP(w, r)
		} else {
			// return 403
			logger.GetLogEntry(r).Info("Invalid ip ", r.RemoteAddr)
			err := render.Render(w, r, ErrForbidden)
			if err != nil {
				logger.GetLogEntry(r).Info("WhiteListMiddleware error", err)
			}
			return
		}
	})
}

/** 将头部身份证信息加到上下文的中间件 */
func InjectIdentifierMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, ok := r.Header["X-Tif-Uinfo"]
		if ok && len(id) == 1 {
			ctx := context.WithValue(r.Context(), "id", id[0])
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			err := render.Render(w, r, NewErrMissUInfo())
			if err != nil {
				logger.GetLogEntry(r).Info("InjectIdentifierMiddleware error", err)
			}
			return
		}
	})
}
