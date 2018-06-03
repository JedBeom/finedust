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
	"time"

	"github.com/go-vgo/robotgo"
)

// 본격적인 데이터가 들어가는 구조체
type item struct {
	Time        string `xml:"dataTime"`
	Pm10Value   int    `xml:"pm10Value"`
	Pm25Value   int    `xml:"pm25Value"`
	Pm10Rate    int
	Pm25Rate    int
	MixedRate   int
	HangulRate  string
	Background  string
	stationName string
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

func thisTime(stationName string) response {
	var full response
	full = response{Body: body{Item: item{stationName: stationName}}}

	link := fmt.Sprintf("http://openapi.airkorea.or.kr/openapi/services/rest/ArpltnInforInqireSvc/getMsrstnAcctoRltmMesureDnsty?serviceKey=OOtkvfDic1VY%%2FlqF%%2Fwf57rsYRL8j5a7zXlqNVby7h9SKOo4Vf0khrnDceMU3%%2FAfnSGxxTAqYF41jf8zb%%2BkuHoQ%%3D%%3D&numOfRows=1&pageSize=1&pageNo=1&startPage=1&stationName=%v&dataTerm=DAILY&ver=1.3", stationName)

	resp, err := http.Get(link)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	err = xml.Unmarshal(data, &full)
	if err != nil {
		log.Printf("I've got an error while Unmarshal XML file. stationName: %v", stationName)
		if stationName != "장천동" {
			full = thisTime("장천동")
		} else {
			full = thisTime("연향동")
		}
	}
	return full
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
	full := thisTime("연향동")
	Item := full.Body.Item
	Item = rater(Item)
	Item = MixingRatesAndGiveAHangulRate(Item)
	t, err := template.ParseFiles("index.html")
	if err != nil {
		log.Println("Error from template:", err)
	}
	log.Println(r.Header.Get("X-FORWARDED-FOR"))
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
	var color string

	switch Item.MixedRate {
	case 1:
		rate = "최고"
		color = "#4b74b8"
	case 2:
		rate = "좋음"
		color = "#5c98d1"
	case 3:
		rate = "양호"
		color = "#63abc0"
	case 4:
		rate = "보통"
		color = "#5b8f4a"
	case 5:
		rate = "나쁨"
		color = "#de8b2f"
	case 6:
		rate = "상당히 나쁨"
		color = "#c94e2c"
	case 7:
		rate = "매우 나쁨"
		color = "#b83134"
	case 8:
		rate = "최악"
		color = "#262626"
	}
	Item.HangulRate = rate
	Item.Background = color
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
	time.Sleep(1 * time.Second)
	go robotgo.KeyTap("f11")
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
