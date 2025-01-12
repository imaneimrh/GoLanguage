package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	. "FinalProject/models"
	. "FinalProject/stores"
	. "FinalProject/utils"
)

var e ErrorResponse

func CreateOrderHandler(w http.ResponseWriter, r *http.Request, customerStore *InMemoryCustomerStore, bookStore *InMemoryBookStore, orderStore *InMemoryOrderStore) {
	log.Println("CreateOrderHandler: Received request to create an order.")
	var order Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		log.Printf("CreateOrderHandler: Invalid request payload. Error: %v\n", err)
		e.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	ctx := r.Context()
	createdOrder, err := orderStore.CreateOrder(ctx, order, customerStore, bookStore)
	if err != nil {
		log.Printf("CreateOrderHandler: Failed to create order. Error: %v\n", err)
		e.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	log.Printf("CreateOrderHandler: Order created successfully. ID: %d\n", createdOrder.ID)
	e.RespondWithJSON(w, http.StatusCreated, createdOrder)
}

func GetOrderHandler(w http.ResponseWriter, r *http.Request, orderStore *InMemoryOrderStore) {
	log.Println("GetOrderHandler: Received request to retrieve an order by ID.")
	orderID, err := ExtractPathParamInt(r)
	if err != nil {
		log.Printf("GetOrderHandler: Invalid path parameter. Error: %v\n", err)
		e.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	order, err := orderStore.GetOrder(r.Context(), orderID)
	if err != nil {
		log.Printf("GetOrderHandler: Order not found. ID: %d. Error: %v\n", orderID, err)
		e.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	log.Printf("GetOrderHandler: Order retrieved successfully. ID: %d\n", orderID)
	e.RespondWithJSON(w, http.StatusOK, order)
}

func GetAllOrdersHandler(w http.ResponseWriter, r *http.Request, orderStore *InMemoryOrderStore) {
	log.Println("GetAllOrdersHandler: Received request to list all orders.")
	orders, err := orderStore.ListOrders(r.Context())
	if err != nil {
		log.Printf("GetAllOrdersHandler: Failed to retrieve orders. Error: %v\n", err)
		e.RespondWithError(w, http.StatusInternalServerError, "Failed to get orders")
		return
	}
	log.Println("GetAllOrdersHandler: Orders retrieved successfully.")
	e.RespondWithJSON(w, http.StatusOK, orders)
}

func UpdateOrderHandler(w http.ResponseWriter, r *http.Request, orderStore *InMemoryOrderStore, customerStore *InMemoryCustomerStore, bookStore *InMemoryBookStore) {
	log.Println("UpdateOrderHandler: Received request to update an order.")
	orderID, err := ExtractPathParamInt(r)
	if err != nil {
		log.Printf("UpdateOrderHandler: Invalid path parameter. Error: %v\n", err)
		e.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	var updatedOrder Order
	if err := json.NewDecoder(r.Body).Decode(&updatedOrder); err != nil {
		log.Printf("UpdateOrderHandler: Invalid request payload. Error: %v\n", err)
		e.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	ctx := r.Context()
	err = orderStore.UpdateOrder(ctx, orderID, updatedOrder, customerStore, bookStore)
	if err != nil {
		log.Printf("UpdateOrderHandler: Failed to update order. ID: %d. Error: %v\n", orderID, err)
		e.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	log.Printf("UpdateOrderHandler: Order updated successfully. ID: %d\n", orderID)
	e.RespondWithJSON(w, http.StatusOK, updatedOrder)
}

func DeleteOrderHandler(w http.ResponseWriter, r *http.Request, orderStore *InMemoryOrderStore) {
	log.Println("DeleteOrderHandler: Received request to delete an order.")
	orderID, err := ExtractPathParamInt(r)
	if err != nil {
		log.Printf("DeleteOrderHandler: Invalid path parameter. Error: %v\n", err)
		e.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	err = orderStore.DeleteOrder(r.Context(), orderID)
	if err != nil {
		log.Printf("DeleteOrderHandler: Failed to delete order. ID: %d. Error: %v\n", orderID, err)
		e.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	log.Printf("DeleteOrderHandler: Order deleted successfully. ID: %d\n", orderID)
	e.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func FetchOrdersWithinTimeHandler(w http.ResponseWriter, r *http.Request, store *InMemoryOrderStore) {
	log.Println("FetchOrdersWithinTimeHandler: Received request to fetch orders within a time range.")
	startTimeStr := r.URL.Query().Get("startTime")
	endTimeStr := r.URL.Query().Get("endTime")

	startTime, err := time.Parse(time.RFC3339, startTimeStr)
	if err != nil {
		log.Printf("FetchOrdersWithinTimeHandler: Invalid start time format. Error: %v\n", err)
		e.RespondWithError(w, http.StatusBadRequest, "Invalid start time format")
		return
	}
	endTime, err := time.Parse(time.RFC3339, endTimeStr)
	if err != nil {
		log.Printf("FetchOrdersWithinTimeHandler: Invalid end time format. Error: %v\n", err)
		e.RespondWithError(w, http.StatusBadRequest, "Invalid end time format")
		return
	}

	ctx := r.Context()
	orders, err := store.FetchOrderWithinTimeLimit(ctx, startTime, endTime)
	if err != nil {
		log.Printf("FetchOrdersWithinTimeHandler: Failed to fetch orders. Error: %v\n", err)
		e.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	log.Println("FetchOrdersWithinTimeHandler: Orders retrieved successfully.")
	e.RespondWithJSON(w, http.StatusOK, orders)
}

func ListReportsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("ListReportsHandler: Received request to list reports within a date range.")

	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")

	startTime, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		log.Printf("ListReportsHandler: Invalid start_date format. Provided: %s. Error: %v\n", startDate, err)
		e.RespondWithError(w, http.StatusBadRequest, "Invalid start_date format. Use YYYY-MM-DD")
		return
	}

	endTime, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		log.Printf("ListReportsHandler: Invalid end_date format. Provided: %s. Error: %v\n", endDate, err)
		e.RespondWithError(w, http.StatusBadRequest, "Invalid end_date format. Use YYYY-MM-DD")
		return
	}

	if endTime.Before(startTime) {
		log.Printf("ListReportsHandler: end_date is before start_date. start_date: %s, end_date: %s\n", startDate, endDate)
		e.RespondWithError(w, http.StatusBadRequest, "end_date must be after start_date")
		return
	}

	outputDir := "output-reports"
	files, err := os.ReadDir(outputDir)
	if err != nil {
		log.Printf("ListReportsHandler: Failed to read reports directory. Path: %s. Error: %v\n", outputDir, err)
		e.RespondWithError(w, http.StatusInternalServerError, "Failed to read reports directory")
		return
	}

	var reports []SalesReport
	for _, file := range files {
		if strings.HasPrefix(file.Name(), "report_") {
			reportDateStr := strings.TrimSuffix(strings.TrimPrefix(file.Name(), "report_"), ".json")
			reportDate, err := time.Parse("2006-01-02", reportDateStr)
			if err != nil {
				log.Printf("Failed to parse date from file: %s. Error: %v\n", file.Name(), err)
				continue
			}
			if !reportDate.Before(startTime) && !reportDate.After(endTime) {
				reportFile, err := os.Open(filepath.Join(outputDir, file.Name()))
				if err != nil {
					log.Printf("ListReportsHandler: Failed to open report file. File: %s. Error: %v\n", file.Name(), err)
					continue
				}
				defer reportFile.Close()

				var report SalesReport
				if err := json.NewDecoder(reportFile).Decode(&report); err != nil {
					log.Printf("ListReportsHandler: Failed to decode report file. File: %s. Error: %v\n", file.Name(), err)
					continue
				}
				reports = append(reports, report)
			}
		}
	}

	if len(reports) == 0 {
		log.Printf("ListReportsHandler: No reports found in the specified date range. start_date: %s, end_date: %s\n", startDate, endDate)
		e.RespondWithError(w, http.StatusNotFound, "No reports found for the given date range")
		return
	}

	log.Printf("ListReportsHandler: Reports retrieved successfully. Count: %d\n", len(reports))
	e.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
		"reports": reports,
	})
}

func ViewOrderHistoryHandler(w http.ResponseWriter, r *http.Request, store *InMemoryOrderStore) {
	log.Println("ViewOrderHistoryHandler: Received request to view order history.")

	ctx := r.Context()
	history, err := store.ViewOrderHistory(ctx)
	if err != nil {
		log.Printf("ViewOrderHistoryHandler: Failed to retrieve order history. Error: %v\n", err)
		e.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve order history")
		return
	}

	if len(history) == 0 {
		log.Println("ViewOrderHistoryHandler: No order history found.")
		e.RespondWithError(w, http.StatusNotFound, "No order history found")
		return
	}

	log.Printf("ViewOrderHistoryHandler: Order history retrieved successfully. Count: %d\n", len(history))
	e.RespondWithJSON(w, http.StatusOK, history)
}
