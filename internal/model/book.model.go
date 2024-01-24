package model

type Book struct {
	BaseModel
	Href       string `json:"href" gorm:"column:href"`               // 链接地址
	Name       string `json:"name" gorm:"column:name"`               // 名称
	Pic        string `json:"pic" gorm:"column:pic"`                 // 图片
	Author     string `json:"author" gorm:"column:author;"`          // 作者
	Speaker    string `json:"speaker" gorm:"column:speaker;"`        // 演播
	CreateTime string `json:"create_time" gorm:"column:create_time"` // 上传时间
}

func (m *Book) TableName() string {
	return "book"
}
