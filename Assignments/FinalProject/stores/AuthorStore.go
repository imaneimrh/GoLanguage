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
)

// ----------------------------------------------Definition of AuthorMethods--------------------------------
type InMemoryAuthorStore struct {
	Mu      sync.RWMutex
	Authors map[int]Author
	NextID  int
}
type AuthorStore interface {
	CreateAuthor(ctx context.Context, author Author) (Author, error)
	GetAuthor(ctx context.Context, id int) (Author, error)
	UpdateAuthor(ctx context.Context, id int, author Author) (Author, error)
	DeleteAuthor(ctx context.Context, id int) error
	ListAuthors(ctx context.Context) ([]Author, error)
	LoadAuthors(filePath string) error
	SaveAuthors(filePath string) error
}

func (s *InMemoryAuthorStore) CreateAuthor(ctx context.Context, author Author) (Author, error) {
	select {
	case <-ctx.Done():
		log.Println("Request canceled during Author creation")
		return Author{}, ctx.Err()
	default:
		s.Mu.Lock()
		defer s.Mu.Unlock()
		author.ID = s.NextID
		s.NextID++
		s.Authors[author.ID] = author
		return author, nil
	}
}

func (s *InMemoryAuthorStore) GetAuthor(ctx context.Context, authorId int) (Author, error) {
	select {
	case <-ctx.Done():
		log.Println("Request canceled during book retrieval of Author")
		return Author{}, ctx.Err()
	default:
		s.Mu.Lock()
		defer s.Mu.Unlock()
		author, ok := s.Authors[authorId]
		if !ok {
			log.Println("Author with ID ", authorId, " not found")
			return Author{}, errors.New("Author with ID " + strconv.Itoa(authorId) + " not found")
		}
		return author, nil
	}
}

func (s *InMemoryAuthorStore) UpdateAuthor(ctx context.Context, authorId int, author Author) (Author, error) {
	select {
	case <-ctx.Done():
		log.Println("Request canceled during Author update")
		return author, ctx.Err()
	default:
		s.Mu.Lock()
		defer s.Mu.Unlock()
		if _, ok := s.Authors[authorId]; ok {
			unchangedAuthor := s.Authors[authorId]
			author.ID = unchangedAuthor.ID
			s.Authors[authorId] = author
			return author, nil
		}
		return Author{}, errors.New("Author with id " + strconv.Itoa(authorId) + "not found")
	}
}

func (s *InMemoryAuthorStore) DeleteAuthor(ctx context.Context, authId int, b *InMemoryBookStore) error {
	select {
	case <-ctx.Done():
		log.Println("Request canceled during Author deletion")
		return ctx.Err()
	default:
		s.Mu.Lock()
		defer s.Mu.Unlock()
		if _, ok := s.Authors[authId]; ok {
			for _, book := range b.Books {
				if book.Author.ID == authId {
					log.Println("You are trying to delete an author that is the author of a book with id ", book.ID, ". Please delete the books related to this author first.")
					return errors.New("You are trying to delete an author that is the author of a book with id " + strconv.Itoa(book.ID) + ". Please delete the books related to this author first.")
				}
			}
			delete(s.Authors, authId)
			return nil
		}
		return errors.New("Author with id " + strconv.Itoa(authId) + "not found")
	}
}

func (s *InMemoryAuthorStore) ListAuthors(ctx context.Context) ([]Author, error) {
	select {
	case <-ctx.Done():
		log.Println("Request canceled during Author list retrieval")
		return nil, ctx.Err()
	default:
		s.Mu.Lock()
		defer s.Mu.Unlock()
		var authors []Author
		for _, author := range s.Authors {
			authors = append(authors, author)
		}
		return authors, nil
	}
}

func (s *InMemoryAuthorStore) LoadAuthors(ctx context.Context, filePath string) error {
	select {
	case <-ctx.Done():
		log.Println("Context canceled during author loading")
		return ctx.Err()
	default:
		dir := "database"
		fullPath := filepath.Join(dir, filePath)

		file, err := os.Open(fullPath)
		if err != nil {
			if os.IsNotExist(err) {
				log.Printf("No existing author database found in %s, starting fresh.\n", fullPath)
				s.NextID = 1
				return nil
			}
			log.Printf("Failed to open file %s: %v\n", fullPath, err)
			return err
		}
		defer file.Close()

		var data struct {
			Authors map[int]Author `json:"authors"`
			NextID  int            `json:"next_id"`
		}

		if err := json.NewDecoder(file).Decode(&data); err != nil {
			log.Printf("Failed to decode authors from file %s: %v\n", fullPath, err)
			return err
		}

		s.Mu.Lock()
		defer s.Mu.Unlock()
		s.Authors = data.Authors
		s.NextID = data.NextID
		log.Printf("Authors loaded successfully from %s\n", fullPath)
		return nil
	}
}

func (s *InMemoryAuthorStore) SaveAuthors(ctx context.Context, filePath string) error {
	select {
	case <-ctx.Done():
		log.Println("Context canceled during author saving")
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
			Authors map[int]Author `json:"authors"`
			NextID  int            `json:"next_id"`
		}{
			Authors: s.Authors,
			NextID:  s.NextID,
		}

		file, err := os.Create(fullPath)
		if err != nil {
			log.Printf("Failed to create file %s: %v\n", fullPath, err)
			return err
		}
		defer file.Close()

		if err := json.NewEncoder(file).Encode(data); err != nil {
			log.Printf("Failed to write authors to file %s: %v\n", fullPath, err)
			return err
		}

		log.Printf("Authors saved successfully to %s\n", fullPath)
		return nil
	}
}
