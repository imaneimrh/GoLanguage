package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	. "FinalProject/models"
	. "FinalProject/stores"
	. "FinalProject/utils"
)

func CreateBookHandler(w http.ResponseWriter, r *http.Request, s *InMemoryBookStore, auth *InMemoryAuthorStore) {
	log.Println("CreateBookHandler: Received request to create a book.")
	var book Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		log.Printf("CreateBookHandler: Invalid request payload. Error: %v\n", err)
		e.RespondWithError(w, http.StatusBadRequest, "Invalid request payload for creating a new book")
		return
	}
	if book.Author.ID != 0 && book.Title != "" && book.Genres != nil && book.PublishedAt != (time.Time{}) && book.Price > 0 && book.Stock > 0 {
		createdBook, err := s.CreateBook(r.Context(), book, auth)
		if err != nil {
			log.Printf("CreateBookHandler: Failed to create book. Error: %v\n", err)
			e.RespondWithError(w, http.StatusInternalServerError, "Failed to create book")
			return
		}
		log.Printf("CreateBookHandler: Book created successfully. ID: %d\n", createdBook.ID)
		e.RespondWithJSON(w, http.StatusCreated, createdBook)
		return
	}
	log.Println("CreateBookHandler: Missing required fields.")
	e.RespondWithError(w, http.StatusBadRequest, "All fields should have a value")
	return
}

func GetBookHandler(w http.ResponseWriter, r *http.Request, s *InMemoryBookStore, auths *InMemoryAuthorStore) {
	log.Println("GetBookHandler: Received request to retrieve a book by ID.")
	bookID, err := ExtractPathParamInt(r)
	if err != nil {
		log.Printf("GetBookHandler: Invalid path parameter. Error: %v\n", err)
		e.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	book, err := s.GetBook(r.Context(), bookID, auths)
	if err != nil {
		log.Printf("GetBookHandler: Book not found. ID: %d. Error: %v\n", bookID, err)
		e.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	log.Printf("GetBookHandler: Book retrieved successfully. ID: %d\n", bookID)
	e.RespondWithJSON(w, http.StatusOK, book)
}

func UpdateBookHandler(w http.ResponseWriter, r *http.Request, s *InMemoryBookStore, auth *InMemoryAuthorStore) {
	log.Println("UpdateBookHandler: Received request to update a book.")
	bookID, err1 := ExtractPathParamInt(r)
	if err1 != nil {
		log.Printf("UpdateBookHandler: Invalid path parameter. Error: %v\n", err1)
		e.RespondWithError(w, http.StatusBadRequest, err1.Error())
		return
	}
	existingBook, err := s.GetBook(r.Context(), bookID, auth)
	if err != nil {
		log.Printf("UpdateBookHandler: Book not found. ID: %d. Error: %v\n", bookID, err)
		e.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	var updatedBook Book
	if err := json.NewDecoder(r.Body).Decode(&updatedBook); err != nil {
		log.Printf("UpdateBookHandler: Invalid request payload. Error: %v\n", err)
		e.RespondWithError(w, http.StatusBadRequest, "Invalid request payload for book Update")
		return
	}

	if updatedBook.Title == "" {
		updatedBook.Title = existingBook.Title
	}
	if updatedBook.Author == (Author{}) {
		updatedBook.Author = existingBook.Author
	}
	if len(updatedBook.Genres) == 0 {
		updatedBook.Genres = existingBook.Genres
	}
	if updatedBook.PublishedAt == (time.Time{}) {
		updatedBook.PublishedAt = existingBook.PublishedAt
	}
	if updatedBook.Price == 0 {
		updatedBook.Price = existingBook.Price
	}
	if updatedBook.Stock == 0 {
		updatedBook.Stock = existingBook.Stock
	}

	b, err := s.UpdateBook(r.Context(), bookID, updatedBook, auth)
	if err != nil {
		log.Printf("UpdateBookHandler: Failed to update book. ID: %d. Error: %v\n", bookID, err)
		e.RespondWithError(w, http.StatusInternalServerError, "Failed to update book")
		return
	}
	log.Printf("UpdateBookHandler: Book updated successfully. ID: %d\n", b.ID)
	e.RespondWithJSON(w, http.StatusOK, b)
}

func DeleteBookHandler(w http.ResponseWriter, r *http.Request, s *InMemoryBookStore) {
	log.Println("DeleteBookHandler: Received request to delete a book.")
	bookID, err1 := ExtractPathParamInt(r)
	if err1 != nil {
		log.Printf("DeleteBookHandler: Invalid path parameter. Error: %v\n", err1)
		e.RespondWithError(w, http.StatusBadRequest, err1.Error())
		return
	}
	err := s.DeleteBook(r.Context(), bookID)
	if err != nil {
		log.Printf("DeleteBookHandler: Failed to delete book. ID: %d. Error: %v\n", bookID, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("DeleteBookHandler: Book deleted successfully. ID: %d\n", bookID)
	e.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func SearchBookHandler(w http.ResponseWriter, r *http.Request, s *InMemoryBookStore) {
	log.Println("SearchBookHandler: Received request to search books.")
	titleParam := r.URL.Query().Get("Title")
	authorParam := r.URL.Query().Get("Author")
	genreParam := r.URL.Query().Get("Genre")
	searchCriteria := SearchCriteria{
		Title:  titleParam,
		Author: authorParam,
		Genre:  genreParam,
	}
	books, err := s.SearchBooks(r.Context(), searchCriteria)
	if err != nil {
		log.Printf("SearchBookHandler: No books found matching criteria. Error: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("SearchBookHandler: Books retrieved successfully.")
	e.RespondWithJSON(w, http.StatusOK, books)
}
