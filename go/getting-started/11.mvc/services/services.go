package services

import (
	"example.com/model"
	"example.com/persistence"
)

// Service Layer
type BookService struct {
	repo *persistence.BookRepository
}

func NewBookService(repo *persistence.BookRepository) *BookService {
	return &BookService{repo: repo}
}

func (s *BookService) GetItems() []model.Book {
	return s.repo.GetAll()
}

func (s *BookService) GetItem(id string) (model.Book, bool) {
	return s.repo.GetByID(id)
}

func (s *BookService) CreateItem(book model.Book) {
	s.repo.Create(book)
}

func (s *BookService) UpdateItem(id string, book model.Book) bool {
	return s.repo.Update(id, book)
}

func (s *BookService) DeleteItem(id string) bool {
	return s.repo.Delete(id)
}
