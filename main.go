package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func PrettyEncode(data interface{}, out io.Writer) error {
	enc := json.NewEncoder(out)
	enc.SetIndent("", "    ")
	if err := enc.Encode(data); err != nil {
		return err
	}
	return nil
}

type Stock struct {
	GlobalQuote struct {
		Zero1Symbol           string `json:"01. symbol"`
		Zero2Open             string `json:"02. open"`
		Zero3High             string `json:"03. high"`
		Zero4Low              string `json:"04. low"`
		Zero5Price            string `json:"05. price"`
		Zero6Volume           string `json:"06. volume"`
		Zero7LatestTradingDay string `json:"07. latest trading day"`
		Zero8PreviousClose    string `json:"08. previous close"`
		Zero9Change           string `json:"09. change"`
		One0ChangePercent     string `json:"10. change percent"`
	} `json:"Global Quote"`
}

func main() {
	fmt.Println("This project runs on the Alpha Vantage API")
	fmt.Println("If you do not already have an API key please visit:")
	fmt.Println("https://www.alphavantage.co/support/#api-key")
	fmt.Println()

	tick, key := Auth()

	getStockInfo(tick, key)

	for range time.Tick(time.Minute) {
		getStockInfo(tick, key)
	}
}

func Auth() (string, string) {
	fmt.Print("Enter ticker you wish to track: ")
	var tick string
	fmt.Scanln(&tick)

	fmt.Println("Enter API Key: ")
	var key string
	fmt.Scanln(&key)

	return tick, key
}

func getStockInfo(symbol string, key string) {
	resp, err := http.Get("https://www.alphavantage.co/query?function=GLOBAL_QUOTE&symbol=" + symbol + "&apikey=" + key)
	if err != nil {
		fmt.Println("No response from request")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body) // response body is []byte
	//fmt.Println(string(body))              // convert to string before print

	var stock Stock
	if err := json.Unmarshal(body, &stock); err != nil { // Parse []byte to go struct pointer
		fmt.Println("Can not unmarshal JSON")
	}

	var buffer bytes.Buffer
	err = PrettyEncode(stock, &buffer)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(buffer.String())
}
