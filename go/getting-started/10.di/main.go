package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

// Model
type Book struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Repository
type BookRepository struct {
	mu    sync.RWMutex
	books map[string]Book
}

func NewBookRepository() *BookRepository {
	return &BookRepository{
		books: make(map[string]Book),
	}
}

func (r *BookRepository) GetAll() []Book {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]Book, 0)
	for _, book := range r.books {
		result = append(result, book)
	}
	return result
}

func (r *BookRepository) GetByID(id string) (Book, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	book, exists := r.books[id]
	return book, exists
}

func (r *BookRepository) Create(book Book) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.books[book.ID] = book
}

func (r *BookRepository) Update(id string, book Book) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.books[id]; !exists {
		return false
	}
	r.books[id] = book
	return true
}

func (r *BookRepository) Delete(id string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.books[id]; !exists {
		return false
	}
	delete(r.books, id)
	return true
}

// Service Layer
type BookService struct {
	repo *BookRepository
}

func NewBookService(repo *BookRepository) *BookService {
	return &BookService{repo: repo}
}

func (s *BookService) GetItems() []Book {
	return s.repo.GetAll()
}

func (s *BookService) GetItem(id string) (Book, bool) {
	return s.repo.GetByID(id)
}

func (s *BookService) CreateItem(book Book) {
	s.repo.Create(book)
}

func (s *BookService) UpdateItem(id string, book Book) bool {
	return s.repo.Update(id, book)
}

func (s *BookService) DeleteItem(id string) bool {
	return s.repo.Delete(id)
}

// Handler Layer
type BookHandler struct {
	service *BookService
}

func NewBookHandler(service *BookService) *BookHandler {
	return &BookHandler{service: service}
}

func (h *BookHandler) GetItems(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(h.service.GetItems())
}

func (h *BookHandler) GetItem(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if book, found := h.service.GetItem(id); found {
		json.NewEncoder(w).Encode(book)
	} else {
		http.Error(w, "Book not found", http.StatusNotFound)
	}
}

func (h *BookHandler) CreateItem(w http.ResponseWriter, r *http.Request) {
	var book Book
	json.NewDecoder(r.Body).Decode(&book)
	h.service.CreateItem(book)
	w.WriteHeader(http.StatusCreated)
}

func (h *BookHandler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var book Book
	json.NewDecoder(r.Body).Decode(&book)
	if h.service.UpdateItem(id, book) {
		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "Book not found", http.StatusNotFound)
	}
}

func (h *BookHandler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if h.service.DeleteItem(id) {
		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "Book not found", http.StatusNotFound)
	}
}

// Main Function
func main() {
	repo := NewBookRepository()
	service := NewBookService(repo)
	handler := NewBookHandler(service)

	r := mux.NewRouter()
	r.HandleFunc("/books", handler.GetItems).Methods("GET")
	r.HandleFunc("/books/{id}", handler.GetItem).Methods("GET")
	r.HandleFunc("/books", handler.CreateItem).Methods("POST")
	r.HandleFunc("/books/{id}", handler.UpdateItem).Methods("PUT")
	r.HandleFunc("/books/{id}", handler.DeleteItem).Methods("DELETE")

	fmt.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
