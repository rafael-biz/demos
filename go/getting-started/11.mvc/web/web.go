package web

import (
	"net/http"

	"example.com/model"
	"example.com/services"
	"github.com/gin-gonic/gin"
)

type BookHandler struct {
	service *services.BookService
}

func NewBookHandler(service *services.BookService) *BookHandler {
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
	var book model.Book
	if err := c.BindJSON(&book); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Bad Request"})
		return
	}
	h.service.CreateItem(book)
	c.IndentedJSON(http.StatusCreated, book)
}

func (h *BookHandler) UpdateItem(c *gin.Context) {
	id := c.Param("id")

	var book model.Book
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

func Run(service *services.BookService) {
	handler := NewBookHandler(service)

	router := gin.Default()
	router.GET("/books", handler.GetItems)
	router.GET("/books/:id", handler.GetItem)
	router.POST("/books", handler.CreateItem)
	router.PUT("/books/:id", handler.UpdateItem)
	router.DELETE("/books/:id", handler.DeleteItem)

	router.Run("localhost:8081")
}
