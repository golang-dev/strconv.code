package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/olivere/elastic/v7"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

var (
	indexName = "subject"
	servers   = []string{"http://localhost:9200/"}
	client, _ = elastic.NewClient(elastic.SetURL(servers...))  
)

type Subject struct {
	ID     int      `json:"id"`
	Title  string   `json:"title"`
	Genres []string `json:"genres"`
}

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

func parseUrls(url string, ch chan bool) {
	doc := fetch(url)
	ctx := context.Background()
	nodes := htmlquery.Find(doc, `//ol[@class="grid_view"]/li`)

	subjects := []Subject{}

	for _, node := range nodes {
		url := htmlquery.FindOne(node, `.//div[@class="hd"]/a/@href`)
		title := htmlquery.FindOne(node, `.//span[@class="title"]/text()`)
		genre := htmlquery.Find(node, `.//div[@class="bd"]/p/text()`)

		id, _ := strconv.Atoi(strings.Split(htmlquery.InnerText(url), "/")[4])
		genreStr := strings.Split(htmlquery.InnerText(genre[1]), "/")[2]
		subject := Subject{id, htmlquery.InnerText(title),
			strings.Split(strings.TrimSpace(genreStr), " ")}
		subjects = append(subjects, subject)
	}
	
	bulkRequest := client.Bulk()
	for _, subject := range subjects {
		doc := elastic.NewBulkIndexRequest().Index(indexName).Id(strconv.Itoa(subject.ID)).Doc(subject)
		bulkRequest = bulkRequest.Add(doc)
	}

	response, err := bulkRequest.Do(ctx)
	if err != nil {
		panic(err)
	}
	failed := response.Failed()
	l := len(failed)
	if l > 0 {
		fmt.Printf("Error(%d)", l, response.Errors)
	}

	log.Println("Finished Url", url)

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
