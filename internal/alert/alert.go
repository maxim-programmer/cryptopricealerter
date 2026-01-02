package alert

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

func HasActiveAlerts(alerts []*Alert) bool {
	for _, a := range alerts {
		if !a.Triggered {
			return true
		}
	}
	return false
}
