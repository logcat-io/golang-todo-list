package todo

import (
	"errors"
	"testing"
)

func TestStore_Add(t *testing.T) {
	tests := []struct {
		name      string
		id        int
		title     string
		wantError error
	}{
		{
			name:      "valid todo is added",
			id:        1,
			title:     "learn go",
			wantError: nil,
		},
		{
			name:      "empty title returns ErrTitleRequired",
			id:        1,
			title:     "",
			wantError: ErrTitleRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := NewStore()

			err := store.Add(tt.id, tt.title)
			if !errors.Is(err, tt.wantError) {
				t.Fatalf("error = %v, want %v", err, tt.wantError)
			}
		})
	}
}

func TestStore_Complete(t *testing.T) {
	store := NewStore()
	if err := store.Add(1, "learn go"); err != nil {
		t.Fatalf("unexpected add error: %v", err)
	}

	if err := store.Complete(1); err != nil {
		t.Fatalf("unexpected complete error: %v", err)
	}

	todos := store.List()
	if len(todos) != 1 {
		t.Fatalf("len(todos) = %d, want 1", len(todos))
	}
	if !todos[0].Done {
		t.Fatalf("todo.Done = false, want true")
	}
}

func TestStore_Complete_NotFound(t *testing.T) {
	store := NewStore()

	err := store.Complete(999)
	if !errors.Is(err, ErrTodoNotFound) {
		t.Fatalf("error = %v, want %v", err, ErrTodoNotFound)
	}
}

func TestStore_List_SortedByID(t *testing.T) {
	store := NewStore()
	_ = store.Add(3, "third")
	_ = store.Add(1, "first")
	_ = store.Add(2, "second")

	todos := store.List()
	gotIDs := []int{todos[0].ID, todos[1].ID, todos[2].ID}
	wantIDs := []int{1, 2, 3}

	for i := range wantIDs {
		if gotIDs[i] != wantIDs[i] {
			t.Fatalf("gotIDs = %v, want %v", gotIDs, wantIDs)
		}
	}
}
