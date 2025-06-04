package utils

import (
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"time"

	vegeta "github.com/tsenart/vegeta/v12/lib"
)

type User struct {
	ID          string  `json:"id"`
	Username    string  `json:"username"`
	Name        string  `json:"name"`
	Phone       string  `json:"phone"`
	Email       string  `json:"email"`
	Address     *string `json:"address"`
	DateOfBirth *string `json:"dateOfBirth"`
}

// Generate a random string of length n
func RandomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// Generate a random Java-compatible date string
func RandomJavaCompatibleDateString() string {
	minYear := 2003
	maxYear := 2025
	year := rand.Intn(maxYear-minYear+1) + minYear

	month := rand.Intn(12) + 1

	day := rand.Intn(28) + 1

	hour := rand.Intn(24)

	minute := rand.Intn(60)

	second := rand.Intn(60)

	millisecond := rand.Intn(1000)

	randomTime := time.Date(year, time.Month(month), day, hour, minute, second, millisecond*1000000, time.UTC)

	return randomTime.Format("2006-01-02T15:04:05.000Z07:00")
}

func GenerateRandomUser(id int, withAddress bool, withDOB bool) User {
	username := RandomString(10)
	name := RandomString(15)
	email := RandomString(20)
	phone := RandomString(10)

	if withAddress {
		address := RandomString(30)
		return User{
			ID:       strconv.Itoa(id),
			Username: username,
			Name:     name,
			Phone:    phone,
			Email:    email,
			Address:  &address,
		}
	} else if withDOB {
		address := RandomString(30)
		dateOfBirth := RandomJavaCompatibleDateString()
		return User{
			ID:          strconv.Itoa(id),
			Username:    username,
			Name:        name,
			Phone:       phone,
			Email:       email,
			Address:     &address,
			DateOfBirth: &dateOfBirth,
		}
	}

	return User{
		ID:       strconv.Itoa(id),
		Username: username,
		Name:     name,
		Phone:    phone,
		Email:    email,
	}
}

func OutputSummary(allMetrics []vegeta.Metrics, allLatencies []time.Duration, runNo int) {
	fmt.Printf("\n=== FINAL SUMMARY RUN #%d ===", runNo)
	var totalRequests, totalSuccess uint64
	var maxLatency, totalLatency time.Duration
	maxRate := 0

	for i, m := range allMetrics {
		rate := (i + 1) * 10
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
