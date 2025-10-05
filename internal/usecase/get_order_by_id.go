package usecase

import (
	"github.com/luuan11/clean-architecture/internal/entity"
)

type GetOrderByIDInputDTO struct {
	ID string `json:"id"`
}

type GetOrderByIDOutputDTO struct {
	ID         string  `json:"id"`
	Price      float64 `json:"price"`
	Tax        float64 `json:"tax"`
	FinalPrice float64 `json:"final_price"`
}

type GetOrderByIDUseCase struct {
	OrderRepository entity.OrderRepositoryInterface
}

func NewGetOrderByIDUseCase(orderRepository entity.OrderRepositoryInterface) *GetOrderByIDUseCase {
	return &GetOrderByIDUseCase{
		OrderRepository: orderRepository,
	}
}

func (g *GetOrderByIDUseCase) Execute(input GetOrderByIDInputDTO) (*GetOrderByIDOutputDTO, error) {
	order, err := g.OrderRepository.FindByID(input.ID)
	if err != nil {
		return nil, err
	}

	return &GetOrderByIDOutputDTO{
		ID:         order.ID,
		Price:      order.Price,
		Tax:        order.Tax,
		FinalPrice: order.FinalPrice,
	}, nil
}
