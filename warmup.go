package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
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

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func randomJavaCompatibleDateString() string {
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

func generateRandomUser(id int, withAddress bool, withDOB bool) User {
	username := randomString(10)
	name := randomString(15)
	email := randomString(20)
	phone := randomString(10)

	if withAddress {
		address := randomString(30)
		return User{
			ID:       strconv.Itoa(id),
			Username: username,
			Name:     name,
			Phone:    phone,
			Email:    email,
			Address:  &address,
		}
	} else if withDOB {
		address := randomString(30)
		dateOfBirth := randomJavaCompatibleDateString()
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

func runPOSTAttack() {
	rand.Seed(time.Now().UnixNano())

	rate := vegeta.Rate{Freq: 10, Per: time.Second}
	duration := 1 * time.Second
	userID := 1

	targeter := func(tgt *vegeta.Target) error {
		user := generateRandomUser(userID, true, true)
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

	fmt.Printf("Starting POST attack: 10 requests\n")

	for res := range attacker.Attack(targeter, rate, duration, "POST User Creation!") {
		metrics.Add(res)
	}
	metrics.Close()

	fmt.Printf("POST attack complete\n")
}

func runGETAttack() {
	rate := vegeta.Rate{Freq: 10, Per: time.Second}
	duration := 1 * time.Second
	userID := 1

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

	fmt.Printf("Starting GET attack: 10 requests\n")

	for res := range attacker.Attack(targeter, rate, duration, "GET User Retrieval!") {
		metrics.Add(res)
	}
	metrics.Close()

	fmt.Printf("GET attack complete\n")
}

func runPUTAttack() {
	rand.Seed(time.Now().UnixNano())
	rate := vegeta.Rate{Freq: 10, Per: time.Second}
	duration := 1 * time.Second
	userID := 1

	targeter := func(tgt *vegeta.Target) error {
		user := generateRandomUser(userID, true, false)

		payload, err := json.Marshal(user)
		if err != nil {
			return err
		}

		tgt.Method = "PUT"
		tgt.URL = "http://localhost:8080/user"
		tgt.Body = payload
		tgt.Header = map[string][]string{
			"Content-Type": {"application/json"},
		}
		userID++
		return nil
	}

	attacker := vegeta.NewAttacker()
	var metrics vegeta.Metrics

	fmt.Printf("Starting PUT attack: 10 requests\n")

	for res := range attacker.Attack(targeter, rate, duration, "PUT User Update!") {
		metrics.Add(res)
	}
	metrics.Close()

	fmt.Printf("PUT attack complete\n")
}

func runDELETEAttack() {
	rate := vegeta.Rate{Freq: 10, Per: time.Second}
	duration := 1 * time.Second
	userID := 1

	targeter := func(tgt *vegeta.Target) error {
		tgt.Method = "DELETE"
		tgt.URL = fmt.Sprintf("http://localhost:8080/user/%d", userID)
		tgt.Header = map[string][]string{
			"Content-Type": {"application/json"},
		}
		userID++
		return nil
	}

	attacker := vegeta.NewAttacker()
	var metrics vegeta.Metrics

	fmt.Printf("Starting DELETE attack: 10 requests\n")

	for res := range attacker.Attack(targeter, rate, duration, "DELETE User Removal!") {
		metrics.Add(res)
	}
	metrics.Close()

	fmt.Printf("DELETE attack complete\n")
}

func main() {
	fmt.Print("Starting warmup\n")
	runPOSTAttack()
	time.Sleep(1 * time.Second)
	runGETAttack()
	time.Sleep(1 * time.Second)
	runPUTAttack()
	time.Sleep(1 * time.Second)
	runDELETEAttack()
	fmt.Print("Warmup complete\n")
}
