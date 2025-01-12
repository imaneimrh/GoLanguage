# Comprehensive Walkthrough for Testing with Postman

---

## **Step 1: Start the Application**
1. Run the application:
   ```bash
   go run main.go
   ```
2. Confirm the server is running:
   - URL: `http://localhost:8080`

---

## **Step 2: Populate the Database**

### **2.1 Create Authors**
- **Endpoint**: `POST http://localhost:8080/authors`
- **Request Body**:
  ```json
  {
    "first_name": "Imane",
    "last_name": "Imrh",
    "bio": "A Go programmer."
  }
  ```
  ```json
  {
    "first_name": "Jane",
    "last_name": "Doe",
    "bio": "A writer specializing in advanced topics."
  }
  ```
- **Expected Responses**:
  ```json
  {
    "id": 1,
    "first_name": "Imane",
    "last_name": "Imrh",
    "bio": "A Go programmer."
  }
  ```
  ```json
  {
    "id": 2,
    "first_name": "Jane",
    "last_name": "Doe",
    "bio": "A writer specializing in advanced topics."
  }
  ```

---

### **2.2 Create Books**
- **Endpoint**: `POST http://localhost:8080/books`
- **Request Body**:
  ```json
  {
    "title": "Go Programming",
    "author": {
      "id": 1
    },
    "genres": ["Programming", "Technology"],
    "published_at": "2023-01-01T00:00:00Z",
    "price": 49.99,
    "stock": 100
  }
  ```
  ```json
  {
    "title": "Advanced Go",
    "author": {
      "id": 2
    },
    "genres": ["Programming"],
    "published_at": "2023-06-01T00:00:00Z",
    "price": 59.99,
    "stock": 80
  }
  ```
- **Expected Responses**:
  ```json
  {
    "id": 1,
    "title": "Go Programming",
    "author": {
      "id": 1,
      "first_name": "Imane",
      "last_name": "Imrh",
      "bio": "A Go programmer."
    },
    "genres": ["Programming", "Technology"],
    "published_at": "2023-01-01T00:00:00Z",
    "price": 49.99,
    "stock": 100
  }
  ```

---

### **2.3 Create Customers**
- **Endpoint**: `POST http://localhost:8080/customers`
- **Request Body**:
  ```json
  {
    "name": "John Smith",
    "email": "john.smith@example.com",
    "address": {
      "street": "123 Main St",
      "city": "New York",
      "state": "NY",
      "postal_code": "10001",
      "country": "USA"
    }
  }
  ```
  ```json
  {
    "name": "Alice Johnson",
    "email": "alice.johnson@example.com",
    "address": {
      "street": "456 Maple Ave",
      "city": "Los Angeles",
      "state": "CA",
      "postal_code": "90001",
      "country": "USA"
    }
  }
  ```
- **Expected Responses**:
  ```json
  {
    "id": 1,
    "name": "John Smith",
    "email": "john.smith@example.com",
    "address": {
      "street": "123 Main St",
      "city": "New York",
      "state": "NY",
      "postal_code": "10001",
      "country": "USA"
    }
  }
  ```

---

### **2.4 Place Orders**
- **Endpoint**: `POST http://localhost:8080/orders`
- **Request Body**:
  ```json
  {
    "customer": {
      "id": 1
    },
    "items": [
      {
        "book": {
          "id": 1
        },
        "quantity": 2
      }
    ],
    "total_price": 99.98,
    "status": "Completed"
  }
  ```
  ```json
  {
    "customer": {
      "id": 2
    },
    "items": [
      {
        "book": {
          "id": 2
        },
        "quantity": 1
      }
    ],
    "total_price": 59.99,
    "status": "Completed"
  }
  ```
- **Expected Responses**:
  ```json
  {
    "id": 1,
    "customer": {
      "id": 1,
      "name": "John Smith",
      "email": "john.smith@example.com"
    },
    "items": [
      {
        "book": {
          "id": 1,
          "title": "Go Programming"
        },
        "quantity": 2
      }
    ],
    "total_price": 99.98,
    "status": "Completed"
  }
  ```
---
### **3.2 Test Stock Reduction**
- **Endpoint**: `GET http://localhost:8080/books/1`
- **Expected Response (After Order)**:
  ```json
  {
    "id": 1,
    "title": "Go Programming",
    "stock": 98
  }
  ```

---

## **Step 3: Update Tests**

### **3.1 Update an Author**
- **Endpoint**: `PUT http://localhost:8080/authors/1`
- **Request Body**:
  ```json
  {
    "first_name": "Imane",
    "last_name": "Imrh",
    "bio": "A Go expert and programming teacher."
  }
  ```
- **Expected Response**:
  ```json
  {
    "id": 1,
    "first_name": "Imane",
    "last_name": "Imrh",
    "bio": "A Go expert and programming teacher."
  }
  ```

---


## **Step 4: Search for Books**
- **Endpoint**: `GET http://localhost:8080/books?title=Go`
- **Expected Response**:
  ```json
  [
    {
      "id": 1,
      "title": "Go Programming"
    },
    {
      "id": 2,
      "title": "Advanced Go"
    }
  ]
  ```

---

## **Step 5: List Reports by Time Range**
- **Endpoint**: `GET http://localhost:8080/reports?start_date=2025-01-12&end_date=2025-01-12`
- **Expected Response**:
  ```json
  {
    "reports": [
      {
        "timestamp": "2025-01-12T14:48:17.4369905+01:00",
        "total_revenue": 159.97,
        "total_orders": 2,
        "top_selling_books": [
          {
            "book": {
              "id": 1,
              "title": "Go Programming"
            },
            "quantity_sold": 2
          }
        ]
      }
    ]
  }
  ```

---

## **Step 6: Save and Reload Data**
1. Ask the professor to stop the application and confirm that data is saved in the JSON files within the `database` directory.
2. Restart the application:
   ```bash
   go run main.go
   ```
3. Perform `GET` requests for authors, books, customers, and orders to confirm the data is reloaded correctly.

---

## **Step 7: Delete Tests**

### **7.1 Delete Authors**
- **Endpoint**: `DELETE http://localhost:8080/authors/1`
- **Expected Response**:
  ```json
  {
    "result": "success"
  }
  ```
- **Special Case**: Deleting an author with associated books
  - Expected Response:
    ```json
    {
      "error": "Cannot delete author with ID 1, as they have associated books."
    }
    ```

---

### **7.2 Delete Customers**
- **Endpoint**: `DELETE http://localhost:8080/customers/1`
- **Expected Response**:
  ```json
  {
    "result": "success"
  }
  ```
- **Special Case**: Deleting a customer with associated orders
  - Expected Response:
    ```json
    {
      "error": "Cannot delete customer with ID 1, as they have associated orders."
    }
    ```

---

### **7.3 Delete Orders**
- **Endpoint**: `DELETE http://localhost:8080/orders/1`
- **Expected Response**:
  ```json
  {
    "result": "success"
    }
 

