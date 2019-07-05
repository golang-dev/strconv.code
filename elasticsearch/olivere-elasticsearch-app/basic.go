package main

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

	"github.com/olivere/elastic"
)

const mapping = `
{
	"mappings": {
		"online": {
		"properties": {
			"id": {
				"type": "long"
			},
			"title": {
				"type": "text"
			},
			"genres": {
				"type": "keyword"
			}
		}}
	}
}`

var (
	subject   Subject
	indexName = "subject"
	typeName = "online"
	servers   = []string{"http://localhost:9200/"}
)

type Subject struct {
	ID     int      `json:"id"`
	Title  string   `json:"title"`
	Genres []string `json:"genres"`
}

func Search(client *elastic.Client, ctx context.Context, genre string) {
	fmt.Printf("Search: %s", genre)
	// Term搜索
	termQuery := elastic.NewTermQuery("genres", genre)
	searchResult, err := client.Search().
		Index(indexName).
		Type(typeName).
		Query(termQuery).
		Sort("id", true). // 按id升序排序
		From(0).Size(10). // 拿前10个结果
		Pretty(true).
		Do(ctx) // 执行
	if err != nil {
		panic(err)
	}
	total := searchResult.TotalHits()
	fmt.Printf("Found %d subjects\n", total)
	if total > 0 {
		for _, item := range searchResult.Each(reflect.TypeOf(subject)) {
			if t, ok := item.(Subject); ok {
				fmt.Printf("Found: Subject(id=%d, title=%s)\n", t.ID, t.Title)
			}
		}

	} else {
		fmt.Println("Not found!")
	}
}

func main() {
	ctx := context.Background()
	client, err := elastic.NewClient(elastic.SetURL(servers...))
	if err != nil {
		panic(err)
	}

	// 用IndexExists检查索引是否存在
	exists, err := client.IndexExists(indexName).Do(ctx)
	if err != nil {
		panic(err)
	}
	if !exists {
		// 用CreateIndex创建索引，mapping内容用BodyString传入
		_, err := client.CreateIndex(indexName).BodyString(mapping).Do(ctx)
		if err != nil {
			panic(err)
		}
	}
	subject = Subject{
		ID:     1,
		Title:  "肖恩克的救赎",
		Genres: []string{"犯罪", "剧情"},
	}

	// 写入
	doc, err := client.Index().
		Index(indexName).
		Type(typeName).
		Id(strconv.Itoa(subject.ID)).
		BodyJson(subject).
		Refresh("wait_for").
		Do(ctx)

	if err != nil {
		panic(err)
	}
	fmt.Printf("Indexed with id=%v, type=%s\n", doc.Id, doc.Type)
	subject = Subject{
		ID:     2,
		Title:  "千与千寻",
		Genres: []string{"剧情", "喜剧", "爱情", "战争"},
	}
	fmt.Println(string(subject.ID))
	doc, err = client.Index().
		Index(indexName).
		Type(typeName).
		Id(strconv.Itoa(subject.ID)).
		BodyJson(subject).
		Refresh("wait_for").
		Do(ctx)

	if err != nil {
		panic(err)
	}

	// 获取
	result, err := client.Get().
		Index(indexName).
		Type(typeName).
		Id(strconv.Itoa(subject.ID)).
		Do(ctx)
	if err != nil {
		panic(err)
	}
	if result.Found {
		fmt.Printf("Got document %v (version=%d, index=%s, type=%s)\n",
			result.Id, result.Version, result.Index, result.Type)
		err := json.Unmarshal(*result.Source, &subject)
		if err != nil {
			panic(err)
		}
		fmt.Println(subject.ID, subject.Title, subject.Genres)
	}

	// 搜索
	Search(client, ctx, "剧情")
	fmt.Println("****")
	Search(client, ctx, "犯罪")

	// 删除
	res, err := client.Delete().
		Index(indexName).
		Type(typeName).
		Id("1").
		Refresh("wait_for").
		Do(ctx)
	if err != nil {
		panic(err)
	}
	if res.Result == "deleted" {
		fmt.Println("Document 1: deleted")
	}
	fmt.Println("****")
	Search(client, ctx, "犯罪")
}
