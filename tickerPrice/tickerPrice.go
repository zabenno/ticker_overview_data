package tickerPrice

type TickerPrice struct {
	SymbolCode     string
	SymbolName     string
	RegMarket      TickerMarketOverview
	PreMarket      TickerMarketOverview
	PostMarket     TickerMarketOverview
	CurrencyCode   string
	CurrencySymbol string
	MarketState    string
}

type TickerMarketOverview struct {
	Price                 string
	PriceChange           string
	PricePercentageChange string
}
