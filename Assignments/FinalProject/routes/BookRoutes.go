package routes

import (
	. "FinalProject/controllers"
	. "FinalProject/stores"
	"net/http"
)

func RegisterBookRoutes(mux *http.ServeMux, bookStore *InMemoryBookStore, authorStore *InMemoryAuthorStore) {
	mux.HandleFunc("/books", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			CreateBookHandler(w, r, bookStore, authorStore)
		case "GET":
			SearchBookHandler(w, r, bookStore)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/books/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			GetBookHandler(w, r, bookStore, authorStore)
		case "PUT":
			UpdateBookHandler(w, r, bookStore, authorStore)
		case "DELETE":
			DeleteBookHandler(w, r, bookStore)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}
