package scripts

import (
	"fmt"
	"github.com/bmbstack/ripple"
	"github.com/inkbamboo/tingshu/internal/controllers"
	"github.com/inkbamboo/tingshu/internal/services"
)

func RunServer() {
	controllers.RouteAPI()
	ripple.Default().Run()
}
func getTabList() {
	tabList, _ := services.GetTingService().GetTabList()
	fmt.Printf("******%+v\n", tabList[0])
}
func getBookListByTab() {
	bookList, count, _ := services.GetTingService().GetBookListByTab("/book/1", 1)
	fmt.Printf("******%+v\n", bookList[0])
	fmt.Printf("******%+v\n", count)
}
func getBookInfo() {
	bookInfo, chapterList, chapterCount, _ := services.GetTingService().GetBookInfo("/book/1", "461")
	fmt.Printf("******bookInfo %+v\n", bookInfo)
	fmt.Printf("******chapterList %+v\n", chapterList[0])
	fmt.Printf("******chapterCount %+v\n", chapterCount)
}
func bookPlay() {
	_ = services.GetTingService().BookPlay(nil, "461", "1-316")
}

func Test() {
	//getTabList()
	//getBookListByTab()
	//getBookInfo()
	bookPlay()
}
