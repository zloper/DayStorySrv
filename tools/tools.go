package tools

import (
	"DayStorySrv/parser"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func GetFormatedDate() string {
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

func GetRandomElem(lst []string) string {
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)
	lenght := len(lst)
	result := ""
	if lenght > 0 {
		result = lst[random.Intn(lenght)]
	}
	return result
}

func GetRandomKV(mp map[string][]string) (string, []string) {
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)

	c := len(mp)
	keys := make([]string, 0, c)
	for k := range mp {
		keys = append(keys, k)
	}

	num := random.Intn(c)
	key := keys[num]

	return key, mp[key]
}

func LinksToImages(lst []string) []string {
	// switch pageUrl <=> imgUrl
	var imgList []string
	for index := range lst {
		img := "https:" + parser.GetImage(lst[index])
		if strings.HasSuffix(img, ".svg.png") {
			fmt.Println("bad image")
		} else {
			imgList = append(imgList, img)
		}
	}
	return imgList
}

func LstToStr(s []string) string {
	return strings.Join(s, " ")
}
