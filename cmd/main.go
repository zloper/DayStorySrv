package main

import (
	"DayStorySrv/parser"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

func main() {

	http.HandleFunc("/", FuncProvider)
	if err := http.ListenAndServe(":8089", nil); err != nil {
		panic(err)
	}

}

func Compile(date string) string {
	url := "https://ru.wikipedia.org/wiki/" + date + "#События"
	news := parser.Parse(url)
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)
	fmt.Println(random.Intn(10), random.Intn(10))
	fmt.Println(news.Day)
	world := 0
	local := 0
	if len(news.WorldHolidays) > 0 {
		world = random.Intn(len(news.WorldHolidays))
	}
	if len(news.LocalHolidays) > 0 {
		local = random.Intn(len(news.LocalHolidays))
	}

	keys := make([]string, 0, len(news.Events))
	for k := range news.Events {
		keys = append(keys, k)
	}
	num := random.Intn(len(news.Events))
	randomKey := keys[num]
	lst := news.Events[randomKey]

	// switch pageUrl <=> imgUrl
	var linksList []string
	for index := range lst {
		img := "https:" + parser.GetImage(lst[index])
		if strings.HasSuffix(img, ".svg.png") {
			fmt.Println("bad image")
		} else {
			linksList = append(linksList, img)
		}
	}

	worldRes := "мировые праздники: "
	localRes := "локальные праздники: "
	if len(news.WorldHolidays) > 0 {
		fmt.Println("международние праздники:", news.WorldHolidays[world])
		worldRes += news.WorldHolidays[world]
	} else {
		worldRes += "None"
	}
	if len(news.LocalHolidays) > 0 {
		fmt.Println("локальные праздники:", news.LocalHolidays[local])
		localRes += news.LocalHolidays[local]
	} else {
		localRes += "None"
	}
	fmt.Println(randomKey, linksList)
	//TODO refactor return

	numLinks := random.Intn(len(linksList))
	return "сегодня " + news.Day + "\n" + worldRes + "\n" + localRes + "\n" + "события: " + randomKey + "\n" + linksList[numLinks]
}

func FuncProvider(writer http.ResponseWriter, rq *http.Request) {
	rqMsg := strings.TrimPrefix(rq.URL.Path, "/")
	err := rq.ParseForm()
	if err != nil {
		panic(err)
	}

	if rqMsg == "GetRandomDayInfo" {
		//help: ...:8089/GetRandomDayInfo?date=15maya
		date := GetCurrentDate()
		rqMsg = Compile(date)
	}

	if _, err := writer.Write([]byte(rqMsg)); err != nil {
		panic(err)
	}
}

func GetCurrentDate() string {
	dt := time.Now()
	mounth := dt.Format("01")
	dct := map[string]string{ // map literal
		"01": "января",
		"02": "февраля",
		"03": "марта",
		"04": "апреля",
		"05": "мая",
		"06": "июня",
		"07": "июля",
		"08": "августа",
		"09": "сентября",
		"10": "октября",
		"11": "ноября",
		"12": "декабря",
	}

	day := dt.Format("02")
	fmt.Println("Current date and time is: ", dct[mounth], day)
	return day + "_" + dct[mounth]
}
