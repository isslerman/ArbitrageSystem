package cex

type Wallet struct {
	Assets map[string]Asset
}

type Asset struct {
	Symbol  string
	Balance Balance
}

type Balance struct {
	Amount float64
	Locked float64
}

func NewAsset(symbol string) Asset {
	return Asset{
		Symbol:  symbol,
		Balance: Balance{},
	}
}

func (b *Balance) Total() float64 {
	return b.Amount + b.Locked
}
