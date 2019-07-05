package main

import (
	"encoding/json"
	"fmt"
)

type Person struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func (p *Person) MarshalJSON() ([]byte, error) {
	type Alias Person
	return json.Marshal(&struct {
		*Alias
		FullName string `json:"full_name"`
	}{
		Alias:    (*Alias)(p),
		FullName: fmt.Sprintf("%s %s", p.FirstName, p.LastName),
	})
}

func main() {
	bytes, _ := json.Marshal(&Person{"Dany", "Boon"})
	fmt.Printf("%s", bytes)
}
