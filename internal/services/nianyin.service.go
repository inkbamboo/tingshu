package services

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/bmbstack/ripple"
	"github.com/go-resty/resty/v2"
	"github.com/inkbamboo/tingshu/dto"
	"github.com/labstack/echo/v4"
	"github.com/tidwall/gjson"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	nianYinService     *NianYinService
	nianYinServiceOnce sync.Once
)

func GetNianYinService() *NianYinService {
	nianYinServiceOnce.Do(func() {
		nianYinService = &NianYinService{}
	})
	return nianYinService
}

type NianYinService struct {
}

func (s *NianYinService) parseTabList(doc *goquery.Selection) (list []*dto.TabItem, err error) {
	list = []*dto.TabItem{}
	doc.Each(func(i int, aItem *goquery.Selection) {
		href, _ := aItem.Attr("href")
		if len(href) <= 2 || !strings.HasSuffix(href, "/") || !strings.HasPrefix(href, "/") {
			return
		}
		list = append(list, &dto.TabItem{
			Name: aItem.Text(),
			Key:  href,
		})
	})
	return
}
func (s *NianYinService) GetTabList() (list []*dto.TabItem, err error) {
	list = []*dto.TabItem{}
	var resp *resty.Response
	resp, err = resty.New().R().Get(fmt.Sprintf(`%s`, ripple.GetConfig().GetString("baseUrl.nianYin")))
	if err != nil {
		return
	}
	htmlStr := string(resp.Body())
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlStr))
	if err != nil {
		fmt.Printf("Failed to parse HTML: %v", err)
	}
	list, err = s.parseTabList(doc.Find(".nav").Find("a"))
	return
}

func (s *NianYinService) GetBookListByTab(tab string, index int64) (bookList []*dto.BookItem, totalCount int64, err error) {
	var resp *resty.Response
	resp, err = resty.New().R().Get(fmt.Sprintf(`%s%sindex%d.html`, ripple.GetConfig().GetString("baseUrl.nianYin"), tab, index))
	if err != nil {
		return
	}
	bookList, totalCount = s.getBookListFromHtml(string(resp.Body()))
	return
}
func (s *NianYinService) getBookListFromHtml(htmlStr string) (bookList []*dto.BookItem, totalCount int64) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlStr))
	if err != nil {
		fmt.Printf("Failed to parse HTML: %v", err)
	}
	doc.Find(".clist").Find("a").Each(func(i int, aItem *goquery.Selection) {
		bookInfo := &dto.BookItem{}
		bookInfo.Pic, _ = aItem.Find("img").Attr("src")
		href, _ := aItem.Attr("href")
		bookInfo.Tab, bookInfo.BookId = s.parseHref(href)
		aItem.Find("dd").Children().Each(func(i int, dItem *goquery.Selection) {
			switch i {
			case 0:
				bookInfo.Name = dItem.Text()
			case 1:
				bookInfo.TabName = strings.Split(dItem.Text(), "：")[1]
			case 2:
				bookInfo.Author = strings.Split(dItem.Text(), "：")[1]
			case 3:
				bookInfo.Speaker = strings.Split(dItem.Text(), "：")[1]
			case 4:
				bookInfo.CreateTime = strings.Split(dItem.Text(), "：")[1]
			}
		})
		bookList = append(bookList, bookInfo)
	})
	reg, _ := regexp.Compile("\\d*/\\d*")
	doc.Find(".cpage").Find("span").Each(func(i int, selection *goquery.Selection) {
		txt := reg.FindString(selection.Text())
		if txt != "" {
			txt = strings.Split(txt, "/")[1]
			totalCount, _ = strconv.ParseInt(txt, 10, 64)
		}
	})
	return
}

