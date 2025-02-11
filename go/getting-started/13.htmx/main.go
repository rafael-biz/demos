package main

import (
	"html/template"
	"log"
	"net/http"
)

type Item struct {
	ID    int
	Name  string
	Price float64
}

var items = []Item{
	{ID: 1, Name: "Item 1", Price: 10.99},
	{ID: 2, Name: "Item 2", Price: 20.99},
	{ID: 3, Name: "Item 3", Price: 30.99},
}

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/items", itemsHandler)

	log.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, nil)
}

func itemsHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/items.html"))
	tmpl.Execute(w, items)
}
