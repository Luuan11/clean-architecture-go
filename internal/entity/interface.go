package entity

type OrderRepositoryInterface interface {
	Save(order *Order) error
	FindAll() ([]*Order, error)
	FindByID(id string) (*Order, error)
	Update(order *Order) error
	Delete(id string) error
}
