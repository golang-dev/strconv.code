package main

import (
	"encoding/json"
	"fmt"

	"github.com/olivere/elastic/v7"
)

func PrintQuery(src interface{}) {
	fmt.Println("*****")
	data, err := json.MarshalIndent(src, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))
}

func main() {
	query := elastic.NewTermQuery("genres", "动画")
	src, err := query.Source()
	if err != nil {
		panic(err)
	}
	PrintQuery(src)

	boolQuery := elastic.NewBoolQuery()
	boolQuery = boolQuery.Must(elastic.NewTermQuery("genres", "剧情"))
	boolQuery = boolQuery.Filter(elastic.NewTermQuery("id", 1))
	src, err = boolQuery.Source()
	if err != nil {
		panic(err)
	}
	PrintQuery(src)

	rangeQuery := elastic.NewRangeQuery("born").
		Gte("2012/01/01").
		Lte("now").
		Format("yyyy/MM/dd")
	src, err = rangeQuery.Source()
	PrintQuery(src)
}
