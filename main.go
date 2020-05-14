package main

import (
	"encoding/xml"
	"fmt"
	ptime "github.com/yaa110/go-persian-calendar"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Item struct {
	XMLName xml.Name `xml:"item"`
	Id      int      `xml:"id"`
	Reg     string   `xml:"reg1"`
	Mag     string   `xml:"mag"`
	Dep     string   `xml:"dep"`
	Long    string   `xml:"long"`
	Lat     string   `xml:"lat"`
	Date    string   `xml:"date"`
}

type Items struct {
	XMLName xml.Name `xml:"items"`
	Items   []Item   `xml:"item"`
}

func getData() (string, error) {
	response, err := http.Get("http://irsc.ut.ac.ir/events_list.xml")
	if err != nil {
		return "", fmt.Errorf("GET error: %v", err)
	}

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("status error: %v", response.StatusCode)
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("read body: %v", err)
	}

	return string(data), nil
}

func main() {
	data, err := getData()
	if err != nil {
		log.Fatalln(err)
	}

	var events Items
	err = xml.Unmarshal([]byte(data), &events)
	if err != nil {
		log.Fatalln(err)
	}

	for i, item := range events.Items {
		if i == 0 {
			continue
		}
		if i > 10 {
			break
		}
		date, _ := time.Parse("2006-01-02 15:04:05", item.Date)
		pt := ptime.New(date.In(ptime.Iran()))
		message := fmt.Sprintf(
			"Region:%s\nDepth:%s\nTime:%s\nLocation:https://www.google.com/maps/search/?api=1&query=%s,%s\n",
			item.Reg,
			item.Dep,
			pt.Format("yyyy/MM/dd hh:mm:ss"),
			item.Lat,
			item.Long,
		)

		_, err := http.Get("https://api.telegram.org/bot1138407370:AAGcehBntpDFAD8fOsRiOf-iLOV3oV0ovJI/sendMessage?chat_id=@IranianEarthquakes&text=" + message)
		if err != nil {
			log.Fatalln(err)
		}
	}
}
