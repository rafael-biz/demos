package main

import (
	"context"
	"html/template"
	"log"
	"net/http"

	"github.com/jackc/pgx/v4"
)

type Book struct {
	ID       int
	Title    string
	AuthorID int
}

type Author struct {
	ID   int
	Name string
}

var db *pgx.Conn
var templates = template.Must(template.ParseGlob("templates/*.html"))

func main() {
	var err error
	db, err = pgx.Connect(context.Background(), "postgres://postgres:postgres@localhost:5432/bookstore")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close(context.Background())

	http.HandleFunc("/", listBooks)
	http.HandleFunc("/book/create", createBook)
	http.HandleFunc("/book/update", updateBook)
	http.HandleFunc("/book/delete", deleteBook)

	http.HandleFunc("/authors", listAuthors)
	http.HandleFunc("/author/create", createAuthor)
	http.HandleFunc("/author/update", updateAuthor)
	http.HandleFunc("/author/delete", deleteAuthor)

	log.Println("Server started at :8080")
	http.ListenAndServe(":8080", nil)
}

func listBooks(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(context.Background(), "SELECT id, title, author_id FROM books")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		var b Book
		if err := rows.Scan(&b.ID, &b.Title, &b.AuthorID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		books = append(books, b)
	}

	templates.ExecuteTemplate(w, "books.html", books)
}

func createBook(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		title := r.FormValue("title")
		authorID := r.FormValue("author_id")
		_, err := db.Exec(context.Background(), "INSERT INTO books (title, author_id) VALUES ($1, $2)", title, authorID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	templates.ExecuteTemplate(w, "create_book.html", nil)
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		id := r.FormValue("id")
		title := r.FormValue("title")
		authorID := r.FormValue("author_id")
		_, err := db.Exec(context.Background(), "UPDATE books SET title = $1, author_id = $2 WHERE id = $3", title, authorID, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	_, err := db.Exec(context.Background(), "DELETE FROM books WHERE id = $1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func listAuthors(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(context.Background(), "SELECT id, name FROM authors")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var authors []Author
	for rows.Next() {
		var a Author
		if err := rows.Scan(&a.ID, &a.Name); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		authors = append(authors, a)
	}

	templates.ExecuteTemplate(w, "authors.html", authors)
}

func createAuthor(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		name := r.FormValue("name")
		_, err := db.Exec(context.Background(), "INSERT INTO authors (name) VALUES ($1)", name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/authors", http.StatusSeeOther)
		return
	}
	templates.ExecuteTemplate(w, "create_author.html", nil)
}

func updateAuthor(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		id := r.FormValue("id")
		name := r.FormValue("name")
		_, err := db.Exec(context.Background(), "UPDATE authors SET name = $1 WHERE id = $2", name, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/authors", http.StatusSeeOther)
		return
	}
}

func deleteAuthor(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	_, err := db.Exec(context.Background(), "DELETE FROM authors WHERE id = $1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/authors", http.StatusSeeOther)
}
