package restacular

import (
	"fmt"
	"log"
	"math/rand"
	"os"
)

const alphanum = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

var Random *os.File

func init() {
	f, err := os.Open("/dev/urandom")
	if err != nil {
		log.Fatal(err)
	}
	Random = f
}

// Generates a short random id, to use instead of incrementing integer primary keys
func NewUid(size int) string {
	buffer := make([]byte, size)
	for i := 0; i < size; i++ {
		buffer[i] = alphanum[rand.Intn(len(alphanum))]
	}
	return string(buffer)
}

// Generates a valid UUID4
func NewUuid() string {
	b := make([]byte, 16)
	Random.Read(b)
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}
