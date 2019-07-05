package main

import (
	"log"
	"time"
	"strings"
	"strconv"
	"net/http"

	"github.com/lestrrat-go/libxml2"
	"github.com/lestrrat-go/libxml2/types"
	"github.com/lestrrat-go/libxml2/xpath"
)

func fetch(url string) types.Document {
	log.Println("Fetch Url", url)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Http get err:", err)
	}
	if resp.StatusCode != 200 {
		log.Fatal("Http status code:", resp.StatusCode)
	}
	defer resp.Body.Close()
	doc, err := libxml2.ParseHTMLReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return doc
}

func parseUrls(url string, ch chan bool) {
	doc := fetch(url)
	defer doc.Free()
	nodes := xpath.NodeList(doc.Find(`//ol[@class="grid_view"]/li//div[@class="hd"]`))
	for _, node := range nodes {
		urls, _ := node.Find("./a/@href")
		titles, _ := node.Find(`.//span[@class="title"]/text()`)
		log.Println(strings.Split(urls.NodeList()[0].TextContent(), "/")[4],
			titles.NodeList()[0].TextContent())
	}
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
	log.Printf("Took %s", elapsed)
}
