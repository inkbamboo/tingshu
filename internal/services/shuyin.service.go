package services

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/bmbstack/ripple"
	"github.com/go-resty/resty/v2"
	"github.com/inkbamboo/tingshu/dto"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	"github.com/tidwall/gjson"
	"strconv"
	"strings"
	"sync"
)

var (
	shuYinService     *ShuYinService
	shuYinServiceOnce sync.Once
)

func GetShuYinService() *ShuYinService {
	shuYinServiceOnce.Do(func() {
		shuYinService = &ShuYinService{}
	})
	return shuYinService
}

type ShuYinService struct {
}

func (s *ShuYinService) parseTabList(doc *goquery.Selection) (list []*dto.TabItem, err error) {
	list = []*dto.TabItem{}
	doc.Each(func(i int, aItem *goquery.Selection) {
		href, _ := aItem.Attr("href")
		href = strings.TrimSuffix(href, "-0.html")
		href = strings.TrimPrefix(href, "/listinfo-")
		list = append(list, &dto.TabItem{
			Name: aItem.Text(),
			Key:  href,
		})
	})
	return
}

func (s *ShuYinService) GetTabList() (list []*dto.TabItem, err error) {
	list = []*dto.TabItem{}
	var resp *resty.Response
	resp, err = resty.New().R().Get(fmt.Sprintf(`%s%s`, ripple.GetConfig().GetString("baseUrl.shuYin"), "/yousheng/"))
	if err != nil {
		return
	}
	htmlStr := string(resp.Body())
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlStr))
	if err != nil {
		fmt.Printf("Failed to parse HTML: %v", err)
	}
	list, err = s.parseTabList(doc.Find("dl").Find(".clearfix").Find("a"))
	return
}

func (s *ShuYinService) GetBookListByTab(tab string, index int64) (bookList []*dto.BookItem, totalCount int64, err error) {
	var resp *resty.Response
	resp, err = resty.New().R().Get(fmt.Sprintf(`%s/listinfo-%s-%d.html`, ripple.GetConfig().GetString("baseUrl.shuYin"), tab, index))
	if err != nil {
		return
	}
	bookList, totalCount = s.getBookListFromHtml(string(resp.Body()))

	bookList = lo.Map(bookList, func(item *dto.BookItem, index int) *dto.BookItem {
		item.Tab = tab
		return item
	})
	return
}
func (s *ShuYinService) getBookListFromHtml(htmlStr string) (bookList []*dto.BookItem, totalCount int64) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlStr))
	if err != nil {
		fmt.Printf("Failed to parse HTML: %v", err)
	}
	tabList, _ := s.parseTabList(doc.Find(".classify-row").Find("a"))
	tabInfos := lo.KeyBy(tabList, func(item *dto.TabItem) string {
		return item.Key
	})
	fmt.Printf("*******%+v\n", htmlStr)
	doc.Find(".qm-mod-tb").Find("li").Each(func(i int, bItem *goquery.Selection) {
		bookInfo := &dto.BookItem{}
		bookInfo.Pic, _ = bItem.Find("img").Attr("src")
		nItem := bItem.Find(".s-tit")
		href, _ := nItem.Find("a").Attr("href")
		bookInfo.Tab, bookInfo.BookId = s.parseHref(href)
		bookInfo.Name = nItem.Text()
		bookInfo.Speaker = bItem.Find(".s-tags").Text()
		bookInfo.Summary = bItem.Find(".s-des").Text()
		statusInfo := bItem.Find(".s-name").Text()
		statusInfo = strings.Replace(statusInfo, string([]byte{0xc2, 0xa0, 0xc2, 0xa0}), " ", -1)
		arr := strings.Split(statusInfo, " ")
		bookInfo.Status = arr[0]
		bookInfo.CreateTime = strings.TrimSuffix(arr[1], "更新")
		if tabInfos[bookInfo.Tab] != nil {
			bookInfo.TabName = tabInfos[bookInfo.Tab].Name
		}
		bookList = append(bookList, bookInfo)
	})
	totalCountStr := doc.Find(".qm-page-number").Find("a").First().Find("b").Text()
	totalCount, err = strconv.ParseInt(totalCountStr, 10, 64)
	return
}
func (s *ShuYinService) parseHref(href string) (tab, bookId string) {
	href = strings.TrimPrefix(href, "/album/")
	href = strings.TrimSuffix(href, ".html")
	hrefArr := strings.Split(href, "-")
	tab = hrefArr[0]
	bookId = hrefArr[len(hrefArr)-1]
	return
}
func (s *ShuYinService) parseIndexList(doc *goquery.Selection) (bookList []*dto.BookItem) {
	doc.Find("li").Each(func(i int, lItem *goquery.Selection) {
		bookInfo := &dto.BookItem{}
		bookInfo.Pic, _ = lItem.Find("img").Attr("src")
		Info := lItem.Find(".info")
		bookInfo.Name = Info.Find("a").Text()
		href, _ := Info.Find("a").Attr("href")
		bookInfo.Tab, bookInfo.BookId = s.parseHref(href)
		bookInfo.Summary = Info.Find(".detail").Text()
		bookList = append(bookList, bookInfo)
	})
	return
}
func (s *ShuYinService) GetBookIndexList() (tabList []*dto.TabItem, recommendList []*dto.BookItem, newList []*dto.BookItem, err error) {
	var resp *resty.Response
	resp, err = resty.New().R().Get(fmt.Sprintf(`%s%s`, ripple.GetConfig().GetString("baseUrl.shuYin"), "/yousheng/"))
	if err != nil {
		fmt.Printf("****************%v\n", err)
		return
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(resp.Body())))
	if err != nil {
		fmt.Printf("Failed to parse HTML: %v", err)
	}
	doc.Find(".unusual-list").Each(func(i int, aItem *goquery.Selection) {
		switch i {
		case 0:
			newList = s.parseIndexList(aItem)
		case 1:
			recommendList = s.parseIndexList(aItem)
		}
	})
	tabList, err = s.parseTabList(doc.Find("dl").Find(".clearfix").Find("a"))
	return
}

