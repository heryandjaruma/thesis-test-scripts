package utils

import (
	"math/rand"
	"strconv"
	"time"
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
