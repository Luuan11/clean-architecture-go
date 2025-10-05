package usecase

import (
	"github.com/luuan11/clean-architecture/internal/entity"
)

type UpdateOrderInputDTO struct {
	ID    string  `json:"id"`
	Price float64 `json:"price"`
	Tax   float64 `json:"tax"`
}

type UpdateOrderOutputDTO struct {
	ID         string  `json:"id"`
	Price      float64 `json:"price"`
	Tax        float64 `json:"tax"`
	FinalPrice float64 `json:"final_price"`
}

type UpdateOrderUseCase struct {
	OrderRepository entity.OrderRepositoryInterface
}

func NewUpdateOrderUseCase(orderRepository entity.OrderRepositoryInterface) *UpdateOrderUseCase {
	return &UpdateOrderUseCase{
		OrderRepository: orderRepository,
	}
}

func (u *UpdateOrderUseCase) Execute(input UpdateOrderInputDTO) (*UpdateOrderOutputDTO, error) {
	order, err := u.OrderRepository.FindByID(input.ID)
	if err != nil {
		return nil, err
	}

	err = order.Update(input.Price, input.Tax)
	if err != nil {
		return nil, err
	}

	err = u.OrderRepository.Update(order)
	if err != nil {
		return nil, err
	}

	return &UpdateOrderOutputDTO{
		ID:         order.ID,
		Price:      order.Price,
		Tax:        order.Tax,
		FinalPrice: order.FinalPrice,
	}, nil
}
