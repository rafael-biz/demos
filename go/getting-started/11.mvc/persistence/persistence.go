package persistence

import (
	"sync"

	"example.com/model"
)

// Repository
type BookRepository struct {
	mu    sync.RWMutex
	books map[string]model.Book
}

func NewBookRepository() *BookRepository {
	return &BookRepository{
		books: make(map[string]model.Book),
	}
}

func (r *BookRepository) GetAll() []model.Book {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]model.Book, 0)
	for _, book := range r.books {
		result = append(result, book)
	}
	return result
}

func (r *BookRepository) GetByID(id string) (model.Book, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	book, exists := r.books[id]
	return book, exists
}

func (r *BookRepository) Create(book model.Book) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.books[book.ID] = book
}

func (r *BookRepository) Update(id string, book model.Book) bool {
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
