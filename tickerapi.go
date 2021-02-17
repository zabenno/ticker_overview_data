package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/zabenno/ticker_overview_data/yahooFinanceAPI"
)

func main() {
	ticker := yahooFinanceAPI.New("ARKG")
	log.Print(ticker.SymbolName)

}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the homepage!")
}
