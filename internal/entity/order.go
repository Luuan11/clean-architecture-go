package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidPrice = errors.New("invalid price")
	ErrInvalidTax   = errors.New("invalid tax")
)

type Order struct {
	ID         string
	Price      float64
	Tax        float64
	FinalPrice float64
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func NewOrder(price, tax float64) (*Order, error) {
	if price <= 0 {
		return nil, ErrInvalidPrice
	}
	if tax < 0 {
		return nil, ErrInvalidTax
	}

	order := &Order{
		ID:         uuid.New().String(),
		Price:      price,
		Tax:        tax,
		FinalPrice: price + tax,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	return order, nil
}

func (o *Order) CalculateFinalPrice() {
	o.FinalPrice = o.Price + o.Tax
}
