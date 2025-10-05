package usecase

import (
	"github.com/luuan11/clean-architecture/internal/entity"
)

type DeleteOrderInputDTO struct {
	ID string `json:"id"`
}

type DeleteOrderUseCase struct {
	OrderRepository entity.OrderRepositoryInterface
}

func NewDeleteOrderUseCase(orderRepository entity.OrderRepositoryInterface) *DeleteOrderUseCase {
	return &DeleteOrderUseCase{
		OrderRepository: orderRepository,
	}
}

func (d *DeleteOrderUseCase) Execute(input DeleteOrderInputDTO) error {
	_, err := d.OrderRepository.FindByID(input.ID)
	if err != nil {
		return err
	}

	return d.OrderRepository.Delete(input.ID)
}
