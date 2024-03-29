package dto

type BookItem struct {
	BookId     string `json:"book_id"`     // 链接地址
	Name       string `json:"name"`        // 名称
	Tab        string `json:"tab"`         // 分类
	TabName    string `json:"tab_name"`    // 分类名称
	Pic        string `json:"pic"`         // 图片
	Author     string `json:"author"`      // 作者
	Speaker    string `json:"speaker"`     // 演播
	Status     string `json:"status"`      // 演播
	Summary    string `json:"summary"`     // 简介
	CreateTime string `json:"create_time"` // 上传时间
}
type ChapterItem struct {
	Name      string `json:"name"`
	ChapterId string `json:"chapter_id"`
}
type GetBookListByTabIn struct {
	Page int64  `json:"page" form:"page" binding:"required"`
	Tab  string `json:"tab" form:"tab" binding:"required"` // 名称
}

type GetBookListByTabOut struct {
	List       []*BookItem `json:"list"`        // 链接地址
	TotalCount int64       `json:"total_count"` // 名称
}
type SearchIn struct {
	Page    int64  `json:"page" form:"page" binding:"required"`
	Keyword string `json:"keyword" form:"keyword" binding:"required"` // 名称
}
type SearchOut struct {
	List       []*BookItem `json:"list"`        // 链接地址
	TotalCount int64       `json:"total_count"` // 名称
}

type GetBookIndexListOut struct {
	TabList       []*TabItem  `json:"tab_list"`       // 链接地址
	RecommendList []*BookItem `json:"recommend_list"` // 链接地址
	NewList       []*BookItem `json:"new_list"`       // 链接地址
}
type GetBookInfoIn struct {
	Tab    string `json:"tab" form:"tab" binding:"required"`         // 名称
	BookId string `json:"book_id" form:"book_id" binding:"required"` // 名称
}
type GetBookInfoOut struct {
	BookInfo     *BookItem      `json:"book_info"`    // 链接地址
	ChapterList  []*ChapterItem `json:"chapter_list"` // 名称
	ChapterCount int64          `json:"chapter_count"`
}
type BookPlayIn struct {
	Chapter string `json:"chapter" form:"chapter" binding:"required"` // 名称
	BookId  string `json:"book_id" form:"book_id" binding:"required"` // 名称
}
type BookPlayOut struct {
	Url string `json:"url"` // 链接地址
}
type TabItem struct {
	Name string `json:"name"` // 名称
	Key  string `json:"key"`  // 分类
}
type GetTabListOut struct {
	List []*TabItem `json:"list"` // 链接地址
}
