package durable

import (
	"math/rand"
	"time"
)

func GenerateAccountNumber() string {
	const digits = "0123456789"
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	accountNumber := make([]byte, 26)
	for i := range accountNumber {
		accountNumber[i] = digits[r.Intn(len(digits))]
	}

	return "TR" + string(accountNumber)
}

func GenerateAmount() float64 {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	return float64(r.Intn(10000) + 50)
}
