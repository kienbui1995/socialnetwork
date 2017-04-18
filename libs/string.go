package libs

import (
	"math/rand"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

//RandStringBytes fun return a random string with n characters
func RandStringBytes(n int) string {
	b := make([]byte, n)
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)
	for i := range b {
		b[i] = letterBytes[r.Intn(len(letterBytes))]
	}
	return string(b)
}
