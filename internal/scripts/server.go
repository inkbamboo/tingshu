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
func Test() {
	_, recommendList, newList, _ := services.GetNianYinService().GetBookIndexList()
	fmt.Printf("******%+v\n", recommendList[0])
	fmt.Printf("******%+v\n", newList[0])

}
