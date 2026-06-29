package money

type Money struct {
	Cents int
}

func (m Money) GetAsEuro() float32 {
	return float32(m.Cents) / 100
}

func (m Money) GetAsCents() int {
	return m.Cents
}

func (m Money) Add(money Money) Money {
	return MoneyFrom(
		m.Cents + money.Cents,
	)
}

func (m Money) Sub(money Money) Money {
	return MoneyFrom(
		m.Cents - m.Cents,
	)
}

func (m Money) Multi(factor int) Money {
	return MoneyFrom(
		m.Cents * factor,
	)
}

func MoneyFrom(cents int) Money {
	return Money{
		Cents: cents,
	}
}
