package main

import (
	"log"
	"strings"
)

func countfreq(sentence []string) map[string]int {
	freq := make(map[string]int)
	for _, word := range sentence {
		word = strings.Trim(word, "@!,;:.-")
		word = strings.ToLower(word)
		freq[word]++
	}
	return freq
}

func main() {

	sentence := "Hello hello, dear friends And Their FRIENDS!"

	words := strings.Split(sentence, " ")

	freq := countfreq(words)
	for word, count := range freq {
		log.Printf("%s -> %d\n", word, count)
	}

}
