package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"log"
	"os"
	"path/filepath"
	"strings"
	"todo-list/internal"
)

type model struct {
	todos    []internal.Todo
	store    *internal.FileStore
	cursor   int
	view     string
	input    string
	filepath string
}

func initialModel(filepath string) model {
	store, _ := internal.Load(filepath)
	todos, _ := store.List(false, false)
	return model{todos: todos, store: store, view: "list", filepath: filepath}
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+q", "ctrl+c":
			return m, tea.Quit
		case "ctrl+n": // Add todo
			m.view = "add"
			m.input = ""
		case "enter":
			if m.view == "add" && strings.TrimSpace(m.input) != "" {
				m.store.Create(m.input)
				m.todos, _ = m.store.List(false, false)
				m.view = "list"
			}
		case "up":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down":
			if m.cursor < len(m.todos)-1 {
				m.cursor++
			}
		case "ctrl+d": // Toggle done
			if len(m.todos) > 0 {
				m.store.ToggleDone(m.todos[m.cursor].ID)
				m.todos, _ = m.store.List(false, false)
			}
		case "ctrl+x": // Remove
			if len(m.todos) > 0 {
				m.store.Delete(m.todos[m.cursor].ID)
				m.todos, _ = m.store.List(false, false)
				if m.cursor >= len(m.todos) {
					m.cursor = len(m.todos) - 1
				}
			}
		case "backspace":
			if m.view == "add" && len(m.input) > 0 {
				m.input = m.input[:len(m.input)-1]
			}
		default:
			if m.view == "add" {
				m.input += msg.String()
			}
		}
	}
	return m, nil
}

func (m model) View() string {
	if m.view == "add" {
		return fmt.Sprintf("Add Todo: %s\n\n(enter to save, esc to cancel)", m.input)
	}
	s := "Todos:\n\n"
	for i, todo := range m.todos {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		status := " "
		if todo.Done {
			status = "✓"
		}
		s += fmt.Sprintf("%s %s %s\n", cursor, status, todo.Text)
	}
	s += "\n(ctrl-n)dd (ctrl-d)one (ctrl-x)emove (ctrl-c)uit"
	return s
}

func main() {
	filepath := getDataPath("")
	p := tea.NewProgram(initialModel(filepath))
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

// Helper function (copy from cmd/todo/main.go)
func getDataPath(customPath string) string {
	if customPath != "" {
		return customPath
	}
	configDir, err := os.UserConfigDir()
	if err != nil {
		home, _ := os.UserHomeDir()
		return filepath.Join(home, ".todo.json")
	}
	return filepath.Join(configDir, "todo-list", "todos.json")
}
