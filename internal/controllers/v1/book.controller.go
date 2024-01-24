package v1

import (
	"github.com/inkbamboo/tingshu/dto"
	"github.com/inkbamboo/tingshu/ecode"
	"github.com/inkbamboo/tingshu/internal/services"
	"github.com/labstack/echo/v4"
)

type BookController struct {
	Group *echo.Group
}

func NewBookController(group *echo.Group) *BookController {
	return &BookController{
		Group: group,
	}
}
func (c BookController) Setup() {
	c.Group.GET("/indexList", c.GetBookIndexList)
	c.Group.GET("/listByTab", c.GetBookListByTab)
	c.Group.GET("/tabList", c.GetTabList)
	c.Group.GET("/search", c.Search)
	c.Group.GET("/bookInfo", c.GetBookInfo)
	c.Group.GET("/play", c.BookPlay)
}

// GetTabList
// @Summary  获取分类列表
// @Description
// @Security  TokenAuth
// @Tags 书籍
// @Accept json
// @Produce  json
// @Param body query dto.BaseIn  true  "请求结构"
// @Success  200  {object}  ecode.Response{data=dto.GetTabListOut}
// @Router /v1/book/tabList [get]
func (c BookController) GetTabList(ctx echo.Context) (err error) {
	req := &dto.BaseIn{}
	res := &dto.GetTabListOut{}
	if err = ctx.Bind(req); err != nil {
		return err
	}
	var tingShuHandler services.TingShuInterface
	if tingShuHandler, err = services.NewInstance(req.Channel); err != nil {
		return
	}
	res.List, _ = tingShuHandler.GetTabList()
	return ecode.SuccessJSON(ctx, res)
}

// GetBookListByTab
// @Summary  获取书籍列表
// @Description
// @Security  TokenAuth
// @Tags 书籍
// @Accept json
// @Produce  json
// @Param body query dto.GetBookListByTabIn  true  "请求结构"
// @Success  200  {object}  ecode.Response{data=dto.GetBookListByTabOut}
// @Router /v1/book/listByTab [get]
func (c BookController) GetBookListByTab(ctx echo.Context) (err error) {
	req := &dto.GetBookListByTabIn{}
	res := &dto.GetBookListByTabOut{}
	if err = ctx.Bind(req); err != nil {
		return err
	}
	var tingShuHandler services.TingShuInterface
	if tingShuHandler, err = services.NewInstance(req.Channel); err != nil {
		return
	}
	if res.List, res.TotalCount, err = tingShuHandler.GetBookListByTab(req.Tab, req.Page); err != nil {
		return
	}
	return ecode.SuccessJSON(ctx, res)
}

// Search
// @Summary  获取书籍列表
// @Description
// @Security  TokenAuth
// @Tags 书籍
// @Accept json
// @Produce  json
// @Param body query dto.SearchIn  true  "请求结构"
// @Success  200  {object}  ecode.Response{data=dto.SearchOut}
// @Router /v1/book/search [get]
func (c BookController) Search(ctx echo.Context) (err error) {
	req := &dto.SearchIn{}
	res := &dto.SearchOut{}
	if err = ctx.Bind(req); err != nil {
		return err
	}
	var tingShuHandler services.TingShuInterface
	if tingShuHandler, err = services.NewInstance(req.Channel); err != nil {
		return
	}
	if res.List, res.TotalCount, err = tingShuHandler.Search(req.Keyword, req.Page); err != nil {
		return
	}
	return ecode.SuccessJSON(ctx, res)
}

// GetBookIndexList
// @Summary  获取分类列表
// @Description
// @Security  TokenAuth
// @Tags 书籍
// @Accept json
// @Produce  json
// @Param body query dto.BaseIn  true  "请求结构"
// @Success  200  {object}  ecode.Response{data=dto.GetBookIndexListOut}
// @Router /v1/book/indexList [get]
func (c BookController) GetBookIndexList(ctx echo.Context) (err error) {
	req := &dto.BaseIn{}
	res := &dto.GetBookIndexListOut{}
	if err = ctx.Bind(req); err != nil {
		return err
	}
	var tingShuHandler services.TingShuInterface
	if tingShuHandler, err = services.NewInstance(req.Channel); err != nil {
		return
	}
	if res.TabList, res.RecommendList, res.NewList, err = tingShuHandler.GetBookIndexList(); err != nil {
		return
	}
	return ecode.SuccessJSON(ctx, res)
}

// GetBookInfo
// @Summary  获取书籍列表
// @Description
// @Security  TokenAuth
// @Tags 书籍
// @Accept json
// @Produce  json
// @Param body query dto.GetBookInfoIn  true  "请求结构"
// @Success  200  {object}  ecode.Response{data=dto.GetBookInfoOut}
// @Router /v1/book/bookInfo [get]
func (c BookController) GetBookInfo(ctx echo.Context) (err error) {
	req := &dto.GetBookInfoIn{}
	res := &dto.GetBookInfoOut{}
	if err = ctx.Bind(req); err != nil {
		return err
	}
	var tingShuHandler services.TingShuInterface
	if tingShuHandler, err = services.NewInstance(req.Channel); err != nil {
		return
	}
	if res.BookInfo, res.ChapterList, res.ChapterCount, err = tingShuHandler.GetBookInfo(req.Tab, req.BookId); err != nil {
		return
	}
	return ecode.SuccessJSON(ctx, res)
}

// BookPlay
// @Summary  获取书籍列表
// @Description
// @Security  TokenAuth
// @Tags 书籍
// @Accept json
// @Produce  json
// @Param body query dto.BookPlayIn  true  "请求结构"
// @Router /v1/book/play [get]
func (c BookController) BookPlay(ctx echo.Context) (err error) {
	req := &dto.BookPlayIn{}
	if err = ctx.Bind(req); err != nil {
		return err
	}
	var tingShuHandler services.TingShuInterface
	if tingShuHandler, err = services.NewInstance(req.Channel); err != nil {
		return
	}
	return tingShuHandler.BookPlay(ctx, req.BookId, req.Chapter)
}
