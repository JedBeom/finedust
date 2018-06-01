package main

import (
	"encoding/xml"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"runtime"
)

// 본격적인 데이터가 들어가는 구조체
type item struct {
	Time       string `xml:"dataTime"`
	Pm10Value  int    `xml:"pm10Value"`
	Pm25Value  int    `xml:"pm25Value"`
	Pm10Rate   int
	Pm25Rate   int
	MixedRate  int
	HangulRate string
}

/*
type items struct {
	Item []item `xml:"item"`
}
*/
type body struct {
	Item item `xml:"items>item"`
}

type response struct {
	XMLName xml.Name `xml:"response"`
	Body    body     `xml:"body"`
}

func thisTime() response {
	var Item response

	resp, err := http.Get("http://openapi.airkorea.or.kr/openapi/services/rest/ArpltnInforInqireSvc/getMsrstnAcctoRltmMesureDnsty?serviceKey=OOtkvfDic1VY%2FlqF%2Fwf57rsYRL8j5a7zXlqNVby7h9SKOo4Vf0khrnDceMU3%2FAfnSGxxTAqYF41jf8zb%2BkuHoQ%3D%3D&numOfRows=1&pageSize=1&pageNo=1&startPage=1&stationName=장천동&dataTerm=DAILY&ver=1.3")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	err = xml.Unmarshal(data, &Item)
	if err != nil {
		panic(err)
	}
	return Item
}

func rater(Item item) item {
	// 미세먼지 등급을 메기기(1-8)
	// 높을수록 공기의 상태가 좋음
	pm10 := Item.Pm10Value
	switch {
	case pm10 <= 15:
		Item.Pm10Rate = 1
	case pm10 <= 30:
		Item.Pm10Rate = 2
	case pm10 <= 40:
		Item.Pm10Rate = 3
	case pm10 <= 50:
		Item.Pm10Rate = 4
	case pm10 <= 75:
		Item.Pm10Rate = 5
	case pm10 <= 100:
		Item.Pm10Rate = 6
	case pm10 <= 150:
		Item.Pm10Rate = 7
	default:
		Item.Pm10Rate = 8
	}
	// 초미세먼지 등급을 메김(1-8)
	// 등급이 높을수록 공기의 상태가 좋음
	pm25 := Item.Pm25Value
	switch {
	case pm25 <= 8:
		Item.Pm25Rate = 1
	case pm25 <= 15:
		Item.Pm25Rate = 2
	case pm25 <= 20:
		Item.Pm25Rate = 3
	case pm25 <= 25:
		Item.Pm25Rate = 4
	case pm25 <= 37:
		Item.Pm25Rate = 5
	case pm25 <= 50:
		Item.Pm25Rate = 6
	case pm25 <= 75:
		Item.Pm25Rate = 7
	default:
		Item.Pm25Rate = 8
	}
	return Item
}

func sender(w http.ResponseWriter, r *http.Request) {
	full := thisTime()
	Item := full.Body.Item
	Item = rater(Item)
	Item = MixingRatesAndGiveAHangulRate(Item)
	t, err := template.ParseFiles("index.html")
	if err != nil {
		fmt.Println("Error from template:", err)
	}
	fmt.Println(Item)
	t.Execute(w, Item)
}

// open opens the specified URL in the default browser of the user.
func open(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}

func MixingRatesAndGiveAHangulRate(Item item) item {
	if Item.Pm10Rate >= Item.Pm25Rate {
		Item.MixedRate = Item.Pm10Rate
	} else if Item.Pm10Rate < Item.Pm25Rate {
		Item.MixedRate = Item.Pm25Rate
	} else {
		log.Println("Error on Logic.")
	}
	var rate string

	switch Item.MixedRate {
	case 1:
		rate = "최고"
	case 2:
		rate = "좋음"
	case 3:
		rate = "양호"
	case 4:
		rate = "보통"
	case 5:
		rate = "나쁨"
	case 6:
		rate = "상당히 나쁨"
	case 7:
		rate = "매우 나쁨"
	case 8:
		rate = "최악"
	}
	Item.HangulRate = rate
	return Item
}

func main() {
	port := ":8080"
	fmt.Println("Server Started at port", port)
	server := http.Server{
		Addr: port,
	}
	files := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", files))
	http.HandleFunc("/", sender)
	go open("http://localhost:8080")
	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("Error on ListenAndServe()")
	}
	/*
		full := thisTime()
		Item := full.Body.Item
		Item = rater(Item)
		fmt.Println(Item)
	*/
}
