package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	. "FinalProject/models"
	. "FinalProject/stores"
	. "FinalProject/utils"
)

func CreateAuthorHandler(w http.ResponseWriter, r *http.Request, auth *InMemoryAuthorStore) {
	log.Println("CreateAuthorHandler: Received request to create an author.")
	var author Author
	if err := json.NewDecoder(r.Body).Decode(&author); err != nil {
		log.Printf("CreateAuthorHandler: Invalid request payload. Error: %v\n", err)
		e.RespondWithError(w, http.StatusBadRequest, "Invalid request payload for creating a new author")
		return
	}
	if author.FirstName != "" && author.LastName != "" {
		createdAuthor, err := auth.CreateAuthor(r.Context(), author)
		if err != nil {
			log.Printf("CreateAuthorHandler: Failed to create author. Error: %v\n", err)
			e.RespondWithError(w, http.StatusInternalServerError, "Failed to create author")
			return
		}
		log.Printf("CreateAuthorHandler: Author created successfully. ID: %d\n", createdAuthor.ID)
		e.RespondWithJSON(w, http.StatusCreated, createdAuthor)
		return
	}
	log.Println("CreateAuthorHandler: Missing required fields.")
	e.RespondWithError(w, http.StatusBadRequest, "All fields should have a value")
}

func GetAuthorByIdHandler(w http.ResponseWriter, r *http.Request, auth *InMemoryAuthorStore) {
	log.Println("GetAuthorByIdHandler: Received request to retrieve an author by ID.")
	authorID, err := ExtractPathParamInt(r)
	if err != nil {
		log.Printf("GetAuthorByIdHandler: Invalid path parameter. Error: %v\n", err)
		e.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	author, err := auth.GetAuthor(r.Context(), authorID)
	if err != nil {
		log.Printf("GetAuthorByIdHandler: Author not found. ID: %d. Error: %v\n", authorID, err)
		e.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	log.Printf("GetAuthorByIdHandler: Author retrieved successfully. ID: %d\n", authorID)
	e.RespondWithJSON(w, http.StatusOK, author)
}

func UpdateAuthorHandler(w http.ResponseWriter, r *http.Request, auth *InMemoryAuthorStore) {
	log.Println("UpdateAuthorHandler: Received request to update an author.")
	authorID, err := ExtractPathParamInt(r)
	if err != nil {
		log.Printf("UpdateAuthorHandler: Invalid path parameter. Error: %v\n", err)
		e.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	existingAuthor, err := auth.GetAuthor(r.Context(), authorID)
	if err != nil {
		log.Printf("UpdateAuthorHandler: Author not found. ID: %d. Error: %v\n", authorID, err)
		e.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	var updatedAuthor Author
	if err := json.NewDecoder(r.Body).Decode(&updatedAuthor); err != nil {
		log.Printf("UpdateAuthorHandler: Invalid request payload. Error: %v\n", err)
		e.RespondWithError(w, http.StatusBadRequest, "Invalid request payload for author Update")
		return
	}

	if updatedAuthor.FirstName == "" {
		updatedAuthor.FirstName = existingAuthor.FirstName
	}
	if updatedAuthor.LastName == "" {
		updatedAuthor.LastName = existingAuthor.LastName
	}
	if updatedAuthor.Bio == "" {
		updatedAuthor.Bio = existingAuthor.Bio
	}

	updatedAuthor, err = auth.UpdateAuthor(r.Context(), authorID, updatedAuthor)
	if err != nil {
		log.Printf("UpdateAuthorHandler: Failed to update author. ID: %d. Error: %v\n", authorID, err)
		e.RespondWithError(w, http.StatusInternalServerError, "Failed to update author")
		return
	}
	log.Printf("UpdateAuthorHandler: Author updated successfully. ID: %d\n", updatedAuthor.ID)
	e.RespondWithJSON(w, http.StatusOK, updatedAuthor)
}

func DeleteAuthorHandler(w http.ResponseWriter, r *http.Request, auth *InMemoryAuthorStore, s *InMemoryBookStore) {
	log.Println("DeleteAuthorHandler: Received request to delete an author.")
	authorID, err := ExtractPathParamInt(r)
	if err != nil {
		log.Printf("DeleteAuthorHandler: Invalid path parameter. Error: %v\n", err)
		e.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	err = auth.DeleteAuthor(r.Context(), authorID, s)
	if err != nil {
		log.Printf("DeleteAuthorHandler: Failed to delete author. ID: %d. Error: %v\n", authorID, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("DeleteAuthorHandler: Author deleted successfully. ID: %d\n", authorID)
	e.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func ListAllHandler(w http.ResponseWriter, r *http.Request, auth *InMemoryAuthorStore) {
	log.Println("ListAllHandler: Received request to list all authors.")
	authors, err := auth.ListAuthors(r.Context())
	if err != nil {
		log.Printf("ListAllHandler: Failed to list authors. Error: %v\n", err)
		e.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	log.Println("ListAllHandler: Authors retrieved successfully.")
	e.RespondWithJSON(w, http.StatusOK, authors)
}
