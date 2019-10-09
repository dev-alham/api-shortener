package utils

import (
	"time"
	"math/rand"
	"strings"
)

func RandStr(length int) string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("ABCDEFGHJKLMNPQRSTWXYZ" +
		"123456789")

	var temp_str strings.Builder

	for a:=0; a < length; a++ {
		temp_str.WriteRune(chars[rand.Intn(len(chars))])
	}
	str := temp_str.String()
	return str
}

