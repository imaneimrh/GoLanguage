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

type InMemoryCustomerStore struct {
	Mu        sync.RWMutex
	Customers map[int]Customer
	NextID    int
}

type CustomerStore interface {
	CreateCustomer(ctx context.Context, customer Customer) (Customer, error)
	GetCustomer(ctx context.Context, id int) (Customer, error)
	UpdateCustomer(ctx context.Context, id int, customer Customer) error
	DeleteCustomer(ctx context.Context, id int, orderStore *InMemoryOrderStore) error
	ListCustomers(ctx context.Context) ([]Customer, error)
	LoadCustomersFromJSON(filePath string) error
	SaveCustomersToJSON(filePath string) error
}

func (s *InMemoryCustomerStore) CreateCustomer(ctx context.Context, customer Customer) (Customer, error) {
	select {
	case <-ctx.Done():
		log.Println("Request canceled during customer creation")
		return Customer{}, ctx.Err()
	default:
		s.Mu.Lock()
		defer s.Mu.Unlock()
		customer.ID = s.NextID
		customer.CreatedAt = time.Now()
		s.NextID++
		s.Customers[customer.ID] = customer
		return s.Customers[customer.ID], nil
	}
}

func (s *InMemoryCustomerStore) GetCustomer(ctx context.Context, customerId int) (Customer, error) {
	select {
	case <-ctx.Done():
		log.Println("Request canceled during customer retrieval:", customerId)
		return Customer{}, ctx.Err()
	default:
		s.Mu.Lock()
		defer s.Mu.Unlock()
		customer, ok := s.Customers[customerId]
		if !ok {
			log.Printf("Customer with ID %d not found", customerId)
			return Customer{}, errors.New("customer with ID " + strconv.Itoa(customerId) + " not found")
		}
		return customer, nil
	}
}

func (s *InMemoryCustomerStore) UpdateCustomer(ctx context.Context, customerId int, customer Customer) error {
	select {
	case <-ctx.Done():
		log.Println("Request canceled during customer update:", customerId)
		return ctx.Err()
	default:
		s.Mu.Lock()
		defer s.Mu.Unlock()
		if _, ok := s.Customers[customerId]; ok {
			unchangedCustomer := s.Customers[customerId]
			customer.CreatedAt = unchangedCustomer.CreatedAt
			customer.ID = customerId
			s.Customers[customerId] = customer
			return nil
		}
		log.Printf("Customer with ID %d not found", customerId)
		return errors.New("customer with ID " + strconv.Itoa(customerId) + " not found")
	}
}

func (s *InMemoryCustomerStore) DeleteCustomer(ctx context.Context, customerId int, orderStore *InMemoryOrderStore) error {
	select {
	case <-ctx.Done():
		log.Printf("Request canceled during deletion of customer ID %d", customerId)
		return ctx.Err()
	default:
		s.Mu.Lock()
		defer s.Mu.Unlock()
		if _, ok := s.Customers[customerId]; ok {
			for _, order := range orderStore.Orders {
				if order.Customer.ID == customerId {
					errMsg := "Cannot delete customer with ID " + strconv.Itoa(customerId) + ", they have an order with ID " + strconv.Itoa(order.ID)
					log.Println(errMsg)
					return errors.New(errMsg)
				}
			}
			delete(s.Customers, customerId)
			return nil
		}
		return errors.New("customer with ID " + strconv.Itoa(customerId) + " not found")
	}
}

func (s *InMemoryCustomerStore) ListCustomers(ctx context.Context) ([]Customer, error) {
	select {
	case <-ctx.Done():
		log.Println("Request canceled during customers list retrieval")
		return nil, ctx.Err()
	default:
		s.Mu.Lock()
		defer s.Mu.Unlock()
		var customers []Customer
		for _, customer := range s.Customers {
			customers = append(customers, customer)
		}
		return customers, nil
	}
}

func (s *InMemoryCustomerStore) LoadCustomers(ctx context.Context, filePath string) error {
	select {
	case <-ctx.Done():
		log.Println("Context canceled during customer loading")
		return ctx.Err()
	default:
		dir := "database"
		fullPath := filepath.Join(dir, filePath)

		file, err := os.Open(fullPath)
		if err != nil {
			if os.IsNotExist(err) {
				log.Printf("No existing customer database found in %s, starting fresh.\n", fullPath)
				s.NextID = 1
				return nil
			}
			log.Printf("Failed to open file %s: %v\n", fullPath, err)
			return err
		}
		defer file.Close()

		var data struct {
			Customers map[int]Customer `json:"customers"`
			NextID    int              `json:"next_id"`
		}

		if err := json.NewDecoder(file).Decode(&data); err != nil {
			log.Printf("Failed to decode customers from file %s: %v\n", fullPath, err)
			return err
		}

		s.Mu.Lock()
		defer s.Mu.Unlock()
		s.Customers = data.Customers
		s.NextID = data.NextID
		log.Printf("Customers loaded successfully from %s\n", fullPath)
		return nil
	}
}

func (s *InMemoryCustomerStore) SaveCustomers(ctx context.Context, filePath string) error {
	select {
	case <-ctx.Done():
		log.Println("Context canceled during customer saving")
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
			Customers map[int]Customer `json:"customers"`
			NextID    int              `json:"next_id"`
		}{
			Customers: s.Customers,
			NextID:    s.NextID,
		}

		file, err := os.Create(fullPath)
		if err != nil {
			log.Printf("Failed to create file %s: %v\n", fullPath, err)
			return err
		}
		defer file.Close()

		if err := json.NewEncoder(file).Encode(data); err != nil {
			log.Printf("Failed to write customers to file %s: %v\n", fullPath, err)
			return err
		}

		log.Printf("Customers saved successfully to %s\n", fullPath)
		return nil
	}
}
