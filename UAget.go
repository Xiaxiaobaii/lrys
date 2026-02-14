package lrys

import (
	"math/rand"
	"time"
)

var UA []string = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36 Edg/122.0.0.0",
}

func GetUa() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	rand := r.Intn(len(UA))
	return UA[rand]
}
