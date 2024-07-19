package external

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
)

type PriceService struct {
	Ch        chan SymbolPrice
	LastPrice LastPrice
}

type LastPrice struct {
	Prices map[string]string `json:"prices"`
	mu     *sync.RWMutex
}

type SymbolPrice struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

const urlGetPrice = "https://api.binance.com/api/v3/ticker/price"

func NewPriceService() *PriceService {
	ch := make(chan SymbolPrice, 5)
	return &PriceService{
		Ch: ch,
		LastPrice: LastPrice{
			Prices: make(map[string]string),
			mu:     &sync.RWMutex{},
		},
	}
}

func (s *PriceService) GetPrice(symbol string) error {
	params := url.Values{}
	params.Add("symbol", symbol)
	finalURL := fmt.Sprintf("%s?%s", urlGetPrice, params.Encode())

	resp, err := http.Get(finalURL)
	if err != nil {
		return fmt.Errorf("error create request for %s: %v", symbol, err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error read response for %s: %v", symbol, err)
	}

	var price SymbolPrice
	err = json.Unmarshal(body, &price)
	if err != nil {
		return fmt.Errorf("error unmarshal JSON for %s: %v", symbol, err)
	}
	currentPrice, ok := s.LastPrice.Get(price.Symbol)
	if ok && currentPrice == price.Price {
		return nil
	}
	fmt.Printf("Symbol: %s, Price: %s\n", price.Symbol, price.Price)
	s.Ch <- price

	return nil
}

func (s *PriceService) ProcessPrices() {
	for price := range s.Ch {
		s.LastPrice.Update(price.Symbol, price.Price)
	}
}

func (lp *LastPrice) GetAll() map[string]string {
	lp.mu.RLock()
	prices := lp.Prices
	lp.mu.RUnlock()

	return prices
}

func (lp *LastPrice) Get(symbol string) (string, bool) {
	lp.mu.RLock()
	price, ok := lp.Prices[symbol]
	lp.mu.RUnlock()

	return price, ok
}

func (lp *LastPrice) Update(symbol, price string) {
	lp.mu.Lock()
	lp.Prices[symbol] = price
	lp.mu.Unlock()
}
