package main

import (
	"crypto/rand"
	"math/big"
)

// https://gist.github.com/dopey/c69559607800d2f2f90b1b1ed4e550fb
// Generates a random, cryptographically secure string of length n
func randomString(n int) string {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			panic(err)
		}
		ret[i] = letters[num.Int64()]
	}

	return string(ret)
}
