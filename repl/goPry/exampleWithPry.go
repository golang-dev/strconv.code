package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/d4l3k/go-pry/pry"
)

func main() {
	args := os.Args
	for i := 1; i < len(args); i++ {
		ret, err := strconv.Atoi(args[i])
		pry.Pry()
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Println(ret)
		}
	}
}
