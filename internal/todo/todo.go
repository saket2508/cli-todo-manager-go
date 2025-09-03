package todo

import (
	"errors"
	"strings"
	"time"
)

// Todo represents a todo item with all its properties
type Todo struct {
	ID          int64      `json:"id"`
	Text        string     `json:"text"`
	Done        bool       `json:"done"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

// Validate checks if a todo has valid data
func (t *Todo) Validate() error {
	if strings.TrimSpace(t.Text) == "" {
		return errors.New("todo text cannot be empty")
	}
	if len(t.Text) > 4096 { // 4KB limit as per DESIGN.md
		return errors.New("todo text is too long (max 4096 characters)")
	}
	return nil
}

// Store defines the interface for todo persistence
type Store interface {
	Create(text string) (Todo, error)
	List(all bool, doneOnly bool) ([]Todo, error)
	Get(id int64) (Todo, error)
	UpdateText(id int64, text string) (Todo, error)
	ToggleDone(id int64) (Todo, error)
	Delete(id int64) error
}
