package routes

import (
	. "FinalProject/controllers"
	. "FinalProject/stores"
	"net/http"
)

func RegisterOrderRoutes(mux *http.ServeMux, orderStore *InMemoryOrderStore, customerStore *InMemoryCustomerStore, bookStore *InMemoryBookStore) {
	mux.HandleFunc("/orders", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			CreateOrderHandler(w, r, customerStore, bookStore, orderStore)
		case "GET":
			GetAllOrdersHandler(w, r, orderStore)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/orders/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			GetOrderHandler(w, r, orderStore)
		case "PUT":
			UpdateOrderHandler(w, r, orderStore, customerStore, bookStore)
		case "DELETE":
			DeleteOrderHandler(w, r, orderStore)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/orders/history", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			ViewOrderHistoryHandler(w, r, orderStore)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/orders/timerange", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			FetchOrdersWithinTimeHandler(w, r, orderStore)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/reports", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			ListReportsHandler(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

}