func (s *NianYinService) parseIndexList(doc *goquery.Selection) (bookList []*dto.BookItem) {
	doc.Each(func(i int, aItem *goquery.Selection) {
		bookInfo := &dto.BookItem{}
		bookInfo.Pic, _ = aItem.Find("img").Attr("src")
		href, _ := aItem.Attr("href")
		bookInfo.Tab, bookInfo.BookId = s.parseHref(href)
		aItem.Find("dd").Children().Each(func(i int, dItem *goquery.Selection) {
			switch i {
			case 0:
				bookInfo.Name = dItem.Text()
			case 1:
				bookInfo.Summary = dItem.Text()
			}
		})
		bookList = append(bookList, bookInfo)
	})
	return
}
func (s *NianYinService) parseHref(href string) (tab, bookId string) {
	hrefArr := strings.Split(href, "/")
	tab = strings.Join(hrefArr[0:len(hrefArr)-1], "/") + "/"
	bookId = hrefArr[len(hrefArr)-1]
	bookId = strings.Split(bookId, ".")[0]
	return
}
func (s *NianYinService) GetBookIndexList() (tabList []*dto.TabItem, recommendList []*dto.BookItem, newList []*dto.BookItem, err error) {
	var resp *resty.Response
	resp, err = resty.New().R().Get(fmt.Sprintf(`%s`, ripple.GetConfig().GetString("baseUrl.nianYin")))
	if err != nil {
		return
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(resp.Body())))
	if err != nil {
		fmt.Printf("Failed to parse HTML: %v", err)
	}
	recommendList = s.parseIndexList(doc.Find(".tlist").Find("a"))
	newList = s.parseIndexList(doc.Find(".nlist").Find("a"))
	tabList, err = s.parseTabList(doc.Find(".nav").Find("a"))
	return
}
func (s *NianYinService) GetBookInfo(tab, bookId string) (bookInfo *dto.BookItem, chapterList []*dto.ChapterItem, chapterCount int64, err error) {
	var resp *resty.Response
	resp, err = resty.New().R().Get(fmt.Sprintf(`%s%s%s.html`, ripple.GetConfig().GetString("baseUrl.nianYin"), tab, bookId))
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
	doc.Find(".binfo").Children().Each(func(i int, item *goquery.Selection) {
		switch i {
		case 0:
			bookInfo.Name = item.Text()
		case 1:
			bookInfo.TabName = strings.Split(item.Text(), "：")[1]
		case 2:
			bookInfo.Author = strings.Split(item.Text(), "：")[1]
		case 3:
			bookInfo.Speaker = strings.Split(item.Text(), "：")[1]
		case 4:
			bookInfo.Status = strings.Split(item.Text(), "：")[1]
		case 5:
			bookInfo.CreateTime = strings.Split(item.Text(), "：")[1]
		}
	})
	bookInfo.Pic, _ = doc.Find(".book").Find("img").Attr("src")
	bookInfo.Summary = doc.Find(".intro").Find("p").Text()
	bookInfo.Summary = strings.TrimSuffix(bookInfo.Summary, "法律声明：点击查看")
	doc.Find(".plist").Find("a").Each(func(i int, item *goquery.Selection) {
		chapterList = append(chapterList, &dto.ChapterItem{
			Name:      item.Text(),
			ChapterId: i + 1,
		})
	})
	chapterCount = int64(len(chapterList))
	//if chapterCount >= 10 {
	//	chapterList = chapterList[0:10]
	//}
	return
}
func (s *NianYinService) BookPlay(ctx echo.Context, bookId string, chapter int64) (err error) {
	var resp *resty.Response
	resp, err = resty.New().R().SetHeader("content-type", "application/x-www-form-urlencoded; charset=UTF-8").SetFormData(map[string]string{
		"bookId": bookId,
		"page":   fmt.Sprintf("%d", chapter),
	}).Post(fmt.Sprintf(`%s/?s=api-getneoplay`, ripple.GetConfig().GetString("baseUrl.nianYin")))
	if err != nil {
		return
	}
	urlStr := gjson.Get(string(resp.Body()), "url").String()
	urlStr = fmt.Sprintf("%s?v=%d", urlStr, time.Now().UnixMilli())
	resp, err = resty.New().R().SetHeaders(map[string]string{
		"Referer": "https://m.nianyin.com/",
	}).Get(urlStr)
	_, err = ctx.Response().Write(resp.Body())
	return
}
func (s *NianYinService) Search(keyword string, index int64) (bookList []*dto.BookItem, totalCount int64, err error) {
	var resp *resty.Response
	resp, err = resty.New().R().Get(fmt.Sprintf(`%s/?s=ting-search-wd-%s-p-%d.html`, ripple.GetConfig().GetString("baseUrl.nianYin"), keyword, index))
	if err != nil {
		return
	}
	bookList, totalCount = s.getBookListFromHtml(string(resp.Body()))
	return
}
