package stores

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	. "FinalProject/models"
)

// ----------------------------------------------Definition of BookMethods--------------------------------
type InMemoryBookStore struct {
	Mu     sync.RWMutex
	Books  map[int]Book
	NextID int
}

type BookStore interface {
	CreateBook(ctx context.Context, book Book) (Book, error)
	GetBook(ctx context.Context, id int) (Book, error)
	UpdateBook(ctx context.Context, id int, book Book) (Book, error)
	DeleteBook(ctx context.Context, id int) error
	SearchBooks(ctx context.Context, criteria SearchCriteria) ([]Book, error)
	LoadBooks(ctx context.Context, filePath string) error
	SaveBooks(ctx context.Context, filePath string) error
}

func (s *InMemoryBookStore) CreateBook(ctx context.Context, book Book, auths *InMemoryAuthorStore) (Book, error) {
	select {
	case <-ctx.Done():
		log.Println("Request canceled during book creation")
		return Book{}, ctx.Err()
	default:
		s.Mu.Lock()
		defer s.Mu.Unlock()

		author, ok := auths.Authors[book.Author.ID]
		if !ok {
			return Book{}, errors.New("Author with ID " + strconv.Itoa(book.Author.ID) + " not found")
		}
		book.Author = author

		book.ID = s.NextID
		s.NextID++
		s.Books[book.ID] = book

		log.Printf("Book created successfully. ID: %d\n", book.ID)
		return book, nil
	}
}

func (s *InMemoryBookStore) GetBook(ctx context.Context, bookId int, auths *InMemoryAuthorStore) (Book, error) {
	select {
	case <-ctx.Done():
		log.Println("Request canceled during book retrieval of book")
		return Book{}, ctx.Err()
	default:
		s.Mu.Lock()
		defer s.Mu.Unlock()
		book, ok := s.Books[bookId]
		if !ok {
			log.Println("Book with ID ", bookId, " not found")
			return Book{}, errors.New("Book with ID " + strconv.Itoa(bookId) + " not found")
		}
		author, ok := auths.Authors[book.Author.ID]
		if ok {
			book.Author = author
		} else {
			log.Printf("Author with ID %d not found for book ID %d\n", book.Author.ID, book.ID)
		}
		return book, nil
	}
}

func (s *InMemoryBookStore) UpdateBook(ctx context.Context, bookId int, book Book, auths *InMemoryAuthorStore) (Book, error) {
	select {
	case <-ctx.Done():
		log.Println("Request canceled during book update")
		return Book{}, ctx.Err()
	default:
		s.Mu.Lock()
		defer s.Mu.Unlock()
		authors, _ := auths.ListAuthors(ctx)
		foundAuthor := false
		for _, a := range authors {
			if a.FirstName == book.Author.FirstName && a.LastName == book.Author.LastName {
				unchangedBook := s.Books[bookId]
				book.ID = unchangedBook.ID
				s.Books[bookId] = book
				foundAuthor = true
				return s.Books[bookId], nil
			}
		}
		if !foundAuthor {
			log.Println("Author with name", book.Author.FirstName, "and last name", book.Author.LastName, "not found")
			_, err := auths.CreateAuthor(ctx, book.Author)
			if err != nil {
				log.Println("Error creating the author for the book you're trying to update")
				return s.Books[bookId], err
			}
			log.Println("Author with name", book.Author.FirstName, "and last name", book.Author.LastName, "was created, in order to update book")
			unchangedBook := s.Books[bookId]
			book.ID = unchangedBook.ID
			s.Books[bookId] = book
			return s.Books[bookId], nil
		}

		return Book{}, errors.New("Book with id " + strconv.Itoa(bookId) + "not found")
	}
}

func (s *InMemoryBookStore) DeleteBook(ctx context.Context, bookId int) error {
	select {
	case <-ctx.Done():
		log.Println("Request canceled during book deletion")
		return ctx.Err()
	default:
		s.Mu.Lock()
		defer s.Mu.Unlock()
		if _, ok := s.Books[bookId]; ok {
			delete(s.Books, bookId)
			return nil
		}
		return errors.New("Book with id " + strconv.Itoa(bookId) + "not found")
	}
}

func (s *InMemoryBookStore) SearchBooks(ctx context.Context, criteria SearchCriteria) ([]Book, error) {
	select {
	case <-ctx.Done():
		log.Println("Request canceled during book search")
		return nil, ctx.Err()
	default:
		s.Mu.Lock()
		defer s.Mu.Unlock()
		var result []Book
		for _, book := range s.Books {
			if strings.Contains(book.Title, criteria.Title) ||
				strings.Contains(book.Author.FirstName, criteria.Author) ||
				strings.Contains(book.Author.LastName, criteria.Author) ||
				strings.Contains(strings.Join(book.Genres, ","), criteria.Genre) {
				result = append(result, book)
			}
		}
		if len(result) == 0 {
			return nil, errors.New("No books found")
		}
		return result, nil
	}
}

func (s *InMemoryBookStore) LoadBooks(ctx context.Context, filePath string) error {
	select {
	case <-ctx.Done():
		log.Println("Context canceled during book loading")
		return ctx.Err()
	default:
		dir := "database"
		fullPath := filepath.Join(dir, filePath)

		file, err := os.Open(fullPath)
		if err != nil {
			if os.IsNotExist(err) {
				log.Printf("No existing book database found in %s, starting fresh.\n", fullPath)
				s.NextID = 1
				return nil
			}
			log.Printf("Failed to open file %s: %v\n", fullPath, err)
			return err
		}
		defer file.Close()

		var data struct {
			Books  map[int]Book `json:"books"`
			NextID int          `json:"next_id"`
		}

		if err := json.NewDecoder(file).Decode(&data); err != nil {
			log.Printf("Failed to decode books from file %s: %v\n", fullPath, err)
			return err
		}

		s.Mu.Lock()
		defer s.Mu.Unlock()
		s.Books = data.Books
		s.NextID = data.NextID
		log.Printf("Books loaded successfully from %s\n", fullPath)
		return nil
	}
}

func (s *InMemoryBookStore) SaveBooks(ctx context.Context, filePath string) error {
	select {
	case <-ctx.Done():
		log.Println("Context canceled during book saving")
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
			Books  map[int]Book `json:"books"`
			NextID int          `json:"next_id"`
		}{
			Books:  s.Books,
			NextID: s.NextID,
		}

		file, err := os.Create(fullPath)
		if err != nil {
			log.Printf("Failed to create file %s: %v\n", fullPath, err)
			return err
		}
		defer file.Close()

		if err := json.NewEncoder(file).Encode(data); err != nil {
			log.Printf("Failed to write books to file %s: %v\n", fullPath, err)
			return err
		}

		log.Printf("Books saved successfully to %s\n", fullPath)
		return nil
	}
}
