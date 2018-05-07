package main

import (
	"encoding/xml"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
)

// 본격적인 데이터가 들어가는 구조체
type item struct {
	Time      string `xml:"dataTime"`
	Pm10Value int    `xml:"pm10Value"`
	Pm25Value int    `xml:"pm25Value"`
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
	var fine response

	resp, err := http.Get("http://openapi.airkorea.or.kr/openapi/services/rest/ArpltnInforInqireSvc/getMsrstnAcctoRltmMesureDnsty?serviceKey=OOtkvfDic1VY%2FlqF%2Fwf57rsYRL8j5a7zXlqNVby7h9SKOo4Vf0khrnDceMU3%2FAfnSGxxTAqYF41jf8zb%2BkuHoQ%3D%3D&numOfRows=1&pageSize=1&pageNo=1&startPage=1&stationName=연향동&dataTerm=DAILY&ver=1.3")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	err = xml.Unmarshal(data, &fine)
	if err != nil {
		panic(err)
	}
	return fine
}

func sender(w http.ResponseWriter, r *http.Request) {
	fine := thisTime()
	fmt.Println(fine.Body.Item)
	t, _ := template.ParseFiles("index.html")
	t.Execute(w, fine.Body.Item)
}

func main() {
	server := http.Server{
		Addr: ":8080",
	}
	http.HandleFunc("/", sender)
	server.ListenAndServe()
}
