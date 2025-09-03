package main

import (
	"fmt"
	"time"
)

type Todo struct {
	id          int64
	text        string
	createdAt   time.Time
	updatedAt   time.Time
	done        bool
	completedAt *time.Time
}

type Store interface {
	Create(text string) (Todo, error)
	List(all bool, doneOnly bool) ([]Todo, error)
	Get(id int64) (Todo, error)
	UpdateText(id int64, text string) (Todo, error)
	ToggleDone(id int64) (Todo, error)
	Delete(id int64) error
}

func main() {
	timestamp := time.Now().Format(time.RFC3339)
	fmt.Printf("Current datetime: %s\n", timestamp)

	fmt.Println("this is a todo app written in go")
}
