package main

import (
	"encoding/xml"
	"fmt"
)

type item struct {
	Time      string `xml:"dataTime"`
	Pm10Value int    `xml:"pm10Value"`
	Pm25Value int    `xml:"pm25Value"`
	Pm10Rate  int
	Pm25Rate  int
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

func main() {
	Response := response{
		Body{
			Item{
				Time:      "",
				Pm10Value: 34,
				Pm25Value: 34,
				Pm10Rate:  0,
				Pm25Rate:  0}}}
	rating(&fine)
	fmt.Println(fine.Body.Item.Pm10Rate)
	fmt.Println(fine.Body.Item.Pm25Rate)
}
func rating(f *response) {
	pm10 := f.Body.Item.Pm10Value
	switch {
	case pm10 <= 15:
		*f.Body.Item.Pm10Rate = 1
	case pm10 <= 30:
		*f.Body.Item.Pm10Rate = 2
	case pm10 <= 40:
		*f.Body.Item.Pm10Rate = 3
	case pm10 <= 50:
		*f.Body.Item.Pm10Rate = 4
	case pm10 <= 75:
		*f.Body.Item.Pm10Rate = 5
	case pm10 <= 100:
		*f.Body.Item.Pm10Rate = 6
	case pm10 <= 150:
		*f.Body.Item.Pm10Rate = 7
	default:
		*f.Body.Item.Pm10Rate = 8
	}

	pm25 := *f.Body.Item.Pm25Value
	switch {
	case pm25 <= 8:
		*f.Body.Item.Pm25Value = 1
	case pm25 <= 15:
		*f.Body.Item.Pm25Value = 2
	case pm25 <= 20:
		*f.Body.Item.Pm25Value = 3
	case pm25 <= 25:
		*f.Body.Item.Pm25Value = 4
	case pm25 <= 37:
		*f.Body.Item.Pm25Value = 5
	case pm25 <= 50:
		*f.Body.Item.Pm25Value = 6
	case pm25 <= 75:
		*f.Body.Item.Pm25Value = 7
	default:
		*f.Body.Item.Pm25Value = 8
	}
}
