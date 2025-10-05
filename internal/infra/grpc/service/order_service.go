package service

import (
	"context"

	"github.com/luuan11/clean-architecture/internal/infra/grpc/pb"
	"github.com/luuan11/clean-architecture/internal/usecase"
)

type OrderService struct {
	pb.UnimplementedOrderServiceServer
	CreateOrderUseCase  *usecase.CreateOrderUseCase
	ListOrdersUseCase   *usecase.ListOrdersUseCase
	GetOrderByIDUseCase *usecase.GetOrderByIDUseCase
	UpdateOrderUseCase  *usecase.UpdateOrderUseCase
	DeleteOrderUseCase  *usecase.DeleteOrderUseCase
}

func NewOrderService(
	createOrderUseCase *usecase.CreateOrderUseCase,
	listOrdersUseCase *usecase.ListOrdersUseCase,
	getOrderByIDUseCase *usecase.GetOrderByIDUseCase,
	updateOrderUseCase *usecase.UpdateOrderUseCase,
	deleteOrderUseCase *usecase.DeleteOrderUseCase,
) *OrderService {
	return &OrderService{
		CreateOrderUseCase:  createOrderUseCase,
		ListOrdersUseCase:   listOrdersUseCase,
		GetOrderByIDUseCase: getOrderByIDUseCase,
		UpdateOrderUseCase:  updateOrderUseCase,
		DeleteOrderUseCase:  deleteOrderUseCase,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	input := usecase.CreateOrderInputDTO{
		Price: req.Price,
		Tax:   req.Tax,
	}

	output, err := s.CreateOrderUseCase.Execute(input)
	if err != nil {
		return nil, err
	}

	return &pb.CreateOrderResponse{
		Id:         output.ID,
		Price:      output.Price,
		Tax:        output.Tax,
		FinalPrice: output.FinalPrice,
	}, nil
}

func (s *OrderService) ListOrders(ctx context.Context, req *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error) {
	orders, err := s.ListOrdersUseCase.Execute()
	if err != nil {
		return nil, err
	}

	var pbOrders []*pb.Order
	for _, order := range orders {
		pbOrders = append(pbOrders, &pb.Order{
			Id:         order.ID,
			Price:      order.Price,
			Tax:        order.Tax,
			FinalPrice: order.FinalPrice,
		})
	}

	return &pb.ListOrdersResponse{
		Orders: pbOrders,
	}, nil
}


