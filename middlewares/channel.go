package middlewares

import (
	"github.com/labstack/echo/v4"
)

const ChannelKey = "channel"

func Channel() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			c.Set(ChannelKey, c.Request().Header.Get(ChannelKey))
			return next(c)
		}
	}
}
