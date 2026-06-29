package money

type Money struct {
	cents int
}

func (m Money) GetAsEuro() float64 {
	return float64(m.cents / 100)
}

func (m Money) GetAsCents() int {
	return m.cents
}

func (m Money) Add(money Money) Money {
	return MoneyFrom(
		m.cents + money.cents,
	)
}

func (m Money) Sub(money Money) Money {
	return MoneyFrom(
		m.cents - m.cents,
	)
}

func (m Money) Multi(factor int) Money {
	return MoneyFrom(
		m.cents * factor,
	)
}

func MoneyFrom(cents int) Money {
	return Money{
		cents: cents,
	}
}
