package mocks

import (
	"errors"

	"github.com/luuan11/clean-architecture/internal/entity"
)

type OrderRepositoryMock struct {
	orders    []*entity.Order
	SaveError error
	FindError error
}

func NewOrderRepositoryMock() *OrderRepositoryMock {
	return &OrderRepositoryMock{
		orders: make([]*entity.Order, 0),
	}
}

func (m *OrderRepositoryMock) Save(order *entity.Order) error {
	if m.SaveError != nil {
		return m.SaveError
	}
	m.orders = append(m.orders, order)
	return nil
}

func (m *OrderRepositoryMock) FindAll() ([]*entity.Order, error) {
	if m.FindError != nil {
		return nil, m.FindError
	}
	return m.orders, nil
}

func (m *OrderRepositoryMock) SetSaveError(err error) {
	m.SaveError = err
}

func (m *OrderRepositoryMock) SetFindError(err error) {
	m.FindError = err
}

func (m *OrderRepositoryMock) Reset() {
	m.orders = make([]*entity.Order, 0)
	m.SaveError = nil
	m.FindError = nil
}

var (
	ErrRepositoryFailure  = errors.New("repository failure")
	ErrDatabaseConnection = errors.New("database connection failed")
)
