package money

type Money struct {
	cents int
}

func (m *Money) GetAsEuro() float64 {
	return float64(m.cents / 100)
}

func (m *Money) GetAsCents() int {
	return m.cents
}

func (m *Money) Add(money Money) Money {
	return NewMoney(
		m.cents + money.cents,
	)
}

func (m *Money) Sub(money Money) Money {
	return NewMoney(
		m.cents - m.cents,
	)
}

func (m *Money) Multi(factor int) Money{
	return	NewMoney(
		m.cents * factor,
	)
}

func NewMoney(cents int) Money {
	return Money{
		cents: cents,
	}
}
