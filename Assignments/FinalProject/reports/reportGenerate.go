package reports

import (
	. "FinalProject/models"
	. "FinalProject/stores"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

func GenerateSalesReport(ctx context.Context, orderStore *InMemoryOrderStore, bookStore *InMemoryBookStore, interval time.Duration) error {
	log.Println("Starting GenerateSalesReport function...")

	startTime := time.Now().Add(-interval).Truncate(0)
	endTime := time.Now().Truncate(0)
	log.Printf("Generating sales report for orders between %s and %s\n", startTime, endTime)

	var orders []Order
	orderStore.Mu.RLock()
	for _, order := range orderStore.Orders {
		log.Printf("Checking order ID %d with CreatedAt %s\n", order.ID, order.CreatedAt)
		if !order.CreatedAt.Before(startTime) && !order.CreatedAt.After(endTime) {
			log.Printf("Order ID %d is within the time range.\n", order.ID)
			orders = append(orders, order)
		} else {
			log.Printf("Order ID %d is outside the time range.\n", order.ID)
		}
	}
	orderStore.Mu.RUnlock()

	if len(orders) == 0 {
		log.Println("No orders found for the specified time range.")
		return nil
	}

	var totalRevenue float64
	var totalOrders int
	bookSalesMap := make(map[int]int)

	for _, order := range orders {
		totalRevenue += order.TotalPrice
		totalOrders++
		for _, item := range order.Items {
			bookSalesMap[item.Book.ID] += item.Quantity
		}
	}

	var topSellingBooks []BookSales
	var maxQuantity int
	for _, quantity := range bookSalesMap {
		if quantity > maxQuantity {
			maxQuantity = quantity
		}
	}

	bookStore.Mu.RLock()
	for bookID, quantity := range bookSalesMap {
		if quantity == maxQuantity {
			if book, exists := bookStore.Books[bookID]; exists {
				topSellingBooks = append(topSellingBooks, BookSales{
					Book:     book,
					Quantity: quantity,
				})
			}
		}
	}
	bookStore.Mu.RUnlock()

	report := SalesReport{
		Timestamp:       time.Now(),
		TotalRevenue:    totalRevenue,
		TotalOrders:     totalOrders,
		TopSellingBooks: topSellingBooks,
	}

	outputDir := "output-reports"
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		log.Printf("Failed to create output directory: %v\n", err)
		return err
	}

	fileName := fmt.Sprintf("report_%s.json", time.Now().Format("2006-01-02"))
	filePath := filepath.Join(outputDir, fileName)
	log.Printf("Attempting to create report file at: %s", filePath)

	file, err := os.Create(filePath)
	if err != nil {
		log.Printf("Failed to create report file: %v\n", err)
		return err
	}
	defer file.Close()

	if err := json.NewEncoder(file).Encode(report); err != nil {
		log.Printf("Failed to write report to file: %v\n", err)
		return err
	}
	log.Printf("Sales report generated successfully: %s\n", filePath)

	return nil
}

func StartSalesReportBackgroundJob(ctx context.Context, orderStore *InMemoryOrderStore, bookStore *InMemoryBookStore, reportGenerationInterval time.Duration) {
	ticker := time.NewTicker(1 * time.Minute) //24*time.Hour
	defer ticker.Stop()

	log.Println("Starting periodic sales report background job...")

	for {
		select {
		case <-ctx.Done():
			log.Println("Stopping sales report background job.")
			return
		case <-ticker.C:
			log.Println("Triggering sales report generation...")
			err := GenerateSalesReport(ctx, orderStore, bookStore, reportGenerationInterval)
			if err != nil {
				log.Printf("Error generating sales report: %v\n", err)
			}
		}
	}
}
