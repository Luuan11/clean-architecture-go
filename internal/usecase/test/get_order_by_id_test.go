package usecase_test

import (
	"testing"

	"github.com/luuan11/clean-architecture/internal/entity"
	"github.com/luuan11/clean-architecture/internal/usecase"
	"github.com/luuan11/clean-architecture/internal/usecase/test/mocks"
)

func TestGetOrderByIDUseCase_Success(t *testing.T) {
	mockRepo := mocks.NewOrderRepositoryMock()
	
	order, _ := entity.NewOrder(100.0, 10.0)
	mockRepo.Save(order)
	
	uc := usecase.NewGetOrderByIDUseCase(mockRepo)
	
	input := usecase.GetOrderByIDInputDTO{
		ID: order.ID,
	}
	
	output, err := uc.Execute(input)
	
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	if output.ID != order.ID {
		t.Errorf("Expected ID %s, got %s", order.ID, output.ID)
	}
	
	if output.Price != 100.0 {
		t.Errorf("Expected price 100.0, got %.2f", output.Price)
	}
	
	if output.Tax != 10.0 {
		t.Errorf("Expected tax 10.0, got %.2f", output.Tax)
	}
	
	if output.FinalPrice != 110.0 {
		t.Errorf("Expected final price 110.0, got %.2f", output.FinalPrice)
	}
}

func TestGetOrderByIDUseCase_OrderNotFound(t *testing.T) {
	mockRepo := mocks.NewOrderRepositoryMock()
	uc := usecase.NewGetOrderByIDUseCase(mockRepo)
	
	input := usecase.GetOrderByIDInputDTO{
		ID: "non-existent-id",
	}
	
	_, err := uc.Execute(input)
	
	if err != entity.ErrOrderNotFound {
		t.Errorf("Expected ErrOrderNotFound, got %v", err)
	}
}

func TestGetOrderByIDUseCase_RepositoryError(t *testing.T) {
	mockRepo := mocks.NewOrderRepositoryMock()
	mockRepo.SetFindError(mocks.ErrDatabaseConnection)
	
	uc := usecase.NewGetOrderByIDUseCase(mockRepo)
	
	input := usecase.GetOrderByIDInputDTO{
		ID: "some-id",
	}
	
	_, err := uc.Execute(input)
	
	if err != mocks.ErrDatabaseConnection {
		t.Errorf("Expected ErrDatabaseConnection, got %v", err)
	}
}

func TestGetOrderByIDUseCase_EmptyID(t *testing.T) {
	mockRepo := mocks.NewOrderRepositoryMock()
	uc := usecase.NewGetOrderByIDUseCase(mockRepo)
	
	input := usecase.GetOrderByIDInputDTO{
		ID: "",
	}
	
	_, err := uc.Execute(input)
	
	if err != entity.ErrOrderNotFound {
		t.Errorf("Expected ErrOrderNotFound, got %v", err)
	}
}

func TestGetOrderByIDUseCase_MultipleOrders(t *testing.T) {
	mockRepo := mocks.NewOrderRepositoryMock()
	
	order1, _ := entity.NewOrder(100.0, 10.0)
	order2, _ := entity.NewOrder(200.0, 20.0)
	order3, _ := entity.NewOrder(300.0, 30.0)
	
	mockRepo.Save(order1)
	mockRepo.Save(order2)
	mockRepo.Save(order3)
	
	uc := usecase.NewGetOrderByIDUseCase(mockRepo)
	
	input := usecase.GetOrderByIDInputDTO{
		ID: order2.ID,
	}
	
	output, err := uc.Execute(input)
	
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	if output.ID != order2.ID {
		t.Errorf("Expected order2, got different order")
	}
	
	if output.Price != 200.0 {
		t.Errorf("Expected price 200.0, got %.2f", output.Price)
	}
}

func TestGetOrderByIDUseCase_NilRepository(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic with nil repository")
		}
	}()
	
	uc := usecase.NewGetOrderByIDUseCase(nil)
	
	input := usecase.GetOrderByIDInputDTO{
		ID: "some-id",
	}
	
	uc.Execute(input)
}
