package main

import (
	"log"
	"strings"

	"github.com/gocolly/colly"
	"github.com/antchfx/htmlquery"
)

func main() {
	c := colly.NewCollector(
		colly.Async(true),
		colly.UserAgent("Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)"),
	)

	c.Limit(&colly.LimitRule{DomainGlob:  "*.douban.*", Parallelism: 5})

	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL)
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})

	c.OnResponse(func(r *colly.Response) {
		doc, err := htmlquery.Parse(strings.NewReader(string(r.Body)))
        if err != nil {
			log.Fatal(err)
        }
		nodes := htmlquery.Find(doc, `//ol[@class="grid_view"]/li//div[@class="hd"]`)
        for _, node := range nodes {
			url := htmlquery.FindOne(node, "./a/@href")
			title := htmlquery.FindOne(node, `.//span[@class="title"]/text()`)
			log.Println(strings.Split(htmlquery.InnerText(url), "/")[4],
				htmlquery.InnerText(title))
        }
	})

	c.OnHTML(".paginator a", func(e *colly.HTMLElement) {
		e.Request.Visit(e.Attr("href"))
	})

	c.Visit("https://movie.douban.com/top250?start=0&filter=")
	c.Wait()
}
