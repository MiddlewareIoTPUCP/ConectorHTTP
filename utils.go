package main

import (
	"log"
	"math/rand"
)

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func randomString(n int) string {
	bytes := make([]byte, n)
	for i := 0; i < n; i++ {
		bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
