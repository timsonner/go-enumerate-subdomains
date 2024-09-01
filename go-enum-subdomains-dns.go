package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"sync"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run <script> <domain> <wordlist>")
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
		subDomainURL := fmt.Sprintf("%s.%s", subdomain, domain)

		// Increment the WaitGroup counter
		wg.Add(1)

		// Launch a goroutine for each subdomain check
		go func(url string) {
			// Decrement the counter when the goroutine completes
			defer wg.Done()

			// Perform DNS lookup
			ips, err := net.LookupIP(url)
			if err != nil {
				// fmt.Printf("No record found for %s: %v\n", url, err)
				return
			}

			// Print the resolved IP addresses
			fmt.Printf("%s resolves to %v\n", url, ips)
		}(subDomainURL) // Pass subDomainURL as an argument to the goroutine
	}

	// Wait for all goroutines to finish
	wg.Wait()

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}
}
