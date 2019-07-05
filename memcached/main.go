package main

import (
	"fmt"
	"log"

	"github.com/bradfitz/gomemcache/memcache"
)

var (
	server = "127.0.0.1:11211"
)

func main() {
	mc := memcache.New(server)
	if mc == nil {
		log.Fatal("New failed")
	}
	// Create
	err := mc.Set(&memcache.Item{Key: "a", Value: []byte("abcd")})
	if err != nil {
		log.Fatal("Set failed", err)
	}
	mc.Set(&memcache.Item{Key: "b", Value: []byte("123")})
	// Retrieve
	item, err := mc.Get("a")
	if err != nil { // err 为 ErrCacheMiss 表示没有这个键的缓存
		log.Printf("Key a 's value is %s", item.Value)
	}
	items, err := mc.GetMulti([]string{"a", "b", "c"})
	for k, v := range items {
		log.Printf("k=%s, v=%s", k, v.Value)
	}
	if err != nil {
		log.Print("GetMulti failed: %v", err)
	}
	// Delete
	err = mc.Delete("a")

	item, err = mc.Get("a")
	log.Printf("Get failed: %v, item=%v", err, item)

	err = mc.Delete("a")
	if err != nil {
		fmt.Println("Delete failed:", err)
	}
	// Incrby
	mc.Set(&memcache.Item{Key: "num", Value: []byte("1")})

	num, err := mc.Increment("num", 7) // 结果为1+7=8
	if err != nil {
		fmt.Println("Increment failed", err)
	} else {
		fmt.Println("The current value is :", num)
	}
	// Decrby
	num, err = mc.Decrement("num", 3) // 结果为8-3=5
	if err != nil {
		fmt.Println("Decrement failed", err)
	} else {
		fmt.Println("The current value is :", num)
	}
}
