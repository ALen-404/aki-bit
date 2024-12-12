package types

type CoinMarketInfo struct {
	ID                int     `json:"id"`
	Name              string  `json:"name"`
	Symbol            string  `json:"symbol"`
	Slug              string  `json:"slug"`
	CmcRank           int     `json:"cmc_rank"`
	MarketPairs       int     `json:"num_market_pairs"`
	CirculatingSupply float64 `json:"circulating_supply"`
	TotalSupply       float64 `json:"total_supply"`
	MaxSupply         float64 `json:"max_supply"`
	LastUpdated       string  `json:"last_updated"`

	Quote map[string]CoinMarketQuote `json:"quote"`
}

type CoinMarketQuote struct {
	Price            float64 `json:"price"`
	Volume24h        float64 `json:"volume_24h"`
	Volume7d         float64 `json:"volume_7d"`
	Volume30d        float64 `json:"volume_30d"`
	MarketCap        float64 `json:"market_cap"`
	FullyMarketCap   float64 `json:"fully_diluted_market_cap"`
	PercentChange1h  float64 `json:"percent_change_1h"`
	PercentChange24h float64 `json:"percent_change_24h"`
	PercentChange7d  float64 `json:"percent_change_7d"`
	LastUpdated      string  `json:"last_updated"`
}

type CategoryMarketData struct {
	TokenNum        int     `json:"num_tokens"`
	AvgPriceChange  float64 `json:"avg_price_change"`
	MarketCap       float64 `json:"market_cap"`
	MarketCapChange float64 `json:"market_cap_change"`
	Volume          float64 `json:"volume"`
	VolumeChange    float64 `json:"volume_change"`
	LastUpdated     string  `json:"last_updated"`

	Coins []CoinMarketInfo `json:"coins"`
}

type CategoryMarketListData struct {
	TokenNum        int     `json:"num_tokens"`
	AvgPriceChange  float64 `json:"avg_price_change"`
	MarketCap       float64 `json:"market_cap"`
	MarketCapChange float64 `json:"market_cap_change"`
	Volume          float64 `json:"volume"`
	VolumeChange    float64 `json:"volume_change"`
	LastUpdated     string  `json:"last_updated"`
}

type OhlcvSeriesData struct {
	ID     int         `json:"id"`
	Name   string      `json:"name"`
	Symbol string      `json:"symbol"`
	Quotes []OhlcvData `json:"quotes"`
}

type OhlcvData struct {
	OpenTime  string `json:"time_open"`
	CloseTime string `json:"time_close"`
	HighTime  string `json:"time_high"`
	LowTime   string `json:"time_low"`

	Quotes map[string]OhlcvQuote `json:"quote"`
}

type OhlcvQuote struct {
	Open      float64 `json:"open"`
	High      float64 `json:"high"`
	Low       float64 `json:"low"`
	Close     float64 `json:"close"`
	Volume    float64 `json:"volume"`
	MarketCap float64 `json:"market_cap"`
	Timestamp string  `json:"timestamp"`
}

type CmcStatus struct {
	Timestamp  string `json:"timestamp"`
	ErrCode    int    `json:"error_code"`
	ErrMessage string `json:"error_message"`
	Elapsed    int    `json:"elapsed"`
	UseCredit  int    `json:"credit_count"`
	Notice     string `json:"notice"`
}

type CmcError struct {
	Status CmcStatus `json:"status"`
}

func (e CmcError) Error() string {
	return e.Status.ErrMessage
}
