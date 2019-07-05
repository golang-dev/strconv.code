package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	n int
    h bool
	q *bool
    s string
)

func init() {
	q = flag.Bool("q", false, "Exit")
	flag.BoolVar(&h, "h", false, "Show help")
	flag.IntVar(&n, "n", 0, "Set number")
	flag.StringVar(&s, "s", "Default string", "Set String")
}

func main() {
    flag.Parse()

    if h {
        flag.Usage()
    } else {
		if *q {
			fmt.Println("q is ", *q)
			os.Exit(0)
		}
		fmt.Println("Number is ", n)
		fmt.Println("String is ", s)
	}
}
