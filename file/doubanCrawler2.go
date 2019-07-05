package main

import (
	"os"
	"log"
	"time"
	"strings"
	"strconv"
	"net/http"
	"encoding/csv"

	"golang.org/x/net/html"
	"github.com/antchfx/htmlquery"
)

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func fetch(url string) *html.Node {
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
	doc, err := htmlquery.Parse(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return doc
}

func parseUrls(url string, ch chan bool, w *csv.Writer) {
	doc := fetch(url)
	nodes := htmlquery.Find(doc, `//ol[@class="grid_view"]/li//div[@class="hd"]`)

	for _, node := range nodes {
		url := htmlquery.FindOne(node, "./a/@href")
		title := htmlquery.FindOne(node, `.//span[@class="title"]/text()`)

		err := w.Write([]string{
			strings.Split(htmlquery.InnerText(url), "/")[4],
			htmlquery.InnerText(title)})
		checkError(err)
	}

	time.Sleep(2 * time.Second)
	ch <- true
}

func main() {
	start := time.Now()
	ch := make(chan bool)
	f, err := os.Create("movie.csv")
    checkError(err)
	defer f.Close()

	writer := csv.NewWriter(f)
    defer writer.Flush()

	err = writer.Write([]string{"ID", "Title"})
	checkError(err)

	for i := 0; i < 10; i++ {
		go parseUrls("https://movie.douban.com/top250?start="+strconv.Itoa(25*i), ch, writer)
	}

	for i := 0; i < 10; i++ {
		<-ch
	}

	elapsed := time.Since(start)
	log.Printf("Took %s", elapsed)
}
