package controllers

import (
	"github.com/bmbstack/ripple"
	"github.com/inkbamboo/tingshu/docs"
	_ "github.com/inkbamboo/tingshu/docs"
	v1 "github.com/inkbamboo/tingshu/internal/controllers/v1"
	"github.com/inkbamboo/tingshu/middlewares"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"net/http"
)

func RouteAPI() {
	echoMux := ripple.Default().GetEcho()
	// 运维用
	echoMux.GET("/ping", func(ctx echo.Context) error {
		return ctx.JSON(http.StatusOK, "OK")
	})
	echoMux.HTTPErrorHandler = middlewares.HTTPErrorHandler()
	echoMux.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	//注册swagger文档
	registerXcxAPI(echoMux)
	registerSwagger(echoMux)
}
func registerXcxAPI(echoMux *echo.Echo) {
	v1Group := echoMux.Group("/v1")
	book := v1.NewBookController(v1Group.Group("/book"))
	book.Setup()
}
func registerSwagger(echoMux *echo.Echo) {
	if !ripple.GetConfig().GetBool("swaggeApi") {
		return
	}
	// swagger 设置
	docs.SwaggerInfo.Title = "管理后台服务"
	docs.SwaggerInfo.Description = "管理后台服务-API文档 Schemes 选择HTTPS"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = ripple.GetConfig().GetString("swaggeApiHost")
	docs.SwaggerInfo.BasePath = "/v1"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
	echoMux.GET("/apidoc/*", echoSwagger.WrapHandler)
}
