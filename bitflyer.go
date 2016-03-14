package bitflyer

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

// URL is a bitFlyer Lightning API base URL.
const URL = "https://api.bitflyer.jp"

// APIClient struct represents bitFlyer Lightning API client.
type APIClient struct {
	key    string
	secret string
	client *http.Client
}

// AskBid represents bitFlyer Lightning order book ask or bid record.
type AskBid struct {
	Price float64
	Size  float64
}

// OrderBook represents bitFlyer Lightning order book.
type OrderBook struct {
	MidPrice int
	Bids     []AskBid
	Asks     []AskBid
}

// AssetBalance represents bitFlyer Lightning asset balance.
type AssetBalance struct {
	Balances []Balance
}

// Balance represents bitFlyer Lightning asset balance record.
type Balance struct {
	CurrencyCode string
	Amount       float64
	Available    float64
}

// New creates a new bitFlyer Lightning API client.
func New(key, secret string) (client *APIClient) {
	client = new(APIClient)
	client.key = key
	client.secret = secret
	client.client = new(http.Client)
	return client
}

// GetOrderBook returns bitFlyer Lightning order book.
func (api APIClient) GetOrderBook() (orderBook OrderBook, err error) {
	err = api.doGetRequest(URL+"/v1/getboard", &orderBook)
	if err != nil {
		return orderBook, err
	}
	return orderBook, nil
}

// GetBalance returns bitFlyer Lightning account asset balance.
func (api APIClient) GetBalance() (balance Balance, err error) {
	err = api.doGetRequest(URL+"/v1/me/getbalance", &balance)
	if err != nil {
		return balance, err
	}
	return balance, nil
}

func (api *APIClient) doGetRequest(endpoint string, data interface{}) (err error) {
	headers := headers(api.key, api.secret, "GET", endpoint, "")
	resp, err := api.doRequest("GET", endpoint, headers)
	if err != nil {
		return err
	}
	err = json.Unmarshal(resp, data)
	if err != nil {
		return err
	}
	return nil
}

func (api *APIClient) doRequest(method, endpoint string, headers map[string]string) ([]byte, error) {
	req, err := http.NewRequest(method, endpoint, nil)
	if err != nil {
		return nil, requestError(err.Error())
	}
	setHeaders(req, headers)
	resp, err := api.client.Do(req)
	if err != nil {
		return nil, requestError(err.Error())
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, requestError(err.Error())
	}
	return body, nil
}

func headers(key, secret, method, uri, body string) map[string]string {
	timestamp := strconv.Itoa(int(time.Now().Unix()))
	message := timestamp + method + uri + body
	signature := computeHmac256(message, secret)
	headers := map[string]string{
		"Content-Type":     "application/json",
		"ACCESS-KEY":       key,
		"ACCESS-TIMESTAMP": timestamp,
		"ACCESS-SIGN":      signature,
	}
	return headers
}

func computeHmac256(message string, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	return hex.EncodeToString(h.Sum(nil))
}

func requestError(err interface{}) error {
	return fmt.Errorf("Could not execute request! (%s)", err)
}

func setHeaders(req *http.Request, headers map[string]string) {
	for key, value := range headers {
		req.Header.Add(key, value)
	}
}
