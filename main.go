package main

import (
	"fmt"
	"time"
)

func main() {
	timestamp := time.Now().Format(time.RFC3339)
	fmt.Printf("Current datetime: %s\n", timestamp)

	fmt.Println("this is a todo app written in go")
}
