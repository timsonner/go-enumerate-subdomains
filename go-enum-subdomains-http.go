package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"sync"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run <script> <domain>")
		return
	}

	domain := os.Args[1]
	wordlist := os.Args[2]

	file, err := os.Open(wordlist)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Create a WaitGroup to manage concurrency
	var wg sync.WaitGroup

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		subdomain := scanner.Text()
		subDomainURL := fmt.Sprintf("http://%s.%s", subdomain, domain)

		// Increment the WaitGroup counter
		wg.Add(1)

		// Launch a goroutine for each subdomain check
		go func(url string) {
			// Decrement the counter when the goroutine completes
			defer wg.Done()

			response, err := http.Get(url)
			if err != nil {
				// If there's an error, it's probably a connection error, so we skip it
				// fmt.Println("Error connecting to:", url, "-", err)
				return
			}
			defer response.Body.Close()

			// Check if we have a successful status code
			if response.StatusCode == http.StatusOK {
				fmt.Println("Valid domain:", url)
			} else {
				fmt.Println("Invalid domain or other status code:", url, "Status Code:", response.StatusCode)
			}
		}(subDomainURL) // Pass subDomainURL as an argument to the goroutine
	}

	// Wait for all goroutines to finish
	wg.Wait()

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}
}
