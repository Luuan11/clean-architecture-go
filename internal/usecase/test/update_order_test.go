package usecase_test

import (
	"testing"

	"github.com/luuan11/clean-architecture/internal/entity"
	"github.com/luuan11/clean-architecture/internal/usecase"
	"github.com/luuan11/clean-architecture/internal/usecase/test/mocks"
)

func TestUpdateOrderUseCase_Success(t *testing.T) {
	mockRepo := mocks.NewOrderRepositoryMock()

	order, _ := entity.NewOrder(100.0, 10.0)
	mockRepo.Save(order)

	uc := usecase.NewUpdateOrderUseCase(mockRepo)

	input := usecase.UpdateOrderInputDTO{
		ID:    order.ID,
		Price: 200.0,
		Tax:   20.0,
	}

	output, err := uc.Execute(input)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if output.Price != 200.0 {
		t.Errorf("Expected price 200.0, got %.2f", output.Price)
	}

	if output.Tax != 20.0 {
		t.Errorf("Expected tax 20.0, got %.2f", output.Tax)
	}

	if output.FinalPrice != 220.0 {
		t.Errorf("Expected final price 220.0, got %.2f", output.FinalPrice)
	}
}

func TestUpdateOrderUseCase_WithDifferentValues(t *testing.T) {
	tests := []struct {
		name               string
		initialPrice       float64
		initialTax         float64
		newPrice           float64
		newTax             float64
		expectedFinalPrice float64
	}{
		{
			name:               "Small to large values",
			initialPrice:       10.0,
			initialTax:         1.0,
			newPrice:           1000.0,
			newTax:             100.0,
			expectedFinalPrice: 1100.0,
		},
		{
			name:               "Large to small values",
			initialPrice:       1000.0,
			initialTax:         100.0,
			newPrice:           50.0,
			newTax:             5.0,
			expectedFinalPrice: 55.0,
		},
		{
			name:               "Same values",
			initialPrice:       100.0,
			initialTax:         10.0,
			newPrice:           100.0,
			newTax:             10.0,
			expectedFinalPrice: 110.0,
		},
		{
			name:               "Decimal values",
			initialPrice:       100.0,
			initialTax:         10.0,
			newPrice:           99.99,
			newTax:             9.99,
			expectedFinalPrice: 109.98,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewOrderRepositoryMock()

			order, _ := entity.NewOrder(tt.initialPrice, tt.initialTax)
			mockRepo.Save(order)

			uc := usecase.NewUpdateOrderUseCase(mockRepo)

			input := usecase.UpdateOrderInputDTO{
				ID:    order.ID,
				Price: tt.newPrice,
				Tax:   tt.newTax,
			}

			output, err := uc.Execute(input)

			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			const epsilon = 0.01
			if diff := output.FinalPrice - tt.expectedFinalPrice; diff > epsilon || diff < -epsilon {
				t.Errorf("Expected final price %.2f, got %.2f", tt.expectedFinalPrice, output.FinalPrice)
			}
		})
	}
}

func TestUpdateOrderUseCase_OrderNotFound(t *testing.T) {
	mockRepo := mocks.NewOrderRepositoryMock()
	uc := usecase.NewUpdateOrderUseCase(mockRepo)

	input := usecase.UpdateOrderInputDTO{
		ID:    "non-existent-id",
		Price: 100.0,
		Tax:   10.0,
	}

	_, err := uc.Execute(input)

	if err != entity.ErrOrderNotFound {
		t.Errorf("Expected ErrOrderNotFound, got %v", err)
	}
}

func TestUpdateOrderUseCase_InvalidPrice(t *testing.T) {
	tests := []struct {
		name  string
		price float64
		tax   float64
	}{
		{"Negative price", -100.0, 10.0},
		{"Zero price", 0.0, 10.0},
		{"Very negative price", -999999.0, 10.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewOrderRepositoryMock()

			order, _ := entity.NewOrder(100.0, 10.0)
			mockRepo.Save(order)

			uc := usecase.NewUpdateOrderUseCase(mockRepo)

			input := usecase.UpdateOrderInputDTO{
				ID:    order.ID,
				Price: tt.price,
				Tax:   tt.tax,
			}

			_, err := uc.Execute(input)

			if err != entity.ErrInvalidPrice {
				t.Errorf("Expected ErrInvalidPrice, got %v", err)
			}
		})
	}
}

func TestUpdateOrderUseCase_InvalidTax(t *testing.T) {
	mockRepo := mocks.NewOrderRepositoryMock()

	order, _ := entity.NewOrder(100.0, 10.0)
	mockRepo.Save(order)

	uc := usecase.NewUpdateOrderUseCase(mockRepo)

	input := usecase.UpdateOrderInputDTO{
		ID:    order.ID,
		Price: 100.0,
		Tax:   -10.0,
	}

	_, err := uc.Execute(input)

	if err != entity.ErrInvalidTax {
		t.Errorf("Expected ErrInvalidTax, got %v", err)
	}
}

func TestUpdateOrderUseCase_RepositoryFindError(t *testing.T) {
	mockRepo := mocks.NewOrderRepositoryMock()
	mockRepo.SetFindError(mocks.ErrDatabaseConnection)

	uc := usecase.NewUpdateOrderUseCase(mockRepo)

	input := usecase.UpdateOrderInputDTO{
		ID:    "some-id",
		Price: 100.0,
		Tax:   10.0,
	}

	_, err := uc.Execute(input)

	if err != mocks.ErrDatabaseConnection {
		t.Errorf("Expected ErrDatabaseConnection, got %v", err)
	}
}

func TestUpdateOrderUseCase_RepositoryUpdateError(t *testing.T) {
	mockRepo := mocks.NewOrderRepositoryMock()

	order, _ := entity.NewOrder(100.0, 10.0)
	mockRepo.Save(order)

	mockRepo.SetUpdateError(mocks.ErrRepositoryFailure)

	uc := usecase.NewUpdateOrderUseCase(mockRepo)

	input := usecase.UpdateOrderInputDTO{
		ID:    order.ID,
		Price: 200.0,
		Tax:   20.0,
	}

	_, err := uc.Execute(input)

	if err != mocks.ErrRepositoryFailure {
		t.Errorf("Expected ErrRepositoryFailure, got %v", err)
	}
}

func TestUpdateOrderUseCase_NilRepository(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic with nil repository")
		}
	}()

	uc := usecase.NewUpdateOrderUseCase(nil)

	input := usecase.UpdateOrderInputDTO{
		ID:    "some-id",
		Price: 100.0,
		Tax:   10.0,
	}

	uc.Execute(input)
}
