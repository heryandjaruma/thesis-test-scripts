package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"strconv"
	"testscripts/utils"
	"time"

	vegeta "github.com/tsenart/vegeta/v12/lib"
)

var (
	withAddress = flag.Bool("withAddress", false, "Include address in user data")
	withDOB     = flag.Bool("withDOB", false, "Include date of birth in user data")
	maxRate     = flag.Int("maxRate", 1000, "Maximum rate of requests per second")
)

func runPOSTAttack() {
	rand.Seed(time.Now().UnixNano())

	currentRate := 10
	duration := 1 * time.Second
	failureThreshold := 0.95    // Stop if success rate drops below 95%
	maxConsecutiveFailures := 3 // Stop after 3 consecutive periods with low success
	consecutiveFailures := 0
	userID := 1 // Starting ID for users

	var allMetrics []vegeta.Metrics

	for {
		fmt.Printf("Running attack at %d requests/second...\n", currentRate)

		rate := vegeta.Rate{Freq: currentRate, Per: time.Second}

		// Create a custom targeter that generates random data for each request
		targeter := func(tgt *vegeta.Target) error {
			user := utils.GenerateRandomUser(userID, *withAddress, *withDOB)
			userID++

			payload, err := json.Marshal(user)
			if err != nil {
				return err
			}

			tgt.Method = "POST"
			tgt.URL = "http://localhost:8080/user"
			tgt.Body = payload
			tgt.Header = map[string][]string{
				"Content-Type": {"application/json"},
			}
			return nil
		}

		attacker := vegeta.NewAttacker()

		var metrics vegeta.Metrics

		for res := range attacker.Attack(targeter, rate, duration, "POST User Attack!") {
			metrics.Add(res)
		}
		metrics.Close()
		allMetrics = append(allMetrics, metrics)

		fmt.Printf("Rate %d req/s - Success: %.2f%%, Mean latency: %s, Max latency: %s\n",
			currentRate, metrics.Success*100, metrics.Latencies.Mean, metrics.Latencies.Max)

		if metrics.Success < failureThreshold {
			consecutiveFailures++
			fmt.Printf("⚠️  Success rate below threshold! Consecutive failures: %d/%d\n",
				consecutiveFailures, maxConsecutiveFailures)

			if consecutiveFailures >= maxConsecutiveFailures {
				fmt.Printf("\n🔴 API FAILURE DETECTED! Stopping test after %d consecutive low-success periods.\n",
					maxConsecutiveFailures)
				break
			}
		}

		if metrics.Requests == 0 {
			fmt.Println("\n🔴 NO REQUESTS COMPLETED! API appears to be down.")
			break
		}

		// Check status codes for server errors (fixed type conversion)
		serverErrors := 0
		for codeStr, count := range metrics.StatusCodes {
			code, err := strconv.Atoi(codeStr)
			if err == nil && code >= 500 {
				serverErrors += int(count)
			}
		}
		if serverErrors > 0 {
			errorRate := float64(serverErrors) / float64(metrics.Requests)
			fmt.Printf("⚠️  Server errors detected: %d (%.1f%%)\n", serverErrors, errorRate*100)
			if errorRate > 0.1 { // More than 10% server errors
				fmt.Println("\n🔴 HIGH SERVER ERROR RATE! Stopping test.")
				break
			}
		}

		fmt.Println("---")

		currentRate += 10

		if currentRate > *maxRate {
			fmt.Println("\n🛑 SAFETY LIMIT REACHED! Stopping at 1000 req/s to prevent system overload.")
			break
		}
	}

	// Print overall summary
	fmt.Println("\n=== FINAL SUMMARY ===")
	var totalRequests, totalSuccess uint64
	var maxLatency, totalLatency time.Duration
	maxRate := 0

	for i, m := range allMetrics {
		rate := (i + 1) * 5
		if rate > maxRate {
			maxRate = rate
		}
		totalRequests += m.Requests
		totalSuccess += uint64(float64(m.Requests) * m.Success)
		totalLatency += m.Latencies.Mean
		if m.Latencies.Max > maxLatency {
			maxLatency = m.Latencies.Max
		}
		fmt.Printf("Second %d (Rate %d): Success %.2f%%, Mean %s\n",
			i+1, rate, m.Success*100, m.Latencies.Mean)
	}

	overallSuccess := float64(totalSuccess) / float64(totalRequests) * 100
	averageLatency := totalLatency / time.Duration(len(allMetrics))

	fmt.Printf("\nTest Duration: %d seconds\n", len(allMetrics))
	fmt.Printf("Maximum Rate Achieved: %d req/s\n", maxRate)
	fmt.Printf("Total Requests: %d\n", totalRequests)
	fmt.Printf("Total Users Created: %d\n", userID-1)
	fmt.Printf("Overall Success Rate: %.2f%%\n", overallSuccess)
	fmt.Printf("Average Latency: %s\n", averageLatency)
	fmt.Printf("Maximum Latency: %s\n", maxLatency)
}

func main() {
	flag.Parse()

	runPOSTAttack()
}
