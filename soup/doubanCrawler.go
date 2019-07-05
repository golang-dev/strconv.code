package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/anaskhan96/soup"
)

func fetch(url string) soup.Root {
	fmt.Println("Fetch Url", url)
	soup.Headers = map[string]string{
		"User-Agent": "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)",
	}

	source, err := soup.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	doc := soup.HTMLParse(source)
	return doc
}

func parseUrls(url string, ch chan bool) {
	doc := fetch(url)
	for _, root := range doc.Find("ol", "class", "grid_view").FindAll("div", "class", "hd") {
		movieUrl, _ := root.Find("a").Attrs()["href"]
		title := root.Find("span", "class", "title").Text()
		fmt.Println(strings.Split(movieUrl, "/")[4], title)
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
	fmt.Printf("Took %s", elapsed)
}
