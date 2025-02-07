package main

import (
	"example.com/persistence"
	"example.com/services"
	"example.com/web"
)

// Main Function
func main() {
	repo := persistence.NewBookRepository()
	service := services.NewBookService(repo)
	web.Run(service)
}
