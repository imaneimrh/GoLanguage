package utils

import (
	"context"
	"log"

	. "FinalProject/stores"
)

func SaveAllData(ctx context.Context, bookStore *InMemoryBookStore, authorStore *InMemoryAuthorStore, customerStore *InMemoryCustomerStore, orderStore *InMemoryOrderStore) {
	log.Println("Saving data to files...")

	if err := bookStore.SaveBooks(ctx, "books.json"); err != nil {
		log.Printf("Failed to save books: %v", err)
	}

	if err := authorStore.SaveAuthors(ctx, "authors.json"); err != nil {
		log.Printf("Failed to save authors: %v", err)
	}

	if err := customerStore.SaveCustomers(ctx, "customers.json"); err != nil {
		log.Printf("Failed to save customers: %v", err)
	}

	if err := orderStore.SaveOrders(ctx, "orders.json"); err != nil {
		log.Printf("Failed to save orders: %v", err)
	}

	log.Println("Data saving completed.")
}
