package stores

import (
	. "FinalProject/models"
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

type InMemoryOrderStore struct {
	Mu     sync.RWMutex
	Orders map[int]Order
	NextID int
}
type OrderStore interface {
	CreateOrder(ctx context.Context, order Order) (Order, error)
	GetOrder(ctx context.Context, id int) (Order, error)
	UpdateOrder(ctx context.Context, id int, order Order) (Order, error)
	DeleteOrder(ctx context.Context, id int) error
	ListOrders(ctx context.Context) ([]Order, error)
	ViewOrderHistory(ctx context.Context) (map[int]time.Time, error)
	FetchOrderWithinTimeLimit(ctx context.Context, startTime time.Time, endTime time.Time) ([]Order, error)
	LoadOrders(ctx context.Context, filePath string) error
	SaveOrders(ctx context.Context, filePath string) error
}

func (s *InMemoryOrderStore) CreateOrder(ctx context.Context, order Order, customerStore *InMemoryCustomerStore, bookStore *InMemoryBookStore) (Order, error) {
	select {
	case <-ctx.Done():
		log.Println("Request canceled during Order creation")
		return Order{}, ctx.Err()
	default:
		s.Mu.Lock()
		defer s.Mu.Unlock()

		customer, ok := customerStore.Customers[order.Customer.ID]
		if !ok {
			return Order{}, errors.New("Customer with ID " + strconv.Itoa(order.Customer.ID) + " not found")
		}
		order.Customer = customer

		for i, item := range order.Items {
			book, ok := bookStore.Books[item.Book.ID]
			if !ok {
				return Order{}, errors.New("Book with ID " + strconv.Itoa(item.Book.ID) + " not found")
			}
			if item.Quantity > book.Stock {
				return Order{}, errors.New("Not enough stock for book " + book.Title)
			}
			book.Stock -= item.Quantity
			bookStore.Books[item.Book.ID] = book

			order.Items[i].Book = book
		}

		order.ID = s.NextID
		s.NextID++
		order.CreatedAt = time.Now()
		s.Orders[order.ID] = order

		log.Printf("Order created successfully. ID: %d\n", order.ID)
		return s.Orders[order.ID], nil
	}
}

func (s *InMemoryOrderStore) GetOrder(ctx context.Context, orderId int) (Order, error) {
	select {
	case <-ctx.Done():
		log.Println("Request canceled during Order retrieval of Order" + strconv.Itoa(orderId))
		return Order{}, ctx.Err()
	default:
		s.Mu.Lock()
		defer s.Mu.Unlock()
		order, ok := s.Orders[orderId]
		if !ok {
			log.Println("Order with ID ", orderId, " not found")
			return Order{}, errors.New("Order with ID " + strconv.Itoa(orderId) + " not found")
		}
		return order, nil
	}
}

func (s *InMemoryOrderStore) UpdateOrder(ctx context.Context, orderId int, order Order, customerStore *InMemoryCustomerStore, bookStore *InMemoryBookStore) error {
	select {
	case <-ctx.Done():
		log.Println("Request canceled during Order %d update", orderId)
		return ctx.Err()
	default:
		for _, item := range order.Items {
			for _, book := range bookStore.Books {

				if book.ID != item.Book.ID || item.Quantity != book.Stock {
					return errors.New("Not enough stock for book " + book.Title)
				}
			}
			for _, customer := range customerStore.Customers {
				if customer.ID != order.Customer.ID {
					return errors.New("Customer with id " + strconv.Itoa(order.Customer.ID) + "not found")
				}
			}
		}
		s.Mu.Lock()
		defer s.Mu.Unlock()
		unchangedOrder := s.Orders[order.ID]
		if unchangedOrder.ID != orderId {
			return errors.New("Order with ID " + strconv.Itoa(orderId) + " not found")
		}
		order.ID = unchangedOrder.ID //logically these two fields can't be open for update
		order.Customer = unchangedOrder.Customer
		s.Orders[order.ID] = order
		return nil
	}
}

func (s *InMemoryOrderStore) DeleteOrder(ctx context.Context, OrderId int) error {
	select {
	case <-ctx.Done():
		log.Println("Request canceled during Order %d deletion", OrderId)
		return ctx.Err()
	default:
		s.Mu.Lock()
		defer s.Mu.Unlock()
		if _, ok := s.Orders[OrderId]; ok {
			delete(s.Orders, OrderId)
			return nil
		}
		return errors.New("Order with id " + strconv.Itoa(OrderId) + "not found")
	}
}

func (s *InMemoryOrderStore) ListOrders(ctx context.Context) ([]Order, error) {
	select {
	case <-ctx.Done():
		log.Println("Request canceled during Orders list retrieval")
		return nil, ctx.Err()
	default:
		s.Mu.Lock()
		defer s.Mu.Unlock()
		var orders []Order
		for _, order := range s.Orders {
			orders = append(orders, order)
		}
		return orders, nil
	}
}

func (s *InMemoryOrderStore) ViewOrderHistory(ctx context.Context) (map[int]time.Time, error) {
	select {
	case <-ctx.Done():
		log.Println("Request canceled during Order history retrieval")
		return nil, ctx.Err()
	default:
		history := make(map[int]time.Time)
		for _, order := range s.Orders {
			s.Mu.Lock()
			defer s.Mu.Unlock()
			history[order.ID] = order.CreatedAt
		}
		if len(history) == 0 {
			return nil, errors.New("No orders found")
		}
		return history, nil
	}
}

func (s *InMemoryOrderStore) FetchOrderWithinTimeLimit(ctx context.Context, startTime time.Time, endTime time.Time) ([]Order, error) {
	select {
	case <-ctx.Done():
		log.Println("Request canceled during Order history retrieval")
		return nil, ctx.Err()
	default:
		s.Mu.Lock()
		defer s.Mu.Unlock()

		var orders []Order
		log.Printf("Fetching orders created between %s and %s\n", startTime, endTime)
		for _, order := range s.Orders {
			log.Printf("Checking order ID %d with CreatedAt %s\n", order.ID, order.CreatedAt)
			if !order.CreatedAt.Before(startTime) && !order.CreatedAt.After(endTime) {
				log.Printf("Order ID %d is within the time range (inclusive).\n", order.ID)
				orders = append(orders, order)
			} else {
				log.Printf("Order ID %d is outside the time range.\n", order.ID)
			}
		}

		if len(orders) == 0 {
			log.Println("No orders found within the specified time range.")
		}

		return orders, nil
	}
}

func (s *InMemoryOrderStore) LoadOrders(ctx context.Context, filePath string) error {
	select {
	case <-ctx.Done():
		log.Println("Context canceled during order loading")
		return ctx.Err()
	default:
		dir := "database"
		fullPath := filepath.Join(dir, filePath)

		file, err := os.Open(fullPath)
		if err != nil {
			if os.IsNotExist(err) {
				log.Printf("No existing order database found in %s, starting fresh.\n", fullPath)
				s.NextID = 1
				return nil
			}
			log.Printf("Failed to open file %s: %v\n", fullPath, err)
			return err
		}
		defer file.Close()

		var data struct {
			Orders map[int]Order `json:"orders"`
			NextID int           `json:"next_id"`
		}

		if err := json.NewDecoder(file).Decode(&data); err != nil {
			log.Printf("Failed to decode orders from file %s: %v\n", fullPath, err)
			return err
		}

		s.Mu.Lock()
		defer s.Mu.Unlock()
		s.Orders = data.Orders
		s.NextID = data.NextID
		log.Printf("Orders loaded successfully from %s\n", fullPath)
		return nil
	}
}

func (s *InMemoryOrderStore) SaveOrders(ctx context.Context, filePath string) error {
	select {
	case <-ctx.Done():
		log.Println("Context canceled during order saving")
		return ctx.Err()
	default:
		dir := "database"
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			log.Printf("Failed to create directory %s: %v\n", dir, err)
			return err
		}

		fullPath := filepath.Join(dir, filePath)

		s.Mu.RLock()
		defer s.Mu.RUnlock()

		data := struct {
			Orders map[int]Order `json:"orders"`
			NextID int           `json:"next_id"`
		}{
			Orders: s.Orders,
			NextID: s.NextID,
		}

		file, err := os.Create(fullPath)
		if err != nil {
			log.Printf("Failed to create file %s: %v\n", fullPath, err)
			return err
		}
		defer file.Close()

		if err := json.NewEncoder(file).Encode(data); err != nil {
			log.Printf("Failed to write orders to file %s: %v\n", fullPath, err)
			return err
		}

		log.Printf("Orders saved successfully to %s\n", fullPath)
		return nil
	}
}
