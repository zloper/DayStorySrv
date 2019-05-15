package parser

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"strings"
)

type Page struct {
	Day           string
	WorldHolidays []string
	LocalHolidays []string
	Events        map[string][]string
	born          map[string]string
	dead          map[string]string
}

func Parse(url string) *Page {
	var dayInfo Page
	respose, _ := http.Get(url)
	defer respose.Body.Close()

	page, _ := goquery.NewDocumentFromReader(respose.Body)
	dayInfo.Day = page.Find("h1").Text()

	dayInfo.WorldHolidays = HolidaysParser(page, "span#Международные")
	dayInfo.LocalHolidays = HolidaysParser(page, "span#Национальные")

	fmt.Println(dayInfo.LocalHolidays)
	fmt.Println(dayInfo.WorldHolidays)

	EventsParser(page, "span#События")
	dayInfo.Events = EventsParser(page, "span#События")
	fmt.Println(dayInfo.Events)
	return &dayInfo
}

func HolidaysParser(pg *goquery.Document, id string) []string {
	var tmpList []string
	pg.Find(id).Parent().NextAllFiltered("ul").First().Find("li").Each(func(i int, selection *goquery.Selection) {
		tmpList = append(tmpList, selection.Text())
	})
	return tmpList
}

func EventsParser(pg *goquery.Document, id string) map[string][]string {
	var tmpList []*goquery.Selection

	pg.Find(id).Parent().NextAllFiltered("ul").Each(func(i int, ul *goquery.Selection) {
		isBlockEnded := ul.PrevAll().Find("span#Родились").Length() > 0
		if isBlockEnded {
			return
		}
		ul.Find("li").Each(func(i int, li *goquery.Selection) {
			tmpList = append(tmpList, li)
		})

	})
	result := make(map[string][]string)

	for index := range tmpList {
		name := tmpList[index].Text()
		linksList := FixLinks(tmpList[index])
		result[name] = linksList
	}

	return result
}

func FixLinks(selection *goquery.Selection) []string {
	mainUrl := "https://ru.wikipedia.org"
	var linksList []string
	//var img string
	selection.Find("a").Each(func(i int, linkObj *goquery.Selection) {
		endUrl, _ := linkObj.Attr("href")
		if strings.Contains(endUrl, "(") {
			newEnd := strings.Split(endUrl, "(")[1]
			newEnd = strings.TrimSuffix(newEnd, ")")
			endUrl = "/wiki/" + newEnd
		}
		// TODO add at final result (load all page is to long)
		//img = GetImage(mainUrl + endUrl)
		//linksList = append(linksList, img)
		linksList = append(linksList, mainUrl+endUrl)
	})
	return linksList
}

func GetImage(url string) string {
	awaitResponse := make(chan *http.Response)
	go func() { awaitResponse <- resp(url) }()
	response := <-awaitResponse
	defer response.Body.Close()
	page, _ := goquery.NewDocumentFromReader(response.Body)

	imgUrl, _ := page.Find("div#mw-content-text").Find("img").First().Attr("src")
	return imgUrl
}

func resp(url string) *http.Response {
	r, _ := http.Get(url)
	return r
}
