package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
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
	Item []item `xml:"items>item"`
}

type response struct {
	XMLName xml.Name `xml:"response"`
	Body    body     `xml:"body"`
}

func main() {
	var fine response

	xp, err := os.Open("test.xml")
	if err != nil {
		panic(err)
	}
	defer xp.Close()

	data, err := ioutil.ReadAll(xp)
	if err != nil {
		panic(err)
	}

	err = xml.Unmarshal(data, &fine)
	if err != nil {
		panic(err)
	}
	fmt.Println(fine)
}
