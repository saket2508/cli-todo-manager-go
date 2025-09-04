// contains types and interfaces

package todo

import (
	"errors"
	"strings"
	"time"
)

type Todo struct {
	ID          int64      `json:"id"`
	Text        string     `json:"text"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Done        bool       `json:"done"`
	CompletedAt *time.Time `json:"completed_at"`
}

// validate todo contents
func (t Todo) Validate() error {
	text := strings.TrimSpace(t.Text)
	if text == "" {
		return errors.New("text cannot be empty")
	}
	if len(text) > 4096 { // 4KB text limit
		return errors.New("text too long (max 4096 characters)")
	}
	return nil
}

// check if todo is marked done
func (t Todo) IsCompleted() bool {
	return t.Done && t.CompletedAt != nil
}

// sets done=true and completed timestamp
func (t *Todo) MarkCompleted() {
	t.Done = true
	now := time.Now()
	t.CompletedAt = &now
	t.UpdatedAt = now
}

func (t *Todo) MarkIncomplete() {
	t.Done = false
	t.CompletedAt = nil
	t.UpdatedAt = time.Now()
}

func (t *Todo) UpdateText(newText string) {
	t.Text = strings.TrimSpace(newText)
	t.UpdatedAt = time.Now()
}
