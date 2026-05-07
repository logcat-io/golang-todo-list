package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"todo-cli/internal/todo"
)

type Server struct {
	store *todo.Store
}

type createTodoRequest struct {
	Title string `json:"title"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func main() {
	store := todo.NewStore()
	_, _ = store.AddNext("learn go")
	_, _ = store.AddNext("write tests")

	server := &Server{store: store}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", server.handleHealth)
	mux.HandleFunc("GET /todos", server.handleListTodos)
	mux.HandleFunc("POST /todos", server.handleCreateTodo)
	mux.HandleFunc("PATCH /todos/", server.handleCompleteTodo)
	mux.HandleFunc("DELETE /todos/", server.handleDeleteTodo)
	mux.Handle("GET /", http.FileServer(http.Dir("web")))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := ":" + port
	log.Printf("todo api listening on %s", addr)
	if err := http.ListenAndServe(addr, logRequest(mux)); err != nil {
		log.Fatal(err)
	}
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) handleListTodos(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, s.store.List())
}

func (s *Server) handleCreateTodo(w http.ResponseWriter, r *http.Request) {
	var req createTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid json body"})
		return
	}

	item, err := s.store.AddNext(strings.TrimSpace(req.Title))
	if err != nil {
		if errors.Is(err, todo.ErrTitleRequired) {
			writeJSON(w, http.StatusBadRequest, errorResponse{Error: "title is required"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "internal server error"})
		return
	}

	writeJSON(w, http.StatusCreated, item)
}

func (s *Server) handleCompleteTodo(w http.ResponseWriter, r *http.Request) {
	id, ok := parseTodoActionPath(w, r.URL.Path, "/todos/", "/complete")
	if !ok {
		return
	}

	if err := s.store.Complete(id); err != nil {
		if errors.Is(err, todo.ErrTodoNotFound) {
			writeJSON(w, http.StatusNotFound, errorResponse{Error: "todo not found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "internal server error"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "completed"})
}

func (s *Server) handleDeleteTodo(w http.ResponseWriter, r *http.Request) {
	id, ok := parseTodoIDPath(w, r.URL.Path, "/todos/")
	if !ok {
		return
	}

	if err := s.store.Remove(id); err != nil {
		if errors.Is(err, todo.ErrTodoNotFound) {
			writeJSON(w, http.StatusNotFound, errorResponse{Error: "todo not found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "internal server error"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func parseTodoActionPath(w http.ResponseWriter, path string, prefix string, suffix string) (int, bool) {
	if !strings.HasPrefix(path, prefix) || !strings.HasSuffix(path, suffix) {
		writeJSON(w, http.StatusNotFound, errorResponse{Error: "not found"})
		return 0, false
	}
	rawID := strings.TrimSuffix(strings.TrimPrefix(path, prefix), suffix)
	id, err := strconv.Atoi(rawID)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid todo id"})
		return 0, false
	}
	return id, true
}

func parseTodoIDPath(w http.ResponseWriter, path string, prefix string) (int, bool) {
	if !strings.HasPrefix(path, prefix) {
		writeJSON(w, http.StatusNotFound, errorResponse{Error: "not found"})
		return 0, false
	}
	rawID := strings.TrimPrefix(path, prefix)
	id, err := strconv.Atoi(rawID)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid todo id"})
		return 0, false
	}
	return id, true
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if status == http.StatusNoContent {
		return
	}
	if err := json.NewEncoder(w).Encode(value); err != nil {
		fmt.Println("write json failed:", err)
	}
}

func logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
