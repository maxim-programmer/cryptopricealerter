package main

import (
	"context"
	"cryptopricealerter/internal/alert"
	"cryptopricealerter/internal/cli"
	"cryptopricealerter/internal/pricefetcher"
	"cryptopricealerter/internal/repository"
	"cryptopricealerter/internal/workerpool"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	var (
		time_tick, _   = strconv.Atoi(os.Getenv("TIME_TICK"))
		api_key        = os.Getenv("API_KEY")
		user           = os.Getenv("POSTGRES_USER")
		password       = os.Getenv("POSTGRES_PASSWORD")
		dbname         = os.Getenv("POSTGRES_DB")
		host           = os.Getenv("DB_HOST")
		port           = os.Getenv("DB_PORT")
		chanSize, _    = strconv.Atoi(os.Getenv("CHAN_SIZE"))
		workerCount, _ = strconv.Atoi(os.Getenv("WORKER_COUNT"))
	)

	if time_tick <= 0 {
		log.Fatal("TIME_TICK must be > 0")
	}
	if workerCount <= 0 {
		log.Fatal("WORKER_COUNT must be > 0")
	}
	
	prices := make(map[string]pricefetcher.Price)
	fetcher := pricefetcher.NewFetcher()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	if err := db.AutoMigrate(&alert.Alert{}); err != nil {
		log.Fatal(err)
	}

	alertRepo := repository.NewAlertRepository(db)
	pool := workerpool.NewWorkerPool(chanSize, workerCount, alertRepo, ctx, stop)
	pool.Start()

	fmt.Println("=====Создание алертов=====")
	for {
		newAlert := cli.ReadAlert()
		if err := cli.ValidateAlert(newAlert); err != nil {
			log.Println("Error:", err)
			continue
		}

		if err := alertRepo.Create(newAlert); err != nil {
			log.Println("Error:", err)
			continue
		}
		fmt.Println("=====Created new alert, ID =", strconv.Itoa(int(newAlert.ID))+"======")

		var exitFlag string
		fmt.Print("Закончить создание алертов? (yes -> закончить): ")
		fmt.Scan(&exitFlag)
		if exitFlag == "yes" {
			fmt.Println("=====Завершение создания алертов=====")
			break
		}
	}

	ticker := time.NewTicker(time.Duration(time_tick) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			alerts, err := alertRepo.GetAll()
			if err != nil {
				log.Println("Error:", err)
				continue
			}

			if !alert.HasActiveAlerts(alerts) {
				pool.Stop()
				fmt.Println("Активные алерты закончились!")
				return
			}

			symbols := make([]string, 0)
			for _, alert := range alerts {
				if alert.Triggered {
					continue
				}
				flag := false
				for _, symbol := range symbols {
					if symbol == alert.Symbol {
						flag = true
						break
					}
				}
				if !flag {
					symbols = append(symbols, alert.Symbol)
				}
			}

			prices, err = fetcher.GetPrices(ctx, symbols, api_key)
			if err != nil {
				log.Println("Error:", err)
				continue
			}

			for _, alert := range alerts {
				if !alert.Triggered {
					actualPrice, ok := prices[alert.Symbol]
					if !ok {
						fmt.Println("Неверное название криптовалюты:", alert.Symbol)
						continue
					}
					job := workerpool.NewJob(alert.ID, alert.Symbol, alert.Condition, alert.Price, alert.Triggered, actualPrice.USD)
					pool.JobChan <- job
				}
			}
		case <-ctx.Done():
			pool.Stop()
			fmt.Println("\nПрограмма завершила работу!")
			return
		}
	}
}
