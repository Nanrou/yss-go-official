package api

import (
	"github.com/go-chi/render"
	"net/http"
)

/* Low level error */
type ErrResponse struct {
	Err            error  `json:"-"`               // low-level runtime error
	HTTPStatusCode int    `json:"-"`               // http response status code
	StatusText     string `json:"status"`          // user-level status message
	AppCode        int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText      string `json:"error,omitempty"` // application-level error message, for debugging
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

// 403
var ErrForbidden = &ErrResponse{HTTPStatusCode: 403, StatusText: "Forbidden"}

/* High level error */
type ErrorCodeABC struct {
	Code int         `json:"errcode"`
	Msg  string      `json:"errmsg"`
	Data interface{} `json:"data"` // slice or struct
}

func (e *ErrorCodeABC) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func insABC(ec int, em string, data interface{}) *ErrorCodeABC {
	// reflect.TypeOf(data).Kind()
	return &ErrorCodeABC{
		Code: ec,
		Msg:  em,
		Data: data,
	}
}

// 正常响应，状态码为0
func NewResponseOK(data interface{}) *ErrorCodeABC {
	return insABC(0, "", data)
}

// 请求头部缺失身份证号码
func NewErrMissUInfo() *ErrorCodeABC {
	return insABC(100001, "请求头部缺失身份证号码", "")
}

// 缺失必要参数
func NewErrRequiredFields() *ErrorCodeABC {
	return insABC(100002, "缺失必要参数", "")
}

// 无绑定关系
func NewErrInvalidBinding() *ErrorCodeABC {
	return insABC(100003, "无绑定关系", "")
}

// 错误的用户信息/预留电话
func NewErrInvalidAccountData() *ErrorCodeABC {
	return insABC(100004, "错误的用户信息/预留电话", "")
}

// 绑定关系已存在
func NewErrAlreadyExists() *ErrorCodeABC {
	return insABC(100005, "绑定关系已存在", "")
}

// 内部接口请求失败
func NewErrInnerException() *ErrorCodeABC {
	return insABC(200001, "内部接口请求失败", "")
}
