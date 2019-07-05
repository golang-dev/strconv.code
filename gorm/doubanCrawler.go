package main

import (
	"log"
	"time"
	"strings"
	"strconv"
	"net/http"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"golang.org/x/net/html"
	"github.com/antchfx/htmlquery"
)

type Movie struct {
	ID int `gorm:"AUTO_INCREMENT"`
	Title string `gorm:"type:varchar(100);unique_index"`
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

func parseUrls(url string, ch chan bool, db *gorm.DB) {
	doc := fetch(url)
	nodes := htmlquery.Find(doc, `//ol[@class="grid_view"]/li//div[@class="hd"]`)

	for _, node := range nodes {
		url := htmlquery.FindOne(node, "./a/@href")
		title := htmlquery.FindOne(node, `.//span[@class="title"]/text()`)

		id, _ := strconv.Atoi(strings.Split(htmlquery.InnerText(url), "/")[4])

		movie := &Movie{
			ID: id,
			Title: htmlquery.InnerText(title),
		}
		db.Create(&movie)
	}

	time.Sleep(2 * time.Second)
	ch <- true
}

func main() {
	start := time.Now()
	ch := make(chan bool)
	db, err := gorm.Open("mysql", "root:@/test?charset=utf8")
	defer db.Close()
    checkError(err)

	db.DropTableIfExists(&Movie{})
	db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&Movie{})

	for i := 0; i < 11; i++ {
		go parseUrls("https://movie.douban.com/top250?start="+strconv.Itoa(25*i), ch, db)
	}

	for i := 0; i < 11; i++ {
		<-ch
	}

	elapsed := time.Since(start)
	log.Printf("Took %s", elapsed)
}
