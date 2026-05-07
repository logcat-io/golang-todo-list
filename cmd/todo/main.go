package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"todo-cli/internal/todo"
)

func main() {
	store := todo.NewStore()

	if err := store.Add(1, "learn go"); err != nil {
		fmt.Println("add failed:", err)
		os.Exit(1)
	}
	if err := store.Add(2, "write tests"); err != nil {
		fmt.Println("add failed:", err)
		os.Exit(1)
	}
	if err := store.Add(3, "write summary"); err != nil {
		fmt.Println("add failed:", err)
		os.Exit(1)
	}
	if err := store.Complete(1); err != nil {
		fmt.Println("complete failed:", err)
		os.Exit(1)
	}

	printTodos(store.List())

	if len(os.Args) >= 3 && os.Args[1] == "complete" {
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("invalid id:", err)
			os.Exit(1)
		}
		if err := store.Complete(id); err != nil {
			if errors.Is(err, todo.ErrTodoNotFound) {
				fmt.Println("todo not found")
				os.Exit(1)
			}
			fmt.Println("complete failed:", err)
			os.Exit(1)
		}
	}
}

func printTodos(todos []todo.Todo) {
	doneCount := 0
	for _, item := range todos {
		if item.IsDone() {
			doneCount++
		}
	}

	fmt.Println("┌──────────────────────────────────────┐")
	fmt.Println("│              TODO LIST               │")
	fmt.Println("├────┬──────────┬──────────────────────┤")
	fmt.Println("│ ID │ STATUS   │ TITLE                │")
	fmt.Println("├────┼──────────┼──────────────────────┤")

	for _, item := range todos {
		status := "TODO"
		if item.IsDone() {
			status = "DONE"
		}
		fmt.Printf("│ %-2d │ %-8s │ %-20s │\n", item.ID, status, item.Title)
	}

	fmt.Println("└────┴──────────┴──────────────────────┘")
	fmt.Printf("total=%d done=%d pending=%d\n", len(todos), doneCount, len(todos)-doneCount)
}
