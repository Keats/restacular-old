package restacular

import (
	"math/rand"
)

const alphanum = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

// Generates a short random id, to use instead of incrementing integer primary keys
func NewUid(size int) string {
	buffer := make([]byte, size)
	for i := 0; i < size; i++ {
		buffer[i] = alphanum[rand.Intn(len(alphanum))]
	}
	return string(buffer)
}
