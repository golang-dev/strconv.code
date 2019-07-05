package main

import (
	"encoding/json"
	"fmt"
	"time"
)

type JSONTime struct {
	time.Time
}

func (t *JSONTime) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, t.Format("2006-01-02 15:04:05"))), nil
}

type Subject struct {
	ID        int `json:id`
	CreatedAt JSONTime
}

func main() {
	bytes, _ := json.Marshal(&Subject{111, JSONTime{time.Now()}})
	fmt.Printf("%s", bytes)
}
