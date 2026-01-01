package alert

import (
	"fmt"
)

type Alert struct {
	ID        uint    `gorm:"primaryKey"`
	Symbol    string  `gorm:"size:100;not null"`
	Condition string  `gorm:"size:1;not null"`
	Price     float64 `gorm:"not null"`
	Triggered bool    `gorm:"default:false"`
}

func NewAlert(symbol string, condition string, price float64) *Alert {
	return &Alert{
		Symbol:    symbol,
		Condition: condition,
		Price:     price,
	}
}

func ReadAlert() *Alert {
	var (
		nameCrypto      string
		conditionCrypto string
		priceCrypto     float64
	)
	fmt.Print("Введите название криптовалюты: ")
	fmt.Scan(&nameCrypto)
	fmt.Print("Введите условие (>/<): ")
	fmt.Scan(&conditionCrypto)
	fmt.Print("Введите цену для условия: ")
	fmt.Scan(&priceCrypto)

	return NewAlert(nameCrypto, conditionCrypto, priceCrypto)
}

func HasActiveAlerts(alerts []*Alert) bool {
	hasPending := false
	for _, a := range alerts {
		if !a.Triggered {
			hasPending = true
			break
		}
	}
	return hasPending
}
