package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	. "FinalProject/models"
	. "FinalProject/stores"
	. "FinalProject/utils"
)

func CreateCustomerHandler(w http.ResponseWriter, r *http.Request, c *InMemoryCustomerStore) {
	log.Println("CreateCustomerHandler: Received request to create a customer.")
	var customer Customer
	if err := json.NewDecoder(r.Body).Decode(&customer); err != nil {
		log.Printf("CreateCustomerHandler: Invalid request payload. Error: %v\n", err)
		e.RespondWithError(w, http.StatusBadRequest, "Invalid request payload for creating a new customer")
		return
	}
	if customer.Name != "" && customer.Email != "" && customer.Address != (Address{}) {
		newCustomer, err := c.CreateCustomer(r.Context(), customer)
		if err != nil {
			log.Printf("CreateCustomerHandler: Failed to create customer. Error: %v\n", err)
			e.RespondWithError(w, http.StatusInternalServerError, "Failed to create customer")
			return
		}
		log.Printf("CreateCustomerHandler: Customer created successfully. ID: %d\n", newCustomer.ID)
		e.RespondWithJSON(w, http.StatusCreated, newCustomer)
		return
	}
	log.Println("CreateCustomerHandler: Missing required fields.")
	e.RespondWithError(w, http.StatusBadRequest, "All fields should have a value")
}

func GetCustomerByIDHandler(w http.ResponseWriter, r *http.Request, c *InMemoryCustomerStore) {
	log.Println("GetCustomerByIDHandler: Received request to retrieve a customer by ID.")
	customerID, err := ExtractPathParamInt(r)
	if err != nil {
		log.Printf("GetCustomerByIDHandler: Invalid path parameter. Error: %v\n", err)
		e.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	customer, err := c.GetCustomer(r.Context(), customerID)
	if err != nil {
		log.Printf("GetCustomerByIDHandler: Customer not found. ID: %d. Error: %v\n", customerID, err)
		e.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	log.Printf("GetCustomerByIDHandler: Customer retrieved successfully. ID: %d\n", customerID)
	e.RespondWithJSON(w, http.StatusOK, customer)
}

func GetAllCustomersHandler(w http.ResponseWriter, r *http.Request, c *InMemoryCustomerStore) {
	log.Println("GetAllCustomersHandler: Received request to list all customers.")
	customers, err := c.ListCustomers(r.Context())
	if err != nil {
		log.Printf("GetAllCustomersHandler: Failed to retrieve customers. Error: %v\n", err)
		e.RespondWithError(w, http.StatusInternalServerError, "Failed to get customers")
		return
	}
	log.Println("GetAllCustomersHandler: Customers retrieved successfully.")
	e.RespondWithJSON(w, http.StatusOK, customers)
}

func UpdateCustomerHandler(w http.ResponseWriter, r *http.Request, c *InMemoryCustomerStore) {
	log.Println("UpdateCustomerHandler: Received request to update a customer.")
	customerID, err := ExtractPathParamInt(r)
	if err != nil {
		log.Printf("UpdateCustomerHandler: Invalid path parameter. Error: %v\n", err)
		e.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	existingCustomer, err := c.GetCustomer(r.Context(), customerID)
	if err != nil {
		log.Printf("UpdateCustomerHandler: Customer not found. ID: %d. Error: %v\n", customerID, err)
		e.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	var updatedCustomer Customer
	if err := json.NewDecoder(r.Body).Decode(&updatedCustomer); err != nil {
		log.Printf("UpdateCustomerHandler: Invalid request payload. Error: %v\n", err)
		e.RespondWithError(w, http.StatusBadRequest, "Invalid request payload for customer Update")
		return
	}

	if updatedCustomer.Name == "" {
		updatedCustomer.Name = existingCustomer.Name
	}
	if updatedCustomer.Email == "" {
		updatedCustomer.Email = existingCustomer.Email
	}
	if updatedCustomer.Address == (Address{}) {
		updatedCustomer.Address = existingCustomer.Address
	}

	err = c.UpdateCustomer(r.Context(), customerID, updatedCustomer)
	if err != nil {
		log.Printf("UpdateCustomerHandler: Failed to update customer. ID: %d. Error: %v\n", customerID, err)
		e.RespondWithError(w, http.StatusInternalServerError, "Failed to update customer")
		return
	}
	log.Printf("UpdateCustomerHandler: Customer updated successfully. ID: %d\n", customerID)
	e.RespondWithJSON(w, http.StatusOK, updatedCustomer)
}

func DeleteCustomerHandler(w http.ResponseWriter, r *http.Request, c *InMemoryCustomerStore, o *InMemoryOrderStore) {
	log.Println("DeleteCustomerHandler: Received request to delete a customer.")
	customerID, err := ExtractPathParamInt(r)
	if err != nil {
		log.Printf("DeleteCustomerHandler: Invalid path parameter. Error: %v\n", err)
		e.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	err = c.DeleteCustomer(r.Context(), customerID, o)
	if err != nil {
		log.Printf("DeleteCustomerHandler: Failed to delete customer. ID: %d. Error: %v\n", customerID, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("DeleteCustomerHandler: Customer deleted successfully. ID: %d\n", customerID)
	e.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}
