package routes

import (
	. "FinalProject/controllers"
	. "FinalProject/stores"
	"net/http"
)

func RegisterAuthorRoutes(mux *http.ServeMux, authorStore *InMemoryAuthorStore, bookStore *InMemoryBookStore) {
	mux.HandleFunc("/authors", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			CreateAuthorHandler(w, r, authorStore)
		case "GET":
			ListAllHandler(w, r, authorStore)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/authors/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			GetAuthorByIdHandler(w, r, authorStore)
		case "PUT":
			UpdateAuthorHandler(w, r, authorStore)
		case "DELETE":
			DeleteAuthorHandler(w, r, authorStore, bookStore)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}
