// Interface for storing todos in memory

package internal

type Store interface {
	Create(text string) (Todo, error)
	List(all bool, doneOnly bool) ([]Todo, error)
	Get(id int64) (Todo, error)
	UpdateText(id int64, text string) (Todo, error)
	ToggleDone(id int64) (Todo, error)
	Delete(id int64) error
}
