package yahooFinanceAPI

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/zabenno/ticker_overview_data/tickerPrice"
)

type APIResponse struct {
	Response QuoteSummary `json:"quoteSummary"`
}

type APIError struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

type QuoteSummary struct {
	Result []ResultItem `json:"result"`
	Error  APIError     `json:"error"`
}

type ResultItem struct {
	Price Price `json:"price"`
}

type Price struct {
	RegMarketChangePercent  market `json:"regularMarketChangePercent"`
	RegMarketChange         market `json:"regularMarketChange"`
	RegMarketChangePrice    market `json:"regularMarketPrice"`
	PreMarketChangePercent  market `json:"preMarketChangePercent"`
	PreMarketChange         market `json:"preMarketChange"`
	PreMarketChangePrice    market `json:"preMarketPrice"`
	PostMarketChangePercent market `json:"postMarketChangePercent"`
	PostMarketChange        market `json:"postMarketChange"`
	PostMarketChangePrice   market `json:"postMarketPrice"`
	TickerName              string `json:"longName"`
	TickerCode              string `json:"symbol"`
	MarketState             string `json:"marketState"`
	CurrencyCode            string `json:"currency"`
	CurrencySymbol          string `json:"currencySymbol"`
}

type market struct {
	Raw float32 `json:"raw"`
	Fmt string  `json:"fmt"`
}

func New(tickerCode string) tickerPrice.TickerPrice {
	url := fmt.Sprintf("https://query2.finance.yahoo.com/v10/finance/quoteSummary/%s?modules=price", tickerCode)
	tickerDetails, _, APIError := queryAPIRaw(url)
	if APIError != nil {
		log.Print(APIError)
	}

	parsedPrice, parsingError := parsePriceData(tickerDetails)
	if parsingError != nil {
		log.Print(parsingError)
	}
	stockData := CreateStandarisedPriceObj(parsedPrice)
	return stockData
}

func CreateStandarisedPriceObj(parsedPrice Price) tickerPrice.TickerPrice {
	standardisedPrice := tickerPrice.TickerPrice{}

	standardisedPrice.SymbolCode = parsedPrice.TickerCode
	standardisedPrice.SymbolName = parsedPrice.TickerName
	standardisedPrice.RegMarket.Price = parsedPrice.RegMarketChangePrice.Fmt
	standardisedPrice.RegMarket.PriceChange = parsedPrice.RegMarketChange.Fmt
	standardisedPrice.RegMarket.PricePercentageChange = parsedPrice.RegMarketChangePercent.Fmt
	standardisedPrice.CurrencyCode = parsedPrice.CurrencyCode
	standardisedPrice.CurrencySymbol = parsedPrice.CurrencySymbol
	standardisedPrice.MarketState = parsedPrice.MarketState

	standardisedPrice.PreMarket = tickerPrice.TickerMarketOverview{}
	standardisedPrice.PostMarket = tickerPrice.TickerMarketOverview{}

	if parsedPrice.MarketState == "PRE" {
		standardisedPrice.PreMarket.Price = parsedPrice.PreMarketChangePrice.Fmt
		standardisedPrice.PreMarket.PriceChange = parsedPrice.PreMarketChange.Fmt
		standardisedPrice.PreMarket.PricePercentageChange = parsedPrice.PreMarketChangePercent.Fmt
	} else if parsedPrice.MarketState == "POST" {
		standardisedPrice.PostMarket.Price = parsedPrice.PostMarketChangePrice.Fmt
		standardisedPrice.PostMarket.PriceChange = parsedPrice.PostMarketChange.Fmt
		standardisedPrice.PostMarket.PricePercentageChange = parsedPrice.PostMarketChangePercent.Fmt
	}

	return standardisedPrice
}

func parsePriceData(responseData []byte) (Price, error) {
	var parsedPrice APIResponse
	parsingError := json.Unmarshal(responseData, &parsedPrice)
	if parsingError != nil {
		return Price{}, parsingError
	}
	return parsedPrice.Response.Result[0].Price, nil
}

func queryAPIRaw(url string) ([]byte, int, error) {
	yahooFinance := http.Client{Timeout: time.Second * 2}

	apiRequest, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	apiResponse, responseErr := yahooFinance.Do(apiRequest)
	if responseErr != nil {
		log.Fatal(responseErr)
	}

	responseBody, readError := ioutil.ReadAll(apiResponse.Body)
	if readError != nil {
		log.Fatal(readError)
	}

	return responseBody, apiResponse.StatusCode, nil
}
