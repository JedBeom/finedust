package main

import (
	"encoding/xml"
	"fmt"
)

var xmlFile = []byte(`
<item>
    <dataTime>2018-05-03 21:00</dataTime>
    <mangName>도시대기</mangName>
    <so2Value>0.003</so2Value>
    <coValue>0.4</coValue>
    <o3Value>0.054</o3Value>
    <no2Value>0.005</no2Value>
    <pm10Value>31</pm10Value>
    <pm10Value24>-</pm10Value24>
    <pm25Value>16</pm25Value>
    <pm25Value24>-</pm25Value24>
    <khaiValue>-</khaiValue>
    <khaiGrade/>
    <so2Grade>1</so2Grade>
    <coGrade>1</coGrade>
    <o3Grade>2</o3Grade>
    <no2Grade>1</no2Grade>
    <pm10Grade/>
    <pm25Grade/>
    <pm10Grade1h>2</pm10Grade1h>
    <pm25Grade1h>2</pm25Grade1h>
</item>`)

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
	var fine item

	/*
		xp, err := os.Open("easiler.xml")
		if err != nil {
			panic(err)
		}
		defer xp.Close()

		data, err := ioutil.ReadAll(xp)
		if err != nil {
			panic(err)
		}
	*/

	err := xml.Unmarshal(xmlFile, &fine)
	if err != nil {
		panic(err)
	}
	fmt.Println(fine)
}
