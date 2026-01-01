package workerpool

type Job struct {
	ID uint
	Symbol string
	Condition string
	Price float64
	Triggered bool
	ActualPrice float64
}

func NewJob(id uint, symbol string, condition string, price float64, triggered bool, actualPrice float64) Job {
	return Job{
		ID: id,
		Symbol: symbol,
		Condition: condition,
		Price: price,
		Triggered: triggered,
		ActualPrice: actualPrice,
	}
}