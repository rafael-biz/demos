package main

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
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

func (h *BookHandler) GetItems(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, h.service.GetItems())
}

func (h *BookHandler) GetItem(c *gin.Context) {
	id := c.Param("id")
	if book, found := h.service.GetItem(id); found {
		c.IndentedJSON(http.StatusOK, book)
		return
	} else {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found"})
	}
}

func (h *BookHandler) CreateItem(c *gin.Context) {
	var book Book
	if err := c.BindJSON(&book); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Bad Request"})
		return
	}
	h.service.CreateItem(book)
	c.IndentedJSON(http.StatusCreated, book)
}

func (h *BookHandler) UpdateItem(c *gin.Context) {
	id := c.Param("id")

	var book Book
	if err := c.BindJSON(&book); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	if h.service.UpdateItem(id, book) {
		c.Status(http.StatusOK)
		return
	} else {
		c.Status(http.StatusBadRequest)
		return
	}
}

func (h *BookHandler) DeleteItem(c *gin.Context) {
	id := c.Param("id")

	if h.service.DeleteItem(id) {
		c.Status(http.StatusOK)
		return
	} else {
		c.Status(http.StatusBadRequest)
		return
	}
}

// Main Function
func main() {
	repo := NewBookRepository()
	service := NewBookService(repo)
	handler := NewBookHandler(service)

	router := gin.Default()
	router.GET("/books", handler.GetItems)
	router.GET("/books/:id", handler.GetItem)
	router.POST("/books", handler.CreateItem)
	router.PUT("/books/:id", handler.UpdateItem)
	router.DELETE("/books/:id", handler.DeleteItem)

	router.Run("localhost:8081")
}
