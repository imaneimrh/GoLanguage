package routes

import (
	. "FinalProject/controllers"
	. "FinalProject/stores"
	"net/http"
)

func RegisterCustomerRoutes(mux *http.ServeMux, customerStore *InMemoryCustomerStore, orderStore *InMemoryOrderStore) {
	mux.HandleFunc("/customers", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			CreateCustomerHandler(w, r, customerStore)
		case "GET":
			GetAllCustomersHandler(w, r, customerStore)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/customers/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			GetCustomerByIDHandler(w, r, customerStore)
		case "PUT":
			UpdateCustomerHandler(w, r, customerStore)
		case "DELETE":
			DeleteCustomerHandler(w, r, customerStore, orderStore)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}
