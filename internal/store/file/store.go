// FileStore implementation
package file

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
	"todo-list/internal/store"
	"todo-list/internal/todo"
)

type FileStore struct {
	filepath string
	todos    []todo.Todo
	nextId   int64
	mu       sync.Mutex
}

// Load and return a new FileStore instance
func Load(filepath string) (*FileStore, error) {
	fs := &FileStore{
		filepath: filepath,
		mu:       sync.Mutex{},
	}

	if err := fs.loadFromFile(); err != nil {
		return nil, fmt.Errorf("failed to load store: %w", err)
	}
	return fs, nil
}

// load todos from json file
func (f *FileStore) loadFromFile() error {
	file, err := os.Open(f.filepath)

	// file does not exist, create it
	if os.IsNotExist(err) {
		f.todos = []todo.Todo{}
		f.nextId = 1
		return nil
	}

	if err != nil {
		return err
	}
	defer file.Close()

	var data struct {
		Version int64       `json:"version"`
		NextId  int64       `json:"next_id"`
		Items   []todo.Todo `json:"items"`
	}

	if err := json.NewDecoder(file).Decode(&data); err != nil {
		return fmt.Errorf("corrupted file: %w", err)
	}

	f.todos = data.Items
	f.nextId = data.NextId
	return nil
}

func (f *FileStore) saveToFile() error {
	dir := filepath.Dir(f.filepath)
	if err := os.MkdirAll(dir, 0755); err !=
		nil {
		return err
	}

	tempPath := f.filepath + ".tmp"
	file, err := os.Create(tempPath)
	if err != nil {
		return err
	}

	// closes file handle when function exits
	defer file.Close()

	data := struct {
		Version int64       `json:"version"`
		NextID  int64       `json:"next_id"`
		Items   []todo.Todo `json:"items"`
	}{Version: 1, NextID: f.nextId, Items: f.
		todos}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		return err
	}

	return os.Rename(tempPath, f.filepath)
}

// // Add Create, List, UpdateText, ToggleDone, Delete
func (f *FileStore) Create(text string) (todo.Todo, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	t := todo.Todo{
		ID:        f.nextId,
		Text:      text,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Done:      false,
	}

	if err := t.Validate(); err != nil {
		return todo.Todo{}, err
	}

	f.todos = append(f.todos, t)
	f.nextId++
	if err := f.saveToFile(); err != nil {
		return todo.Todo{}, err
	}

	return t, nil
}

// // Add Create, List, UpdateText, ToggleDone, Delete
// func (f *FileStore) List() ([]todo.Todo, error) {
//
// }
//
// // Add Create, List, UpdateText, ToggleDone, Delete
// func (f *FileStore) UpdateText(id int64) (todo.Todo, error) {
//
// }
//
// // Add Create, List, UpdateText, ToggleDone, Delete
// func (f *FileStore) ToggleDone(id int64) (todo.Todo, error) {
//
// }
//
// // Add Create, List, UpdateText, ToggleDone, Delete
// func (f *FileStore) Delete(id int64) error {
//
// }
