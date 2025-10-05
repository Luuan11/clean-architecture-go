package mocks

import (
	"errors"

	"github.com/luuan11/clean-architecture/internal/entity"
)

type OrderRepositoryMock struct {
	orders      []*entity.Order
	SaveError   error
	FindError   error
	UpdateError error
	DeleteError error
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

func (m *OrderRepositoryMock) FindByID(id string) (*entity.Order, error) {
	if m.FindError != nil {
		return nil, m.FindError
	}
	
	for _, order := range m.orders {
		if order.ID == id {
			return order, nil
		}
	}
	
	return nil, entity.ErrOrderNotFound
}

func (m *OrderRepositoryMock) Update(order *entity.Order) error {
	if m.UpdateError != nil {
		return m.UpdateError
	}
	
	for i, o := range m.orders {
		if o.ID == order.ID {
			m.orders[i] = order
			return nil
		}
	}
	
	return entity.ErrOrderNotFound
}

func (m *OrderRepositoryMock) Delete(id string) error {
	if m.DeleteError != nil {
		return m.DeleteError
	}
	
	for i, order := range m.orders {
		if order.ID == id {
			m.orders = append(m.orders[:i], m.orders[i+1:]...)
			return nil
		}
	}
	
	return entity.ErrOrderNotFound
}

func (m *OrderRepositoryMock) SetSaveError(err error) {
	m.SaveError = err
}

func (m *OrderRepositoryMock) SetFindError(err error) {
	m.FindError = err
}

func (m *OrderRepositoryMock) SetUpdateError(err error) {
	m.UpdateError = err
}

func (m *OrderRepositoryMock) SetDeleteError(err error) {
	m.DeleteError = err
}

func (m *OrderRepositoryMock) Reset() {
	m.orders = make([]*entity.Order, 0)
	m.SaveError = nil
	m.FindError = nil
	m.UpdateError = nil
	m.DeleteError = nil
}

var (
	ErrRepositoryFailure  = errors.New("repository failure")
	ErrDatabaseConnection = errors.New("database connection failed")
)
