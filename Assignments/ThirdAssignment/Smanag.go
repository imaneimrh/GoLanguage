package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Address struct {
	Street     string
	City       string
	State      string
	PostalCode string
}

type Course struct {
	Code   int
	Name   string
	Credit int
}

type Student struct {
	ID int
	StudentInput
}

type StudentInput struct {
	Name    string
	Age     int
	Major   string
	Address Address
	Courses []Course
}

type Response interface {
	respondWithJSON(w http.ResponseWriter, status int, payload interface{})
	respondWithError(w http.ResponseWriter, code int, message string)
}

type ErrorResponse struct {
	error string
}

var Students []Student
var e ErrorResponse

func (e ErrorResponse) respondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}
func (e ErrorResponse) respondWithError(w http.ResponseWriter, code int, message string) {
	e.respondWithJSON(w, code, map[string]string{"error": message})
}

func getStudents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Students)
}

func CreateStudent(w http.ResponseWriter, r *http.Request) {
	var Student Student
	err := json.NewDecoder(r.Body).Decode(&Student)
	if err != nil {
		e.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	//Student.ID = len(Students) + 1
	if Student.Name == "" || Student.Major == "" || Student.Age == 0 || Student.Address.Street == "" || Student.Courses == nil || Student.Address.City == "" || Student.Address.State == "" || Student.Address.PostalCode == "" {
		e.respondWithError(w, http.StatusBadRequest, "All fileds should have a value")
		return
	}
	Students = append(Students, Student)

}

func GetStudentByID(w http.ResponseWriter, r *http.Request) {

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 || parts[2] == "" {
		http.Error(w, "Student ID not provided", http.StatusBadRequest)
		return
	}
	StudentID := parts[2]

	response := map[string]string{
		"message": "Student fetched successfully",
		"user_id": StudentID,
	}
	e.respondWithJSON(w, http.StatusOK, response)
}

func UpdateStudentByID(w http.ResponseWriter, r *http.Request) {
	var student Student
	//var s []Student
	err := json.NewDecoder(r.Body).Decode(&student)
	if err != nil {
		e.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 || parts[2] == "" {
		http.Error(w, "Student ID not provided", http.StatusBadRequest)
		return
	}
	student.ID, _ = strconv.Atoi(parts[2])
	for index, existingStudent := range Students {
		if existingStudent.ID == student.ID {
			Students = append(Students[:index], Students[index+1:]...)
			Students = append(Students, student)
			e.respondWithJSON(w, http.StatusOK, student)
			return
		}
	}
	e.respondWithError(w, http.StatusNotFound, "Student not found")
}

func DeleteStudentByID(w http.ResponseWriter, r *http.Request) {
	var student Student
	err := json.NewDecoder(r.Body).Decode(&student)
	if err != nil {
		e.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 || parts[2] == "" {
		http.Error(w, "Student ID not provided", http.StatusBadRequest)
		return
	}
	student.ID, _ = strconv.Atoi(parts[2])
	for index, existingStudent := range Students {
		if existingStudent.ID == student.ID {
			Students = append(Students[:index], Students[index+1:]...)
			e.respondWithJSON(w, http.StatusOK, student)
			return
		}
	}
	e.respondWithError(w, http.StatusNotFound, "Student not found")
}

func main() {
	s := Student{
		ID: 1,
		StudentInput: StudentInput{
			Name:  "Imane Imrharn",
			Age:   21,
			Major: "Computer Science",
			Address: Address{
				Street:     "123 Main St",
				City:       "City",
				State:      "State",
				PostalCode: "12345",
			},
			Courses: []Course{
				{
					Code:   1,
					Name:   "Course 1",
					Credit: 3,
				},
				{
					Code:   2,
					Name:   "Course 2",
					Credit: 4,
				},
			},
		},
	}
	Students = append(Students, s)
	http.HandleFunc("/students", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getStudents(w, r)
		case http.MethodPost:
			CreateStudent(w, r)
		default:
			http.Error(w, "Method not supported", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/students/{id}", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			GetStudentByID(w, r)
		case http.MethodPut:
			UpdateStudentByID(w, r)
		case http.MethodDelete:
			DeleteStudentByID(w, r)
			return
		default:
			http.Error(w, "Method not supported", http.StatusMethodNotAllowed)
		}
	})
	log.Println("Server started on port 8080")
	http.ListenAndServe(":8080", nil)

}
