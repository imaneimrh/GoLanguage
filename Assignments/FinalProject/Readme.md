# Project README

## **Overview**
This project is a RESTful web application that manages authors, books, customers, and orders. It includes CRUD operations, real-time inventory management, order processing, and a sales reporting system. The application uses JSON files as a database for persistence and incorporates logging for debugging and monitoring.

---

## **Key Functionalities**

### **1. CRUD Operations**
- **Authors**: Add, update, fetch, and delete authors. Prevent deletion of authors with associated books.
- **Books**: Manage inventory with create, update, fetch, and delete functionality.
- **Customers**: Manage customer data. Prevent deletion of customers with associated orders.
- **Orders**: Place orders that automatically adjust book inventory.

### **2. Inventory Management**
- Placing an order reduces the stock of the ordered books.
- Prevents orders if stock is insufficient.

### **3. Sales Reports**
- Automatically generates sales reports periodically.
- Reports include:
  - Total revenue
  - Total orders
  - Top-selling books
- Allows fetching reports within a specific date range.

### **4. Persistent Data Storage**
- Data for authors, books, customers, and orders is stored in JSON files located in the `database` directory.
- On application startup, data is loaded from these files, ensuring no data loss between restarts.

### **5. Logging**
- Comprehensive logging captures all significant events and errors.
- Logs are written to `bookstore.log` for debugging and monitoring.
  - Examples:
    - Server startup and shutdown
    - Successful and failed API requests
    - Report generation status
    - File I/O operations (e.g., saving/loading JSON files)

---

## **System Flow**
1. **Populate Database**:
   - Start by creating authors, then books, customers, and finally orders.
2. **Real-Time Processing**:
   - Placing an order updates the book stock.
   - Reports are generated automatically and stored in the `output-reports` directory.
3. **Restart and Reload**:
   - After stopping the application, confirm that all data is saved in the JSON files.
   - Upon restarting, the application reloads data seamlessly from the JSON files.
4. **Cleanup**:
   - Delete authors, books, customers, and orders, following constraints (e.g., authors with books cannot be deleted).

---

## **Special Features**

### **1. JSON Database**
- All entities (authors, books, customers, orders) are stored in JSON files.
- Example file structure:
  ```json
  {
    "orders": {
      "1": {
        "id": 1,
        "customer": {
          "id": 1,
          "name": "John Smith"
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
        "created_at": "2025-01-12T13:13:07Z"
      }
    },
    "next_id": 2
  }
  ```

### **2. Logging**
- The application logs events and errors to `bookstore.log`.
- Example log entries:
  ```
  2025/01/12 13:13:07 Server running on http://localhost:8080
  2025/01/12 13:14:00 Sales report generated successfully: output-reports/report_20250112.json
  2025/01/12 13:15:23 Error generating sales report: No orders found for the specified time range.
  ```
- Logs assist in debugging and provide insight into application behavior.

### **3. Automated Report Generation**
- Reports are generated periodically and saved in the `output-reports` directory.
- Filename format: `report_YYYY-MM-DD.json`
- Reports include:
  ```json
  {
    "timestamp": "2025-01-12T14:48:17Z",
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
  ```

---

## **Instructions for the Professor**
I have included a testing file with example json data to use when working with postman
1. **Start the Application**:
   - Run:
     ```bash
     go run main.go
     ```

2. **Test CRUD Operations**:
   - Create authors, books, customers, and orders.
   - Update and fetch entities.
   - Verify that placing an order reduces book stock.

3. **Verify Report Generation** (after 1 minute for the sake of testing, but can be changed to 24 hours)
   - Check the `output-reports` directory for generated reports.
   - Use the `GET /reports` API to fetch reports by date range.

4. **Test Data Persistence**:
   - Stop the application and verify that JSON files in the `database` directory contain saved data.
   - Restart the application and confirm that data is reloaded.

5. **Perform Deletions**:
   - Attempt to delete authors, books, customers, and orders.
   - Test special cases where constraints prevent deletion (e.g., authors with books).

---

Thank you for reviewing this work!
I tried to give an overview of how it functions but, Please, feel free to reach out in case you have a question on how to use it.

