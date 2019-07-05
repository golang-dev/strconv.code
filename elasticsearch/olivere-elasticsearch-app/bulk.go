package main

import (
	"context"
	"fmt"
	"strconv"
	"reflect"

	"github.com/olivere/elastic/v7"
)

var (
	indexName = "subject"
	servers   = []string{"http://localhost:9200/"}
)

type Subject struct {
	ID     int      `json:"id"`
	Title  string   `json:"title"`
	Genres []string `json:"genres"`
}

func main() {
	ctx := context.Background()
	client, err := elastic.NewClient(elastic.SetURL(servers...))
	if err != nil {
		panic(err)
	}

	subjects := []Subject{
		Subject{
			ID:     1,
			Title:  "肖恩克的救赎",
			Genres: []string{"犯罪", "剧情"},
		},
		Subject{
			ID:     2,
			Title:  "千与千寻",
			Genres: []string{"剧情", "喜剧", "爱情", "战争"},
		},
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

	subject3 := Subject{
		ID:     3,
		Title:  "这个杀手太冷",
		Genres: []string{"剧情", "动作", "犯罪"},
	}
	subject4 := Subject{
		ID:     4,
		Title:  "阿甘正传",
		Genres: []string{"剧情", "爱情"},
	}

	subject5 := subject3
	subject5.Title = "这个杀手不太冷"

	index1Req := elastic.NewBulkIndexRequest().Index(indexName).Id("3").Doc(subject3)
    index2Req := elastic.NewBulkIndexRequest().OpType("create").Index(indexName).Id("4").Doc(subject4)
	delete1Req := elastic.NewBulkDeleteRequest().Index(indexName).Id("1")
	update2Req := elastic.NewBulkUpdateRequest().Index(indexName).Id("3").
                Doc(subject5)

	bulkRequest = client.Bulk()
	bulkRequest = bulkRequest.Add(index1Req)
	bulkRequest = bulkRequest.Add(index2Req)
	bulkRequest = bulkRequest.Add(delete1Req)
	bulkRequest = bulkRequest.Add(update2Req)

	_, err = bulkRequest.Refresh("wait_for").Do(ctx)
	if err != nil {
		panic(err)
	}
  
  	if bulkRequest.NumberOfActions() == 0 {
		fmt.Println("Actions all clear!")
	}

	searchResult, err := client.Search().
		Index(indexName).
		Sort("id", false). // 按id升序排序
		Pretty(true).
		Do(ctx) // 执行
	if err != nil {
		panic(err)
	}
	var subject Subject
	for _, item := range searchResult.Each(reflect.TypeOf(subject)) {
		if t, ok := item.(Subject); ok {
			fmt.Printf("Found: Subject(id=%d, title=%s)\n", t.ID, t.Title)
		}
	}

}
