package external

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type BinanceService struct {
	ExchangeInfo *ExchangeInfo `json:"exchangeInfo"`
	PriceService *PriceService `json:"priceService"`
}

func NewBinanceService(ctx context.Context) (*BinanceService, error) {
	service := &BinanceService{
		ExchangeInfo: NewExchangeInfo(),
		PriceService: NewPriceService(),
	}
	err := service.ProcessExchangeInfo()
	if err != nil {
		return nil, err
	}

	return service, nil
}

func (s *BinanceService) ProcessGetPrices(ctx context.Context, globalWG *sync.WaitGroup) {
	defer globalWG.Done()
	var wg sync.WaitGroup

	for _, symbol := range s.ExchangeInfo.Symbols {
		wg.Add(1)
		symbol := symbol // Create a new variable i for each goroutine to avoid problems with closure
		go func() {
			defer wg.Done()
			ticker := time.NewTicker(1 * time.Second)

			for {
				select {
				case <-ctx.Done():
					return
				case <-ticker.C:
					err := s.PriceService.GetPrice(symbol.Symbol)
					if err != nil {
						fmt.Println(err)
					}
				}
			}
		}()
	}

	wg.Wait()
}
