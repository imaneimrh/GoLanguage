package main

import (
	"context"
	"log"
	"sync"
	"time"
)

var wg sync.WaitGroup

func fetchURL(ctx context.Context, url string) {
	defer wg.Done()
	select {
	case <-ctx.Done():
		log.Println("Graceful shutdown", ctx.Err())
	default:
		log.Println("Fetched data, from this url: "+url, nil)
	}
}

func fetchURLtoTest(ctx context.Context, url string) {
	defer wg.Done()
	select {
	case <-ctx.Done():
		log.Println("Graceful shutdown", ctx.Err())
	case <-time.After(6 * time.Second):
		log.Println("Fetched data, from this url: "+url, nil)
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	ctx2, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()
	urls := []string{}
	url1 := "https://imane.com"
	url2 := "https://golang.com"
	url3 := "https://google.com"
	urls = append(urls, url1, url2, url3)

	log.Println("\nStarting with fetching the URLS...")
	for _, url := range urls {
		select {
		case <-ctx.Done():
			log.Println("The canceling of the request has been done")
		default:
			wg.Add(1)
			fetchURL(ctx, url)
		}
	}

	log.Println("\nStarting with the testing of the fetching of the URLS, after the context Timeout...")
	for _, url := range urls {
		select {
		case <-ctx2.Done():
			log.Println("The canceling of the request has been done via test functions")
		default:
			wg.Add(1)
			fetchURLtoTest(ctx2, url)
		}
	}

	wg.Wait()
}
