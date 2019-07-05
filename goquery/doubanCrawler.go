package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func fetch(url string) *goquery.Document {
	fmt.Println("Fetch Url", url)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1)")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Http get err:", err)
	}
	if resp.StatusCode != 200 {
		log.Fatal("Http status code:", resp.StatusCode)
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return doc
}

func parseUrls(url string, ch chan bool) {
	doc := fetch(url)
	doc.Find("ol.grid_view li").Find(".hd").Each(func(index int, ele *goquery.Selection) {
		movieUrl, _ := ele.Find("a").Attr("href")
		fmt.Println(strings.Split(movieUrl, "/")[4], ele.Find(".title").Eq(0).Text())
	})
	time.Sleep(2 * time.Second)
	ch <- true
}

func main() {
	start := time.Now()
	ch := make(chan bool)
	for i := 0; i < 10; i++ {
		go parseUrls("https://movie.douban.com/top250?start="+strconv.Itoa(25*i), ch)
	}

	for i := 0; i < 10; i++ {
		<-ch
	}

	elapsed := time.Since(start)
	fmt.Printf("Took %s", elapsed)
}
