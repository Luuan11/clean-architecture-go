package usecase_test

import (
	"errors"
	"testing"

	"github.com/luuan11/clean-architecture/internal/entity"
	"github.com/luuan11/clean-architecture/internal/usecase"
	"github.com/luuan11/clean-architecture/internal/usecase/test/mocks"
)

func TestListOrdersUseCase_Success(t *testing.T) {
	mockRepo := mocks.NewOrderRepositoryMock()

	order1, _ := entity.NewOrder(100.0, 10.0)
	order2, _ := entity.NewOrder(200.0, 20.0)
	mockRepo.Save(order1)
	mockRepo.Save(order2)

	uc := usecase.NewListOrdersUseCase(mockRepo)

	orders, err := uc.Execute()

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if orders == nil {
		t.Fatal("Expected orders list, got nil")
	}

	if len(orders) != 2 {
		t.Errorf("Expected 2 orders, got %d", len(orders))
	}

	if orders[0].Price != 100.0 {
		t.Errorf("Expected first order price 100.0, got %f", orders[0].Price)
	}

	if orders[1].Price != 200.0 {
		t.Errorf("Expected second order price 200.0, got %f", orders[1].Price)
	}
}

func TestListOrdersUseCase_EmptyList(t *testing.T) {
	mockRepo := mocks.NewOrderRepositoryMock()
	uc := usecase.NewListOrdersUseCase(mockRepo)

	orders, err := uc.Execute()

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if orders != nil && len(orders) != 0 {
		t.Errorf("Expected empty list, got %d orders", len(orders))
	}
}

func TestListOrdersUseCase_MultipleOrders(t *testing.T) {
	mockRepo := mocks.NewOrderRepositoryMock()

	expectedCount := 10
	for i := 1; i <= expectedCount; i++ {
		order, err := entity.NewOrder(float64(i*100), float64(i*10))
		if err != nil {
			t.Fatalf("Failed to create order %d: %v", i, err)
		}
		mockRepo.Save(order)
	}

	uc := usecase.NewListOrdersUseCase(mockRepo)

	orders, err := uc.Execute()

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(orders) != expectedCount {
		t.Errorf("Expected %d orders, got %d", expectedCount, len(orders))
	}

	for i, order := range orders {
		expectedPrice := float64((i + 1) * 100)
		if order.Price != expectedPrice {
			t.Errorf("Order %d: expected price %.2f, got %.2f", i, expectedPrice, order.Price)
		}
	}
}

func TestListOrdersUseCase_VerifyOutputDTO(t *testing.T) {
	mockRepo := mocks.NewOrderRepositoryMock()

	originalOrder, _ := entity.NewOrder(150.75, 15.25)
	mockRepo.Save(originalOrder)

	uc := usecase.NewListOrdersUseCase(mockRepo)

	orders, err := uc.Execute()

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(orders) != 1 {
		t.Fatalf("Expected 1 order, got %d", len(orders))
	}

	output := orders[0]

	if output.ID != originalOrder.ID {
		t.Errorf("Expected ID %s, got %s", originalOrder.ID, output.ID)
	}

	if output.Price != originalOrder.Price {
		t.Errorf("Expected price %.2f, got %.2f", originalOrder.Price, output.Price)
	}

	if output.Tax != originalOrder.Tax {
		t.Errorf("Expected tax %.2f, got %.2f", originalOrder.Tax, output.Tax)
	}

	if output.FinalPrice != originalOrder.FinalPrice {
		t.Errorf("Expected final price %.2f, got %.2f", originalOrder.FinalPrice, output.FinalPrice)
	}
}

func TestListOrdersUseCase_RepositoryError(t *testing.T) {
	mockRepo := mocks.NewOrderRepositoryMock()
	mockRepo.SetFindError(mocks.ErrRepositoryFailure)

	uc := usecase.NewListOrdersUseCase(mockRepo)

	orders, err := uc.Execute()

	if err == nil {
		t.Fatal("Expected repository error, got nil")
	}

	if orders != nil {
		t.Error("Expected nil orders on error")
	}

	if !errors.Is(err, mocks.ErrRepositoryFailure) {
		t.Errorf("Expected ErrRepositoryFailure, got %v", err)
	}
}

func TestListOrdersUseCase_DatabaseConnectionError(t *testing.T) {
	mockRepo := mocks.NewOrderRepositoryMock()
	mockRepo.SetFindError(mocks.ErrDatabaseConnection)

	uc := usecase.NewListOrdersUseCase(mockRepo)

	orders, err := uc.Execute()

	if err == nil {
		t.Fatal("Expected database error, got nil")
	}

	if orders != nil {
		t.Error("Expected nil orders on error")
	}

	if !errors.Is(err, mocks.ErrDatabaseConnection) {
		t.Errorf("Expected ErrDatabaseConnection, got %v", err)
	}
}

func TestListOrdersUseCase_OrdersWithDifferentValues(t *testing.T) {
	mockRepo := mocks.NewOrderRepositoryMock()

	testCases := []struct {
		price float64
		tax   float64
	}{
		{10.50, 1.05},
		{999.99, 99.99},
		{0.01, 0.001},
		{5000.00, 500.00},
	}

	for _, tc := range testCases {
		order, _ := entity.NewOrder(tc.price, tc.tax)
		mockRepo.Save(order)
	}

	uc := usecase.NewListOrdersUseCase(mockRepo)

	orders, err := uc.Execute()

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(orders) != len(testCases) {
		t.Errorf("Expected %d orders, got %d", len(testCases), len(orders))
	}

	for i, order := range orders {
		expectedFinalPrice := testCases[i].price + testCases[i].tax
		if order.FinalPrice != expectedFinalPrice {
			t.Errorf("Order %d: expected final price %.3f, got %.3f",
				i, expectedFinalPrice, order.FinalPrice)
		}
	}
}

func TestListOrdersUseCase_NilRepository(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic with nil repository")
		}
	}()

	uc := usecase.NewListOrdersUseCase(nil)
	uc.Execute()
}
