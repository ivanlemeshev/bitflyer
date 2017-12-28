package bitflyer

import (
	"log"
	"os"
	"testing"
)

func TestGetBalance(t *testing.T) {
	key := os.Getenv("BITFLYER_KEY")
	secret := os.Getenv("BITFLYER_SECRET")

	api := New(key, secret)
	ret, body, err := api.GetBalance()
	if err != nil {
		panic(err)
	}

	log.Printf("err:%v", err)
	log.Printf("ret:%v", ret)
	log.Printf("body:%s", string(body))

	return
}
