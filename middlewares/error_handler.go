package middlewares

import (
	"github.com/inkbamboo/tingshu/ecode"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

func HTTPErrorHandler() func(error, echo.Context) {
	return func(err error, ctx echo.Context) {
		if ctx.Response().Committed {
			return
		}
		ctx.Response().Header().Set("Access-Control-Allow-Origin", "*")
		he, ok := err.(*echo.HTTPError)
		code := http.StatusOK
		var message interface{}
		if ok {
			if he.Internal != nil {
				if herr, ok := he.Internal.(*echo.HTTPError); ok {
					he = herr
				}
			}
			code = he.Code
			message = he.Message
			if m, ok := he.Message.(string); ok {
				message = map[string]interface{}{"message": m, "error": err.Error()}
			}
		} else {
			message, ok = err.(*ecode.Response)
			if !ok {
				message = &ecode.Response{
					Code: ecode.ServerError.Code,
					Msg:  err.Error(),
					Err:  err.Error(),
					Now:  time.Now().Unix(),
				}
			}
		}

		// Send response
		if ctx.Request().Method == http.MethodHead { // Issue #608
			err = ctx.NoContent(code)
		} else {
			err = ctx.JSON(code, message)
		}
		if err != nil {
			ctx.Logger().Error(err)
		}
	}
}
