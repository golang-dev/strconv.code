package main

import (
	"encoding/json"
	"fmt"
)

type Person struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func main() {
	bytes, _ := json.Marshal(&Person{"Dany", "Boon"})
	fmt.Printf("%s", bytes)
}
