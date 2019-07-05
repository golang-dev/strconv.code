package main

import (
	"encoding/json"
	"fmt"
	"time"
)

type Subject struct {
	ID        int `json:id`
	CreatedAt time.Time
}

func (s *Subject) MarshalJSON() ([]byte, error) {

	return json.Marshal(struct {
		ID        int `json:id`
		CreatedAt string
	}{
		ID:        s.ID,
		CreatedAt: s.CreatedAt.Format("2006-01-02 15:04:05"),
	})
}

func main() {
	bytes, _ := json.Marshal(&Subject{111, time.Now()})
	fmt.Printf("%s", bytes)
}
