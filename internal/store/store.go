// Interface for storing todos in memory

package store

import "todo-list/internal/todo"

type Store interface {
	Create(text string) (todo.Todo, error)
	List(all bool, doneOnly bool) ([]todo.Todo, error)
	Get(id int64) (todo.Todo, error)
	UpdateText(id int64, text string) (todo.Todo, error)
	ToggleDone(id int64) (todo.Todo, error)
	Delete(id int64) error
}
