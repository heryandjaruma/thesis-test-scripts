package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"testscripts/utils"
	"time"

	vegeta "github.com/tsenart/vegeta/v12/lib"
)

var (
	userCount   = flag.Int("userCount", 150000, "Number of users to create")
	withAddress = flag.Bool("withAddress", false, "Include address in user data")
	withDOB     = flag.Bool("withDOB", false, "Include date of birth in user data")
)

func updatePrecondition() {
	fmt.Printf("Creating %d users with rate of 1000 requests per second...\n", *userCount)

	// Create a new targeter
	targeter := func(t *vegeta.Target) error {
		// Generate random user data
		user := utils.GenerateRandomUser(int(time.Now().UnixNano()), *withAddress, *withDOB)

		// Convert user to JSON
		userJSON, err := json.Marshal(user)
		if err != nil {
			return err
		}

		t.Method = "POST"
		t.URL = "http://localhost:8080/user"
		t.Body = userJSON
		t.Header = http.Header{"Content-Type": {"application/json"}}
		return nil
	}

	// Create an attacker
	rate := vegeta.Rate{Freq: 1000, Per: time.Second}
	duration := time.Duration(float64(*userCount) / 1000 * float64(time.Second))
	attacker := vegeta.NewAttacker()

	// Run the attack
	var metrics vegeta.Metrics
	for res := range attacker.Attack(targeter, rate, duration, "Load Test") {
		metrics.Add(res)
	}
	metrics.Close()

	fmt.Printf("\n=== SUMMARY ===\n")
	fmt.Printf("Total requests: %d\n", metrics.Requests)
}

func main() {
	flag.Parse()
	updatePrecondition()
}
