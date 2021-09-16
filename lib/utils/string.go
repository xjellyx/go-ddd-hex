package utils

import (
	"math/rand"
	"time"
)

var (
	str = "qwertyuiopasdfghjklzxcvbnm1234567890QAZWSXEDCRFVTGBYHNUJMIKOLP"
)

func RandString(length int) string {
	rand.Seed(time.Now().UnixNano())
	var (
		res string
	)
	for length > 0 {
		i := rand.Intn(len(str))
		res += string(str[i])
		length--
	}
	return res
}
