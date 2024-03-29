package services

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/bmbstack/ripple"
	"github.com/go-resty/resty/v2"
	"github.com/inkbamboo/tingshu/dto"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	"github.com/spf13/cast"
	"regexp"
	"strings"
	"sync"
)

var (
	tingService     *TingService
	tingServiceOnce sync.Once
)

func GetTingService() *TingService {
	tingServiceOnce.Do(func() {
		tingService = &TingService{}
	})
	return tingService
}

type TingService struct {
}

func (s *TingService) parseTabList(doc *goquery.Selection) (list []*dto.TabItem, err error) {
	list = []*dto.TabItem{}
	doc.Each(func(i int, aItem *goquery.Selection) {
		href, _ := aItem.Find("a").Attr("href")
		href = strings.TrimSuffix(href, ".html")
		href = strings.TrimPrefix(href, ripple.GetConfig().GetString("baseUrl.ting"))
		if !strings.HasPrefix(href, "/book") {
			return
		}
		list = append(list, &dto.TabItem{
			Name: aItem.Text(),
			Key:  href,
		})
	})
	return
}

func (s *TingService) GetTabList() (list []*dto.TabItem, err error) {
	list = []*dto.TabItem{}
	var resp *resty.Response
	resp, err = resty.New().R().Get(ripple.GetConfig().GetString("baseUrl.ting"))
	if err != nil {
		return
	}
	htmlStr := string(resp.Body())
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlStr))
	if err != nil {
		fmt.Printf("Failed to parse HTML: %v", err)
	}
	list, err = s.parseTabList(doc.Find(".nav").Find(".clearfix").Find("li"))
	return
}

