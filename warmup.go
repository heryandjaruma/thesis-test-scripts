package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"testscripts/utils"
	"time"

	vegeta "github.com/tsenart/vegeta/v12/lib"
)

func runPOSTAttack() {
	rand.Seed(time.Now().UnixNano())

	rate := vegeta.Rate{Freq: 10, Per: time.Second}
	duration := 1 * time.Second
	userID := 1

	targeter := func(tgt *vegeta.Target) error {
		user := utils.GenerateRandomUser(userID, true, true)
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
		user := utils.GenerateRandomUser(userID, true, false)

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
