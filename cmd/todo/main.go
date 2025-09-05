// this is the CLI entry point
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"todo-list/internal/store/file"
)

func main() {
	// Define flags
	var dataPath = flag.String("data", "", "Path to data file (default: ~/.config/todo-list/todos.json)")
	var jsonOutput = flag.Bool("json", false, "Output in JSON format")
	flag.Parse()

	// Determine data file path
	filepath := getDataPath(*dataPath)

	// Load store
	store, err := file.Load(filepath)
	if err != nil {
		log.Fatalf("Failed to load store: %v", err)
	}

	// Parse subcommands
	args := flag.Args()
	if len(args) == 0 {
		fmt.Println("Usage: todo <command> [args]")
		fmt.Println("Commands: add, list, done, undone, rm")
		os.Exit(1)
	}

	cmd := args[0]
	switch cmd {
	case "add":
		if len(args) < 2 {
			fmt.Println("Usage: todo add <text>")
			os.Exit(1)
		}
		text := strings.Join(args[1:], " ")
		todo, err := store.Create(text)
		if err != nil {
			log.Fatalf("Failed to add todo: %v", err)
		}
		if *jsonOutput {
			fmt.Printf("{\"id\": %d, \"text\": \"%s\"}\n", todo.ID, todo.Text)
		} else {
			fmt.Printf("Added todo %d: %s\n", todo.ID, todo.Text)
		}
	case "list":
		all := false
		doneOnly := false
		for _, arg := range args[1:] {
			if arg == "--all" {
				all = true
			}
			if arg == "--done" {
				doneOnly = true
			}
		}
		todos, err := store.List(all, doneOnly)
		if err != nil {
			log.Fatalf("Failed to list todos: %v", err)
		}
		if *jsonOutput {
			fmt.Print("[")
			for i, t := range todos {
				if i > 0 {
					fmt.Print(",")
				}
				fmt.Printf("{\"id\": %d, \"text\": \"%s\", \"done\": %t}", t.ID, t.Text, t.Done)
			}
			fmt.Println("]")
		} else {
			for _, t := range todos {
				status := " "
				if t.Done {
					status = "✓"
				}
				fmt.Printf("%d %s %s\n", t.ID, status, t.Text)
			}
		}
	case "done":
		if len(args) < 2 {
			fmt.Println("Usage: todo done <id>")
			os.Exit(1)
		}
		id, err := strconv.ParseInt(args[1], 10, 64)
		if err != nil {
			log.Fatalf("Invalid ID: %v", err)
		}
		todo, err := store.ToggleDone(id)
		if err != nil {
			log.Fatalf("Failed to mark done: %v", err)
		}
		if *jsonOutput {
			fmt.Printf("{\"id\": %d, \"done\": %t}\n", todo.ID, todo.Done)
		} else {
			fmt.Printf("Marked todo %d as done\n", todo.ID)
		}
	case "undone":
		if len(args) < 2 {
			fmt.Println("Usage: todo undone <id>")
			os.Exit(1)
		}
		id, err := strconv.ParseInt(args[1], 10, 64)
		if err != nil {
			log.Fatalf("Invalid ID: %v", err)
		}
		todo, err := store.ToggleDone(id)
		if err != nil {
			log.Fatalf("Failed to mark undone: %v", err)
		}
		if *jsonOutput {
			fmt.Printf("{\"id\": %d, \"done\": %t}\n", todo.ID, todo.Done)
		} else {
			fmt.Printf("Marked todo %d as undone\n", todo.ID)
		}
	case "rm":
		if len(args) < 2 {
			fmt.Println("Usage: todo rm <id>")
			os.Exit(1)
		}
		id, err := strconv.ParseInt(args[1], 10, 64)
		if err != nil {
			log.Fatalf("Invalid ID: %v", err)
		}
		err = store.Delete(id)
		if err != nil {
			log.Fatalf("Failed to delete todo: %v", err)
		}
		if *jsonOutput {
			fmt.Printf("{\"deleted\": %d}\n", id)
		} else {
			fmt.Printf("Deleted todo %d\n", id)
		}
	default:
		fmt.Printf("Unknown command: %s\n", cmd)
		os.Exit(1)
	}
}

func getDataPath(customPath string) string {
	if customPath != "" {
		return customPath
	}
	configDir, err := os.UserConfigDir()
	if err != nil {
		// Fallback to home dir
		home, _ := os.UserHomeDir()
		return filepath.Join(home, ".todo.json")
	}
	return filepath.Join(configDir, "todo-list", "todos.json")
}
