package main

import (
	"testing"
)

// Unit Tests for Service Layer
func TestBookService(t *testing.T) {
	repo := NewBookRepository()
	service := NewBookService(repo)

	if service.GetItems() == nil {
		t.Error("Expected 0 book, got nil")
	}

	if len(service.GetItems()) != 0 {
		t.Errorf("Expected 0 book, got %d", len(service.GetItems()))
	}

	book := Book{ID: "1", Name: "Test Book"}
	service.CreateItem(book)

	if len(service.GetItems()) != 1 {
		t.Errorf("Expected 1 book, got %d", len(service.GetItems()))
	}

	savedItem, found := service.GetItem("1")
	if !found || savedItem.Name != "Test Book" {
		t.Errorf("Book retrieval failed")
	}

	updatedItem := Book{ID: "1", Name: "Updated Book"}
	if !service.UpdateItem("1", updatedItem) {
		t.Errorf("Book update failed")
	}

	if item, found := service.GetItem("1"); found && item.Name != "Updated Book" {
		t.Errorf("Book update did not persist")
	}

	if !service.DeleteItem("1") || len(service.GetItems()) != 0 {
		t.Errorf("Book deletion failed")
	}
}
