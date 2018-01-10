package bitflyer

import (
	"log"
	"os"
	"testing"
)

func TestGetBalanceFail(t *testing.T) {

	api := New("wrong_key", "wrong_password")
	ret, body, err := api.GetBalance()
	log.Printf("err:%v", err)
	if err == nil {
		panic("should be error")
	}

	log.Printf("ret:%v", ret)
	log.Printf("body:%s", string(body))

	return
}

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
