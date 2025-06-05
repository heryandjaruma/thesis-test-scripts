package main

import (
	"flag"
	"fmt"
	"math/rand"
	"strconv"
	"testscripts/utils"
	"time"

	vegeta "github.com/tsenart/vegeta/v12/lib"
)

var (
	maxRate      = flag.Int("maxRate", 5000, "Maximum rate of requests per second")
	testCaseName = flag.String("testCaseName", "", "Test case name")
	runNo        = flag.Int("runNo", 1, "Run number")
)

func runGETAttack() {
	rand.Seed(time.Now().UnixNano())

	multiplier := 100
	currentRate := 100
	duration := 1 * time.Second
	failureThreshold := 0.95
	maxConsecutiveFailures := 3
	consecutiveFailures := 0
	userID := 1

	var allMetrics []vegeta.Metrics
	var allLatencies []time.Duration

	for {
		fmt.Printf("Running attack at %d requests/second...\n", currentRate)

		rate := vegeta.Rate{Freq: currentRate, Per: time.Second}

		targeter := func(tgt *vegeta.Target) error {
			tgt.Method = "GET"
			tgt.URL = fmt.Sprintf("http://localhost:8080/user/%d", userID)
			tgt.Header = map[string][]string{
				"Content-Type": {"application/json"},
			}
			userID++
			return nil
		}

		attacker := vegeta.NewAttacker()

		var metrics vegeta.Metrics

		for res := range attacker.Attack(targeter, rate, duration, "POST User Attack!") {
			metrics.Add(res)
			allLatencies = append(allLatencies, res.Latency)
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
			if errorRate > 0.1 {
				fmt.Println("\n🔴 HIGH SERVER ERROR RATE! Stopping test.")
				break
			}
		}

		fmt.Println("---")

		currentRate += multiplier

		if currentRate > *maxRate {
			fmt.Printf("\n🛑 MAX RATE REACHED! Stopping at %d req/s.\n", *maxRate)
			break
		}
	}

	utils.OutputSummary(allMetrics, allLatencies, *testCaseName, *runNo)
}

func main() {
	flag.Parse()

	runGETAttack()
}
