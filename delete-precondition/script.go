package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"testscripts/utils"
)

var (
	userCount   = flag.Int("userCount", 150000, "Number of users to create")
	withAddress = flag.Bool("withAddress", false, "Include address in user data")
	withDOB     = flag.Bool("withDOB", false, "Include date of birth in user data")
)

func deletePrecondition() {
	fmt.Printf("Creating %d users...\n", *userCount)

	successCount := 0
	failureCount := 0

	for i := 1; i <= *userCount; i++ {
		// Generate random user data
		user := utils.GenerateRandomUser(i, *withAddress, *withDOB)

		// Convert user to JSON
		userJSON, err := json.Marshal(user)
		if err != nil {
			fmt.Printf("Error marshaling user %d: %v\n", i, err)
			failureCount++
			continue
		}

		// Create HTTP POST request
		resp, err := http.Post("http://localhost:8080/user", "application/json", bytes.NewBuffer(userJSON))
		if err != nil {
			fmt.Printf("Error creating user %d: %v\n", i, err)
			failureCount++
			continue
		}

		// Check response status
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			successCount++
			if i%100 == 0 {
				fmt.Printf("Created %d users so far...\n", i)
			}
		} else {
			fmt.Printf("Failed to create user %d: HTTP %d\n", i, resp.StatusCode)
			failureCount++
		}

		resp.Body.Close()
	}

	fmt.Printf("\n=== SUMMARY ===\n")
	fmt.Printf("Total users attempted: %d\n", *userCount)
	fmt.Printf("Successfully created: %d\n", successCount)
	fmt.Printf("Failed: %d\n", failureCount)
	fmt.Printf("Success rate: %.2f%%\n", float64(successCount)/float64(*userCount)*100)
}

func main() {
	flag.Parse()
	deletePrecondition()
}
