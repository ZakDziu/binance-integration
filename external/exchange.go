package external

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const urlGetExchangeInfo = "https://api.binance.com/api/v3/exchangeInfo"

type ExchangeInfo struct {
	Symbols []Symbol `json:"symbols"`
}

type Symbol struct {
	Symbol string `json:"symbol"`
}

func NewExchangeInfo() *ExchangeInfo {
	symbols := make([]Symbol, 0, 5)

	return &ExchangeInfo{
		Symbols: symbols,
	}
}

func (s *BinanceService) ProcessExchangeInfo() error {
	resp, err := http.Get(urlGetExchangeInfo)
	if err != nil {
		return fmt.Errorf("error get exchange info: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error read response: %v", err)
	}

	var exchangeInfo ExchangeInfo
	err = json.Unmarshal(body, &exchangeInfo)
	if err != nil {
		return fmt.Errorf("error unmarshal JSON: %v", err)
	}

	for i, symbol := range exchangeInfo.Symbols {
		if i >= 5 {
			break
		}
		s.ExchangeInfo.Symbols = append(s.ExchangeInfo.Symbols, symbol)
	}

	return nil
}
