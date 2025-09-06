// FileStore implementation
package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"
)

type FileStore struct {
	filepath string
	todos    []Todo
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
		f.todos = []Todo{}
		f.nextId = 1
		return nil
	}

	if err != nil {
		return err
	}
	defer file.Close()

	var data struct {
		Version int64       `json:"version"`
		NextID  int64       `json:"next_id"`
		Items   []Todo `json:"items"`
	}

	if err := json.NewDecoder(file).Decode(&data); err != nil {
		return fmt.Errorf("corrupted file: %w", err)
	}

	f.todos = data.Items
	f.nextId = data.NextID
	return nil
}

func (f *FileStore) saveToFile() error {
	dir := filepath.Dir(f.filepath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	tempPath := f.filepath + ".tmp"
	file, err := os.Create(tempPath)
	if err != nil {
		return err
	}
	defer file.Close()

	data := struct {
		Version int64       `json:"version"`
		NextID  int64       `json:"next_id"`
		Items   []Todo `json:"items"`
	}{Version: 1, NextID: f.nextId, Items: f.todos}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		return err
	}

	return os.Rename(tempPath, f.filepath)
}

// // Add Create, List, UpdateText, ToggleDone, Delete
func (f *FileStore) Create(text string) (Todo, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	t := Todo{
		ID:        f.nextId,
		Text:      text,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Done:      false,
	}

	if err := t.Validate(); err != nil {
		return Todo{}, err
	}

	f.todos = append(f.todos, t)
	f.nextId++
	if err := f.saveToFile(); err != nil {
		return Todo{}, err
	}

	return t, nil
}

// Add Create, List, UpdateText, ToggleDone, Delete
func (f *FileStore) List(all bool, doneOnly bool) ([]Todo, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	var result []Todo
	for _, t := range f.todos {
		if doneOnly && !t.Done {
			continue
		}
		if !all && t.Done {
			continue
		}
		result = append(result, t)
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].Done != result[j].Done {
			return !result[i].Done
		}
		if result[i].Done {
			return result[i].CompletedAt.After(*result[j].CompletedAt)
		}
		return result[i].UpdatedAt.After(result[j].UpdatedAt)
	})

	return result, nil
}

func (f *FileStore) Get(id int64) (Todo, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	for _, t := range f.todos {
		if t.ID == id {
			return t, nil
		}
	}
	return Todo{}, fmt.Errorf("todo with id %d not found", id)
}

func (f *FileStore) UpdateText(id int64, text string) (Todo, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	for i, t := range f.todos {
		if t.ID == id {
			f.todos[i].UpdateText(text)
			if err := f.todos[i].Validate(); err != nil {
				return Todo{}, err
			}
			if err := f.saveToFile(); err != nil {
				return Todo{}, err
			}
			return f.todos[i], nil
		}
	}
	return Todo{}, fmt.Errorf("todo with id %d not found", id)
}

func (f *FileStore) ToggleDone(id int64) (Todo, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	for i, t := range f.todos {
		if t.ID == id {
			if t.Done {
				f.todos[i].MarkIncomplete()
			} else {
				f.todos[i].MarkCompleted()
			}
			if err := f.saveToFile(); err != nil {
				return Todo{}, err
			}
			return f.todos[i], nil
		}
	}
	return Todo{}, fmt.Errorf("todo with id %d not found", id)
}

func (f *FileStore) Delete(id int64) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	for i, t := range f.todos {
		if t.ID == id {
			f.todos = append(f.todos[:i], f.todos[i+1:]...)
			return f.saveToFile()
		}
	}
	return fmt.Errorf("todo with id %d not found", id)
}
