package services

import (
	"fmt"
	"github.com/inkbamboo/tingshu/consts"
	"github.com/inkbamboo/tingshu/dto"
	"github.com/labstack/echo/v4"
)

type TingShuInterface interface {
	GetTabList() (list []*dto.TabItem, err error)
	GetBookListByTab(tab string, index int64) (bookList []*dto.BookItem, totalCount int64, err error)
	GetBookIndexList() (tabList []*dto.TabItem, recommendList []*dto.BookItem, newList []*dto.BookItem, err error)
	GetBookInfo(tab, bookId string) (bookInfo *dto.BookItem, chapterList []*dto.ChapterItem, chapterCount int64, err error)
	BookPlay(ctx echo.Context, bookId string, chapter string) (err error)
	Search(keyword string, index int64) (bookList []*dto.BookItem, totalCount int64, err error)
}

func NewInstance(tingShuType string) (tingShuInterface TingShuInterface, err error) {
	if tingShuType == consts.NianYin.Name() {
		return GetNianYinService(), nil
	}
	switch tingShuType {
	case consts.NianYin.Name():
		tingShuInterface = GetNianYinService()
	case consts.ShuYin.Name():
		tingShuInterface = GetShuYinService()
	case consts.Ting.Name():
		tingShuInterface = GetTingService()
	}
	if tingShuInterface == nil {
		return tingShuInterface, fmt.Errorf("current placeorder type:%v is not match", tingShuType)
	}
	return tingShuInterface, nil
}
