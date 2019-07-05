package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
	fmt.Println(time.Now().Format("2020-01-02 15:04:05"))

	t, _ := time.Parse("2006-01-02 15:04:05", "2006-01-02 03:04:05")

	fmt.Println(t)
	fmt.Println(t.Format("2006-01-02 03:04:05"))
	fmt.Println(t.Format("2006-01-02 15:04:05"))

	t, _ = time.Parse("2006-01-02 15:04:05", "2006-01-02 15:04:05")

	fmt.Println(t)
	fmt.Println(t.Format("2006-01-02 03:04:05"))
	fmt.Println(t.Format("2006-01-02 15:04:05"))
}
