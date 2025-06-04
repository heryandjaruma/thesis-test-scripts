package main

import (
	"flag"
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"time"

	vegeta "github.com/tsenart/vegeta/v12/lib"
)

var (
	withAddress = flag.Bool("withAddress", false, "Include address in user data")
	withDOB     = flag.Bool("withDOB", false, "Include date of birth in user data")
	maxRate     = flag.Int("maxRate", 1000, "Maximum rate of requests per second")
	runNo       = flag.Int("runNo", 1, "Run number")
)

func runGETAttack() {
	rand.Seed(time.Now().UnixNano())

	currentRate := 10
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

		currentRate += 10

		if currentRate > *maxRate {
			fmt.Printf("\n🛑 MAX RATE REACHED! Stopping at %d req/s.\n", *maxRate)
			break
		}
	}

	fmt.Printf("\n=== FINAL SUMMARY RUN #%d ===", *runNo)
	var totalRequests, totalSuccess uint64
	var maxLatency, totalLatency time.Duration
	maxRate := 0

	for i, m := range allMetrics {
		rate := 10 + (i * 10)
		if rate > maxRate {
			maxRate = rate
		}
		totalRequests += m.Requests
		totalSuccess += uint64(float64(m.Requests) * m.Success)
		totalLatency += m.Latencies.Mean
		if m.Latencies.Max > maxLatency {
			maxLatency = m.Latencies.Max
		}
	}

	overallSuccess := float64(totalSuccess) / float64(totalRequests) * 100

	sort.Slice(allLatencies, func(i, j int) bool {
		return allLatencies[i] < allLatencies[j]
	})

	minLatency := allLatencies[0]
	maxLatency = allLatencies[len(allLatencies)-1]
	p50 := allLatencies[int(float64(len(allLatencies))*0.50)]
	p90 := allLatencies[int(float64(len(allLatencies))*0.90)]
	p95 := allLatencies[int(float64(len(allLatencies))*0.95)]
	p99 := allLatencies[int(float64(len(allLatencies))*0.99)]

	var totalLatencySum time.Duration
	for _, lat := range allLatencies {
		totalLatencySum += lat
	}
	meanLatency := totalLatencySum / time.Duration(len(allLatencies))

	fmt.Printf("\nTest Duration: %d seconds\n", len(allMetrics))
	fmt.Printf("Maximum Rate Achieved: %d req/s\n", maxRate)
	fmt.Printf("Total Requests: %d\n", totalRequests)
	fmt.Printf("Total Users Queried: %d\n", userID-1)
	fmt.Printf("Overall Success Rate: %.2f%%\n", overallSuccess)
	fmt.Printf("\n=== LATENCY PERCENTILES ===\n")
	fmt.Printf("Min:  %s\n", minLatency)
	fmt.Printf("Mean: %s\n", meanLatency)
	fmt.Printf("50th: %s\n", p50)
	fmt.Printf("90th: %s\n", p90)
	fmt.Printf("95th: %s\n", p95)
	fmt.Printf("99th: %s\n", p99)
	fmt.Printf("Max:  %s\n", maxLatency)
}

func main() {
	flag.Parse()

	runGETAttack()
}
