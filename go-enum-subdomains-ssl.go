package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

// CertEntry represents the JSON structure of the crt.sh API response
type CertEntry struct {
	CommonName string `json:"common_name"`
	NameValue  string `json:"name_value"`
}

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Usage: go run <script> <domain>")
		return
	}

	domain := os.Args[1]

	// Query crt.sh for all certificates associated with the domain
	url := fmt.Sprintf("https://crt.sh/?q=%%25.%s&output=json", domain)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to query crt.sh: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// Decode the JSON response
	var certs []CertEntry
	if err := json.NewDecoder(resp.Body).Decode(&certs); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to decode JSON response: %v\n", err)
		return
	}

	// Use a map to store unique subdomains
	subdomainSet := make(map[string]struct{})

	// Extract unique common names and name values (subdomains)
	for _, cert := range certs {
		// Add common name to the set
		if cert.CommonName != "" && strings.Contains(cert.CommonName, domain) {
			subdomainSet[cert.CommonName] = struct{}{}
		}

		// Add each name value to the set
		nameValues := strings.Split(cert.NameValue, "\n")
		for _, name := range nameValues {
			name = strings.TrimSpace(name)
			if name != "" && strings.Contains(name, domain) {
				subdomainSet[name] = struct{}{}
			}
		}
	}

	// Display unique subdomains
	fmt.Println("Discovered subdomains:")
	for subdomain := range subdomainSet {
		fmt.Println("- " + subdomain)
	}
}