func (s *TingService) GetBookListByTab(tab string, index int64) (bookList []*dto.BookItem, totalCount int64, err error) {
	var resp *resty.Response
	resp, err = resty.New().R().Get(fmt.Sprintf(`%s%s-%d.html`, ripple.GetConfig().GetString("baseUrl.ting"), tab, index))
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
func (s *TingService) getBookListFromHtml(htmlStr string) (bookList []*dto.BookItem, totalCount int64) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlStr))
	if err != nil {
		fmt.Printf("Failed to parse HTML: %v", err)
	}
	tabInfo := &dto.TabItem{
		Name: "",
		Key:  "",
	}
	doc.Find(".place").Find("a").Each(func(i int, item *goquery.Selection) {
		href, _ := item.Attr("href")
		if strings.HasPrefix(href, "/book/") {
			tabInfo.Name = item.Text()
			tabInfo.Key = strings.TrimSuffix(href, ".html")
		}
	})
	doc.Find(".row3").Find(".style-img").Each(func(i int, bItem *goquery.Selection) {
		bookInfo := s.parseBookInfo(bItem)
		bookInfo.Tab = tabInfo.Key
		bookInfo.TabName = tabInfo.Name
		bookList = append(bookList, bookInfo)
	})
	lastTab, _ := doc.Find(".pagebar").Find("a").Last().Attr("href")
	lastTab = strings.TrimPrefix(lastTab, fmt.Sprintf("%s-", tabInfo.Key))
	lastTab = strings.TrimSuffix(lastTab, ".html")
	totalCount = cast.ToInt64(lastTab)
	return
}
func (s *TingService) parseHref(href string) (tab, bookId string) {
	bookId = strings.TrimPrefix(href, "/show/")
	bookId = strings.TrimSuffix(href, ".html")
	return
}
func (s *TingService) parseIndexList(doc *goquery.Selection) (bookList []*dto.BookItem) {
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
func (s *TingService) GetBookIndexList() (tabList []*dto.TabItem, recommendList []*dto.BookItem, newList []*dto.BookItem, err error) {
	var resp *resty.Response
	resp, err = resty.New().R().Get(fmt.Sprintf(`%s%s`, ripple.GetConfig().GetString("baseUrl.ting"), "/yousheng/"))
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
func (s *TingService) parseBookInfo(doc *goquery.Selection) (bookInfo *dto.BookItem) {
	speakerMatch, _ := regexp.Compile(`由.*播音`)
	bookInfo = &dto.BookItem{}
	bookInfo.Pic, _ = doc.Find("img").Attr("src")
	bookInfo.Pic = fmt.Sprintf("%s%s", ripple.GetConfig().GetString("baseUrl.ting"), bookInfo.Pic)
	nItem := doc.Find("section")
	bookInfo.Author = nItem.Find(".fr").Text()
	href, _ := nItem.Find("a").Attr("href")
	if strings.HasPrefix(href, "/show/") {
		bookInfo.BookId = strings.TrimPrefix(href, "/show/")
		bookInfo.BookId = strings.TrimSuffix(bookInfo.BookId, ".html")
	}
	bookInfo.Name = nItem.Find(".f-bold").Text()
	bookInfo.Speaker = speakerMatch.FindString(doc.Find(".f-gray").Text())
	bookInfo.Speaker = strings.TrimPrefix(bookInfo.Speaker, "由")
	bookInfo.Speaker = strings.TrimSuffix(bookInfo.Speaker, "播音")
	return
}
func (s *TingService) GetBookInfo(tab, bookId string) (bookInfo *dto.BookItem, chapterList []*dto.ChapterItem, chapterCount int64, err error) {
	var resp *resty.Response
	resp, err = resty.New().R().Get(fmt.Sprintf(`%s/show/%s.html`, ripple.GetConfig().GetString("baseUrl.ting"), bookId))
	if err != nil {
		return
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(resp.Body())))
	if err != nil {
		fmt.Printf("Failed to parse HTML: %v", err)
	}
	tabInfo := &dto.TabItem{
		Name: "",
		Key:  tab,
	}
	doc.Find(".place").Find("a").Each(func(i int, item *goquery.Selection) {
		href, _ := item.Attr("href")
		if strings.HasPrefix(href, "/book/") {
			tabInfo.Name = item.Text()
			tabInfo.Key = strings.TrimSuffix(href, ".html")
		}
	})
	bookInfo = s.parseBookInfo(doc.Find(".style-img").First())
	bookInfo.BookId = bookId
	bookInfo.Tab = tabInfo.Key
	bookInfo.TabName = tabInfo.Name
	doc.Find(".ul-36").Find("li").Find("a").Each(func(i int, item *goquery.Selection) {
		title, _ := item.Attr("title")
		href, _ := item.Attr("href")
		href = strings.TrimPrefix(href, fmt.Sprintf("/jpplay/%s-", bookId))
		href = strings.TrimSuffix(href, ".html")
		chapterList = append(chapterList, &dto.ChapterItem{
			Name:      title,
			ChapterId: href,
		})
	})
	chapterCount = int64(len(chapterList))
	if chapterCount >= 10 {
		chapterList = chapterList[0:10]
	}
	return
}
func (s *TingService) BookPlay(ctx echo.Context, bookId string, chapter string) (err error) {
	var resp *resty.Response
	resp, err = resty.New().R().SetHeader("content-type", "application/x-www-form-urlencoded; charset=UTF-8").Get(fmt.Sprintf(`%s/jpplay/%s-%s.html`, ripple.GetConfig().GetString("baseUrl.ting"), bookId, chapter))
	if err != nil {
		return
	}
	if err != nil {
		return
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(resp.Body())))
	if err != nil {
		return
	}
	var urlStr string
	strArr := strings.Split(doc.Find(".tx-text").Find("script").First().Text(), ";")
	for _, item := range strArr {
		if strings.HasPrefix(item, "var now") {
			urlStr = strings.TrimPrefix(item, "var now=\"")
			urlStr = strings.TrimSuffix(urlStr, "\"")
		}
	}
	resp, err = resty.New().R().Get(urlStr)
	_, err = ctx.Response().Write(resp.Body())
	return
}
func (s *TingService) Search(keyword string, index int64) (bookList []*dto.BookItem, totalCount int64, err error) {
	var resp *resty.Response
	resp, err = resty.New().R().SetFormData(map[string]string{
		"show":     "title,newstext,player,playadmin",
		"keyboard": keyword,
	}).Post(fmt.Sprintf(`%s/e/search/index.php`, ripple.GetConfig().GetString("baseUrl.ting")))
	if err != nil {
		return
	}

	fmt.Printf(string(resp.Body()))
	//bookList, totalCount = s.getBookListFromHtml(string(resp.Body()))
	return
}
