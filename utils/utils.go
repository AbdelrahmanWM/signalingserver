package utils

import (
	"crypto/rand"
	"log"
	"math/big"
)

func GenerateRandomID(length int) string {

	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		randInt, err := rand.Int(rand.Reader,big.NewInt(int64(len(chars))))
		if err!=nil{
			log.Fatal(err)
		}
		result[i]=chars[randInt.Int64()]
	}
	return string(result)
}