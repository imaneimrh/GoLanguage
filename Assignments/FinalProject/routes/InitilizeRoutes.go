package routes

import (
	. "FinalProject/models"
	. "FinalProject/stores"
	"context"
	"log"
	"net/http"
	"sync"
)

func InitializeRoutes() (*http.ServeMux, *InMemoryBookStore, *InMemoryAuthorStore, *InMemoryCustomerStore, *InMemoryOrderStore) {
	bookStore := &InMemoryBookStore{
		Mu:     sync.RWMutex{},
		Books:  make(map[int]Book),
		NextID: 1,
	}
	authorStore := &InMemoryAuthorStore{
		Mu:      sync.RWMutex{},
		Authors: make(map[int]Author),
		NextID:  1,
	}
	customerStore := &InMemoryCustomerStore{
		Mu:        sync.RWMutex{},
		Customers: make(map[int]Customer),
		NextID:    1,
	}
	orderStore := &InMemoryOrderStore{
		Mu:     sync.RWMutex{},
		Orders: make(map[int]Order),
		NextID: 1,
	}

	ctx := context.Background()
	if err := bookStore.LoadBooks(ctx, "books.json"); err != nil {
		log.Fatalf("Failed to load books: %v", err)
	}
	if err := authorStore.LoadAuthors(ctx, "authors.json"); err != nil {
		log.Fatalf("Failed to load authors: %v", err)
	}
	if err := customerStore.LoadCustomers(ctx, "customers.json"); err != nil {
		log.Fatalf("Failed to load customers: %v", err)
	}
	if err := orderStore.LoadOrders(ctx, "orders.json"); err != nil {
		log.Fatalf("Failed to load orders: %v", err)
	}

	router := http.NewServeMux()

	RegisterBookRoutes(router, bookStore, authorStore)
	RegisterAuthorRoutes(router, authorStore, bookStore)
	RegisterOrderRoutes(router, orderStore, customerStore, bookStore)
	RegisterCustomerRoutes(router, customerStore, orderStore)

	return router, bookStore, authorStore, customerStore, orderStore
}
