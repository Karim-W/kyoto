package utils

import (
	"math/rand"
	"strings"
)

const alphaPool = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

// geneates a random string given a particular length
// Param: length: uint
// Returns: string
func GenerateRandomString(length uint) string {
	res := strings.Builder{}
	for i := 0; i < int(length); i++ {
		idx := rand.Intn(62) // length of alphaPool
		res.WriteByte(alphaPool[idx])
	}
	return res.String()
}
