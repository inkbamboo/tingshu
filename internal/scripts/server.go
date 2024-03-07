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
	tabList, _ := services.GetShuYinService().GetTabList()
	fmt.Printf("******%+v\n", tabList[0])
}
func getBookListByTab() {
	bookList, count, _ := services.GetShuYinService().GetBookListByTab("1", 1)
	fmt.Printf("******%+v\n", bookList[0])
	fmt.Printf("******%+v\n", count)
}
func getBookInfo() {
	bookInfo, chapterList, chapterCount, _ := services.GetShuYinService().GetBookInfo("1", "121450")
	fmt.Printf("******%+v\n", bookInfo)
	fmt.Printf("******%+v\n", chapterList[0])
	fmt.Printf("******%+v\n", chapterCount)
}
func bookPlay() {
	_ = services.GetShuYinService().BookPlay(nil, "121450", 1)
}

func Test() {
	//getTabList()
	//getBookListByTab()
	//getBookInfo()
	bookPlay()
}
