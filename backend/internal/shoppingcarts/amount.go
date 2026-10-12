package shoppingcarts

type Amount struct {
	payedAmount int
	freeAmount  int
}

/*
	func NewAmount() Amount {
		return Amount{
			payedAmount: 0,
			freeAmount:  0,
		}
	}
*/
func AmountFrom(payedAmount int, freeAmount int) Amount {
	return Amount{
		payedAmount: payedAmount,
		freeAmount:  freeAmount,
	}
}

func (a *Amount) GetFreeAmount() int {
	return a.freeAmount
}

func (a *Amount) GetPayedAmount() int {
	return a.payedAmount
}

func (a *Amount) GetTotalAmount() int {
	return a.payedAmount + a.freeAmount
}

func (a *Amount) IncreaseFreeAmount() {
	a.freeAmount++

}

func (a *Amount) IncreasePayedAmount() {
	a.payedAmount++
}

// First decrease payed products then free products
func (a *Amount) DecreaseAmount() {

	if a.payedAmount > 0 {
		a.payedAmount--
	} else if a.freeAmount > 0 {
		a.freeAmount--
	}
}