func (s *ShuYinService) GetBookInfo(tab, bookId string) (bookInfo *dto.BookItem, chapterList []*dto.ChapterItem, chapterCount int64, err error) {
	var resp *resty.Response
	resp, err = resty.New().R().Get(fmt.Sprintf(`%s/album/%s-%s.html`, ripple.GetConfig().GetString("baseUrl.shuYin"), tab, bookId))
	if err != nil {
		return
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(resp.Body())))
	if err != nil {
		fmt.Printf("Failed to parse HTML: %v", err)
	}
	bookInfo = &dto.BookItem{
		BookId: bookId,
		Tab:    tab,
	}
	bookInfo.Pic, _ = doc.Find(".bookprofile").Find("dt").Find("img").Attr("src")
	bookInfoDoc := doc.Find(".bookprofile").Find("dd")
	bookInfo.Name = bookInfoDoc.Find(".title").Text()
	bookInfoDoc.Find(".sub-cols").Find("span").Each(func(i int, item *goquery.Selection) {
		switch i {
		case 0:
			bookInfo.Status = strings.Split(item.Text(), "：")[1]
		case 1:
			bookInfo.TabName = strings.Split(item.Text(), "：")[1]
			bookInfo.TabName = strings.Split(bookInfo.TabName, "-")[1]
		case 2:
			bookInfo.CreateTime = strings.Split(item.Text(), "：")[1]
		}
	})
	var speakers []string
	bookInfoDoc.Find(".sub-tags").Find("a").Each(func(i int, aItem *goquery.Selection) {
		speakers = append(speakers, aItem.Text())
	})
	bookInfo.Speaker = strings.Join(speakers, ",")
	bookInfo.Summary = doc.Find(".introcontent").Find("dd").Text()
	doc.Find(".compress").Find("a").Each(func(i int, item *goquery.Selection) {
		chapterList = append(chapterList, &dto.ChapterItem{
			Name:      item.Text(),
			ChapterId: i + 1,
		})
	})
	chapterCount = int64(len(chapterList))
	if chapterCount >= 10 {
		chapterList = chapterList[0:10]
	}
	return
}
func (s *ShuYinService) BookPlay(ctx echo.Context, bookId string, chapter int64) (err error) {
	var resp *resty.Response
	resp, err = resty.New().R().SetHeader("content-type", "application/x-www-form-urlencoded; charset=UTF-8").SetFormData(map[string]string{
		"bookId": bookId,
		"page":   fmt.Sprintf("%d", chapter),
	}).Post(fmt.Sprintf(`%s/?s=api-getneoplay`, ripple.GetConfig().GetString("baseUrl.shuYin")))
	if err != nil {
		return
	}
	urlStr := gjson.Get(string(resp.Body()), "url").String()
	fmt.Printf("********%+v\n", urlStr)
	return
}
func (s *ShuYinService) Search(keyword string, index int64) (bookList []*dto.BookItem, totalCount int64, err error) {
	var resp *resty.Response
	resp, err = resty.New().R().SetFormData(map[string]string{
		"show":     "title,newstext,player,playadmin",
		"keyboard": keyword,
	}).Post(fmt.Sprintf(`%s/e/search/index.php`, ripple.GetConfig().GetString("baseUrl.shuYin")))
	if err != nil {
		return
	}

	fmt.Printf(string(resp.Body()))
	//bookList, totalCount = s.getBookListFromHtml(string(resp.Body()))
	return
}
