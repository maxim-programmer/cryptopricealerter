package main

import (
	"context"
	"cryptopricealerter/internal/pricefetcher"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	var (
		time_tick int
		coins     string
		api_key   string = os.Getenv("API_KEY")
	)
	prices := make(map[string]pricefetcher.Price)
	fetcher := pricefetcher.NewFetcher()

	fmt.Print("Введите время в секундах для интервала отправки запросов: ")
	_, _ = fmt.Scan(&time_tick)
	fmt.Print("Введите названия криптовалюты (если их несколько - через запятую, без пробела): ")
	_, _ = fmt.Scan(&coins)

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan bool, 1)

	ticker := time.NewTicker(time.Duration(time_tick) * time.Second)
	defer ticker.Stop()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT)

	go func() {
		<-sigs
		done <- true
	}()

	for {
		select {
		case <-done:
			cancel()
			fmt.Println("\nПрограмма завершила работу!")

		case <-ticker.C:
			prices, err := fetcher.GetPrices(ctx, coins, prices, api_key)
			if err != nil {
				fmt.Println("Error:", err)
			}
			for i, v := range prices {
				fmt.Println(i, v.USD)
			}

		case <-ctx.Done():
			return
		}
	}
}
