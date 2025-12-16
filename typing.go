package main

type Level interface {
	Price() float64
	Size() float64
	OrderSize() float64
	Type() string
}

type OrderBookLevel struct {
	price     float64
	size      float64
	orderSize float64
	name      string
}

func (this *OrderBookLevel) Price() float64 {
	return this.price
}
func (this *OrderBookLevel) Size() float64      { return this.size }
func (this *OrderBookLevel) OrderSize() float64 { return this.orderSize }
func (this *OrderBookLevel) Type() string       { return this.name }

func NewOrderBookLevel(price, size, orderSize float64, isBid bool) *OrderBookLevel {
	return &OrderBookLevel{
		price:     price,
		size:      size,
		orderSize: orderSize,
		name: func() string {
			if isBid {
				return "bids"
			}
			return "asks"
		}(),
	}
}
