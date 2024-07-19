package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"binance-integrate/external"
	"binance-integrate/pkg/api"
	"binance-integrate/pkg/config"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	conf, err := config.New()
	if err != nil {
		log.Fatalf("Can't read config file: %s", err)
	}

	binanceService, err := external.NewBinanceService(ctx)
	if err != nil {
		log.Fatalf("Can't create binance service: %s", err)
	}

	apiServer := api.NewServer(conf.Server, binanceService)

	var wg sync.WaitGroup

	wg.Add(1)
	go binanceService.ProcessGetPrices(ctx, &wg)
	go binanceService.PriceService.ProcessPrices()

	runErr := make(chan error, 1)
	quitCh := make(chan os.Signal, 1)
	signal.Notify(quitCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err = apiServer.ListenAndServe()
		if err != nil {
			runErr <- fmt.Errorf("can't start http server: %w", err)
		}
	}()

	select {
	case err = <-runErr:
		cancel()
		wg.Wait()
		close(binanceService.PriceService.Ch)

		log.Fatalf("Running error: %s", err)
	case s := <-quitCh:
		cancel()
		wg.Wait()

		close(binanceService.PriceService.Ch)

		log.Printf("Received signal: %v. Running graceful shutdown...", s)
		ctx := context.Background()

		err = apiServer.Shutdown(ctx)
		if err != nil {
			log.Printf("Can't shutdown API server: %s", err)
		}
	}

}
