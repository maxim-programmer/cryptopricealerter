package cli

import (
	"cryptopricealerter/internal/alert"
	"errors"
	"fmt"
	"strings"
)

func ReadAlert() *alert.Alert {
	var (
		nameCrypto      string
		conditionCrypto string
		priceCrypto     float64
	)
	fmt.Print("Введите название криптовалюты: ")
	fmt.Scan(&nameCrypto)
	nameCrypto = strings.ToLower(nameCrypto)
	fmt.Print("Введите условие (>/<): ")
	fmt.Scan(&conditionCrypto)
	fmt.Print("Введите цену для условия: ")
	fmt.Scan(&priceCrypto)

	return alert.NewAlert(nameCrypto, conditionCrypto, priceCrypto)
}

func ValidateAlert(alert *alert.Alert) error {
	if !(alert.Condition == ">" || alert.Condition == "<") {
		return errors.New("unknown condition")
	}
	if alert.Price < 0 {
		return errors.New("negative price")
	}
	if strings.Contains(alert.Symbol, " ") {
		return errors.New("space in symbol")
	}

	return nil
}
