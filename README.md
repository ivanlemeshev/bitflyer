# Client for [bitFlyer Lightning API](https://lightning.bitflyer.jp/docs?lang=en)

Example of usage:

```go
package main

import (
	"fmt"

	"github.com/ivanlemeshev/bitflyer"
)

const (
	key    = "KEY"
	secret = "SECRET"
)

func main() {
	api := bitflyer.New(key, secret)
	orderBook, respBody, err := api.GetOrderBook()
	if err != nil {
		fmt.Println("Error: ", err.Error())
		return
	}

	fmt.Printf("Order Book: %+v\n", orderBook)
	fmt.Printf("respBody: %s\n", string(respBody))
}
```

Todo:
- [X] Order Book
- [X] Ticker
- [X] Get Account Asset Balance
- [x] Send a New Order
- [x] List Orders
- [ ] Full public API
- [ ] Full private API
