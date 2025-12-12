package main

import (
	"context"
	"cryptopricealerter/internal/models"
	"cryptopricealerter/internal/pricefetcher"
	"cryptopricealerter/internal/repository"
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
		time_tick, _ = strconv.Atoi(os.Getenv("TIME_TICK"))
		api_key   = os.Getenv("API_KEY")
		user      = os.Getenv("POSTGRES_USER")
		password  = os.Getenv("POSTGRES_PASSWORD")
		dbname    = os.Getenv("POSTGRES_DB")
		host      = os.Getenv("DB_HOST")
		port      = os.Getenv("DB_PORT")
	)
	prices := make(map[string]pricefetcher.Price)
	fetcher := pricefetcher.NewFetcher()

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan bool, 1)

	ticker := time.NewTicker(time.Duration(time_tick) * time.Second)
	defer ticker.Stop()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT)

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(&models.Alert{})
	if err != nil {
		log.Fatal(err)
	}

	alertRepo := repository.NewAlertRepository(db)

	newAlert := &models.Alert{
		Symbol:    "bitcoin",
		Condition: ">",
		Price:     8000,
	}

	err = alertRepo.Create(newAlert)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Created new alert, ID =", newAlert.ID)

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
			alerts, err := alertRepo.GetAll()
			if err != nil {
				log.Fatal(err)
			}
			for _, alert := range alerts {
				if !alert.Triggered {
					prices, err := fetcher.GetPrices(ctx, alert.Symbol, prices, api_key)
					if err != nil {
						log.Fatal("Error:", err)
					}
					for _, v := range prices {
						switch alert.Condition {
						case ">":
							if v.USD > alert.Price {
								_ = alertRepo.MarkTriggered(alert.ID)
								fmt.Println("Триггер отработал!")
							}
						case "<":
							if v.USD < alert.Price {
								_ = alertRepo.MarkTriggered(alert.ID)
								fmt.Println("Триггер отработал!")
							}
						}
					}
				} else {
					_ = alertRepo.Delete(alert.ID)
					fmt.Println("Алерт удален!")
				}
			}

		case <-ctx.Done():
			return
		}
	}
}
