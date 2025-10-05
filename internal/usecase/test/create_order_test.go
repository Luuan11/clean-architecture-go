package usecase_test

import (
	"errors"
	"testing"

	"github.com/luuan11/clean-architecture/internal/entity"
	"github.com/luuan11/clean-architecture/internal/usecase"
	"github.com/luuan11/clean-architecture/internal/usecase/test/mocks"
)

func TestCreateOrderUseCase_Success(t *testing.T) {
	mockRepo := mocks.NewOrderRepositoryMock()
	uc := usecase.NewCreateOrderUseCase(mockRepo)

	input := usecase.CreateOrderInputDTO{
		Price: 100.0,
		Tax:   10.0,
	}

	output, err := uc.Execute(input)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if output == nil {
		t.Fatal("Expected output, got nil")
	}

	if output.Price != 100.0 {
		t.Errorf("Expected price 100.0, got %f", output.Price)
	}

	if output.Tax != 10.0 {
		t.Errorf("Expected tax 10.0, got %f", output.Tax)
	}

	if output.FinalPrice != 110.0 {
		t.Errorf("Expected final price 110.0, got %f", output.FinalPrice)
	}

	if output.ID == "" {
		t.Error("Expected ID to be set")
	}
}

func TestCreateOrderUseCase_WithDifferentValues(t *testing.T) {
	tests := []struct {
		name               string
		price              float64
		tax                float64
		expectedFinalPrice float64
	}{
		{
			name:               "Small values",
			price:              10.0,
			tax:                1.0,
			expectedFinalPrice: 11.0,
		},
		{
			name:               "Large values",
			price:              1000.0,
			tax:                100.0,
			expectedFinalPrice: 1100.0,
		},
		{
			name:               "Decimal values",
			price:              99.99,
			tax:                9.99,
			expectedFinalPrice: 109.98,
		},
		{
			name:               "Zero tax",
			price:              50.0,
			tax:                0.0,
			expectedFinalPrice: 50.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewOrderRepositoryMock()
			uc := usecase.NewCreateOrderUseCase(mockRepo)

			input := usecase.CreateOrderInputDTO{
				Price: tt.price,
				Tax:   tt.tax,
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

func TestCreateOrderUseCase_InvalidPrice(t *testing.T) {
	tests := []struct {
		name  string
		price float64
		tax   float64
	}{
		{
			name:  "Negative price",
			price: -100.0,
			tax:   10.0,
		},
		{
			name:  "Zero price",
			price: 0.0,
			tax:   10.0,
		},
		{
			name:  "Very negative price",
			price: -999999.99,
			tax:   10.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewOrderRepositoryMock()
			uc := usecase.NewCreateOrderUseCase(mockRepo)

			input := usecase.CreateOrderInputDTO{
				Price: tt.price,
				Tax:   tt.tax,
			}

			output, err := uc.Execute(input)

			if err == nil {
				t.Fatal("Expected error for invalid price, got nil")
			}

			if output != nil {
				t.Error("Expected nil output on error")
			}

			if !errors.Is(err, entity.ErrInvalidPrice) {
				t.Errorf("Expected ErrInvalidPrice, got %v", err)
			}
		})
	}
}

func TestCreateOrderUseCase_InvalidTax(t *testing.T) {
	mockRepo := mocks.NewOrderRepositoryMock()
	uc := usecase.NewCreateOrderUseCase(mockRepo)

	input := usecase.CreateOrderInputDTO{
		Price: 100.0,
		Tax:   -10.0,
	}

	output, err := uc.Execute(input)

	if err == nil {
		t.Fatal("Expected error for invalid tax, got nil")
	}

	if output != nil {
		t.Error("Expected nil output on error")
	}

	if !errors.Is(err, entity.ErrInvalidTax) {
		t.Errorf("Expected ErrInvalidTax, got %v", err)
	}
}

func TestCreateOrderUseCase_RepositoryError(t *testing.T) {
	mockRepo := mocks.NewOrderRepositoryMock()
	mockRepo.SetSaveError(mocks.ErrRepositoryFailure)
	uc := usecase.NewCreateOrderUseCase(mockRepo)

	input := usecase.CreateOrderInputDTO{
		Price: 100.0,
		Tax:   10.0,
	}

	output, err := uc.Execute(input)

	if err == nil {
		t.Fatal("Expected repository error, got nil")
	}

	if output != nil {
		t.Error("Expected nil output on repository error")
	}

	if !errors.Is(err, mocks.ErrRepositoryFailure) {
		t.Errorf("Expected ErrRepositoryFailure, got %v", err)
	}
}

func TestCreateOrderUseCase_NilRepository(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic with nil repository")
		}
	}()

	uc := usecase.NewCreateOrderUseCase(nil)
	input := usecase.CreateOrderInputDTO{
		Price: 100.0,
		Tax:   10.0,
	}
	uc.Execute(input)
}
