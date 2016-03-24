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
	Price float64 `json:"price"`
	Size  float64 `json:"size"`
}

// OrderBook represents bitFlyer Lightning order book.
type OrderBook struct {
	MidPrice int      `json:"mid_price"`
	Bids     []AskBid `json:"bids"`
	Asks     []AskBid `json:"asks"`
}

// AssetBalance represents bitFlyer Lightning asset balance.
type AssetBalance struct {
	Balances []Balance `json:"balances"`
}

// Balance represents bitFlyer Lightning asset balance record.
type Balance struct {
	CurrencyCode string  `json:"currency_code"`
	Amount       float64 `json:"amount"`
	Available    float64 `json:"available"`
}

// Ticker represents bitFlyer Lightning ticker.
type Ticker struct {
	ProductCode     string  `json:"product_code"`
	Timestamp       string  `json:"timestamp"`
	TickID          int     `json:"tick_id"`
	BestBid         float64 `json:"best_bid"`
	BestAsk         float64 `json:"best_ask"`
	BestBidSize     float64 `json:"best_bid_size"`
	BestAskSize     float64 `json:"best_ask_size"`
	TotalBidDepth   float64 `json:"total_bid_depth"`
	TotalAskDepth   float64 `json:"total_ask_depth"`
	LTP             float64 `json:"ltp"`
	Volume          float64 `json:"volume"`
	VolumeByProduct float64 `json:"volume_by_product"`
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

// GetTicker returns bitFlyer Lightning ticker.
func (api APIClient) GetTicker() (ticker Ticker, err error) {
	err = api.doGetRequest(URL+"/v1/getticker", &ticker)
	if err != nil {
		return ticker, err
	}
	return ticker, nil
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
