package main

import (
	"DayStorySrv/parser"
	"DayStorySrv/tools"
	"fmt"
	"net/http"
	"strings"
)

func main() {

	http.HandleFunc("/", FuncProvider)
	if err := http.ListenAndServe(":8089", nil); err != nil {
		panic(err)
	}

}

func Compile(date string) string {
	news := parser.Parse("https://ru.wikipedia.org/wiki/" + date + "#События")

	worldHoliday := "мировые праздники: " + tools.GetRandomElem(news.WorldHolidays)
	localHoliday := "локальные праздники: " + tools.GetRandomElem(news.LocalHolidays)
	event, links := tools.GetRandomKV(news.Events)
	links = tools.LinksToImages(links)

	link := tools.GetRandomElem(links)

	result := fmt.Sprintf("сегодня %s \n %s \n %s \n события: %s \n %s",
		news.Day,
		worldHoliday,
		localHoliday,
		event,
		link)

	return result
}

func FuncProvider(writer http.ResponseWriter, rq *http.Request) {
	rqMsg := strings.TrimPrefix(rq.URL.Path, "/")
	err := rq.ParseForm()
	if err != nil {
		panic(err)
	}
	args := rq.Form

	if rqMsg == "GetRandomDayInfo" {
		//help: ...:8089/GetRandomDayInfo
		date := tools.GetFormatedDate()
		rqMsg = Compile(date)
	}
	if rqMsg == "GetDayInfo" {
		//help: ...:8089/GetDayInfo?date=15_мая
		date := tools.LstToStr(args["date"])
		rqMsg = Compile(date)
	}

	if _, err := writer.Write([]byte(rqMsg)); err != nil {
		panic(err)
	}
}
