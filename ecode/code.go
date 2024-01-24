package ecode

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

var (
	Success              = Error(1, "成功")
	AuthFailed           = Error(10, "Jwt鉴权失败")
	ServerError          = Error(500, "服务器错误")
	ParamError           = Error(501, "参数错误")
	RpcRequestFailed     = Error(502, "RPC请求失败")
	SignError            = Error(503, "Sign参数错误")
	JwtTokenExpireFailed = Error(602, "用户登陆已失效")
	PowerFailed          = Error(403, "您没有当前操作权限")

	ErrorRpcRequestFailed = Error(107, "失败")
	IPForbidden           = Error(108, "IPForbidden,禁止的IP访问")
	RateLimitError        = Error(109, "您请求太过频繁了")
)

//==============================================
//                  OK
//==============================================

type Response struct {
	Code int         `json:"code"`                                // code码, OK为成功
	Data interface{} `json:"data,omitempty" swaggertype:"object"` // 返回数据
	Msg  string      `json:"msg,omitempty"`                       // 返回消息
	Err  string      `json:"error,omitempty"`                     // 返回内部错误
	Now  int64       `json:"now"`                                 // 服务器时间戳
}

type PageResponse struct {
	Response
	Meta PageMeta `json:"meta" swaggertype:"object"` // 分页数据
}

type PageMeta struct {
	Page      int   `json:"page"`       // 当前页码
	Size      int   `json:"size"`       // 每页长度
	TotalPage int   `json:"total_page"` // 总页数
	TotalItem int64 `json:"total_item"` // 总数据数
}

func Error(code int, msg string) *Response {
	return &Response{Code: code, Msg: msg, Now: time.Now().Unix()}
}

// Status implement error interface
func (r *Response) Error() string {
	return fmt.Sprintf("code=%d, message=%s", r.Code, r.Msg)
}
func (r *Response) ErrorWithMessage(msg string) *Response {
	return Error(r.Code, msg)
}
func SuccessJSON(ctx echo.Context, data interface{}) error {
	return ctx.JSON(http.StatusOK, Response{
		Code: Success.Code,
		Data: data,
		Msg:  Success.Msg,
		Err:  "",
		Now:  time.Now().Unix(),
	})
}
