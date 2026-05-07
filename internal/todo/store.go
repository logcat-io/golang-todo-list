package todo

import (
	"errors"
	"fmt"
	"sort"
	"sync"
)

var (
	ErrTitleRequired = errors.New("title is required")
	ErrTodoNotFound  = errors.New("todo not found")
)

type Store struct {
	mu     sync.Mutex
	todos  map[int]Todo
	nextID int
}

func NewStore() *Store {
	return &Store{
		todos:  make(map[int]Todo),
		nextID: 1,
	}
}

func (s *Store) Add(id int, title string) error {
	if title == "" {
		return ErrTitleRequired
	}

	if _, exists := s.todos[id]; exists {
		return fmt.Errorf("todo with id %d already exists", id)
	}
	s.todos[id] = Todo{
		ID:    id,
		Title: title,
		Done:  false,
	}

	return nil
}

func (s *Store) AddNext(title string) (Todo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if title == "" {
		return Todo{}, ErrTitleRequired
	}

	id := s.nextID
	s.nextID++

	todo := Todo{
		ID:    id,
		Title: title,
		Done:  false,
	}

	s.todos[id] = todo
	return todo, nil
}

func (s *Store) Complete(id int) error {
	todo, exists := s.todos[id]
	if !exists {
		return fmt.Errorf("%w: %d", ErrTodoNotFound, id)
	}

	todo.MarkDone()
	s.todos[id] = todo
	return nil
}

func (s *Store) Remove(id int) error {
	if _, exists := s.todos[id]; !exists {
		return fmt.Errorf("%w: %d", ErrTodoNotFound, id)
	}

	delete(s.todos, id)
	return nil
}

func (s *Store) List() []Todo {
	todos := make([]Todo, 0, len(s.todos))

	for _, todo := range s.todos {
		todos = append(todos, todo)
	}

	sort.Slice(todos, func(i, j int) bool {
		return todos[i].ID < todos[j].ID
	})

	return todos
}
