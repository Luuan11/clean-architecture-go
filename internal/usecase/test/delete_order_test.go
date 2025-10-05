package usecase_test

import (
	"testing"

	"github.com/luuan11/clean-architecture/internal/entity"
	"github.com/luuan11/clean-architecture/internal/usecase"
	"github.com/luuan11/clean-architecture/internal/usecase/test/mocks"
)

func TestDeleteOrderUseCase_Success(t *testing.T) {
	mockRepo := mocks.NewOrderRepositoryMock()

	order, _ := entity.NewOrder(100.0, 10.0)
	mockRepo.Save(order)

	uc := usecase.NewDeleteOrderUseCase(mockRepo)

	input := usecase.DeleteOrderInputDTO{
		ID: order.ID,
	}

	err := uc.Execute(input)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	_, err = mockRepo.FindByID(order.ID)
	if err != entity.ErrOrderNotFound {
		t.Error("Expected order to be deleted")
	}
}

func TestDeleteOrderUseCase_OrderNotFound(t *testing.T) {
	mockRepo := mocks.NewOrderRepositoryMock()
	uc := usecase.NewDeleteOrderUseCase(mockRepo)

	input := usecase.DeleteOrderInputDTO{
		ID: "non-existent-id",
	}

	err := uc.Execute(input)

	if err != entity.ErrOrderNotFound {
		t.Errorf("Expected ErrOrderNotFound, got %v", err)
	}
}

func TestDeleteOrderUseCase_VerifyDeletion(t *testing.T) {
	mockRepo := mocks.NewOrderRepositoryMock()

	order1, _ := entity.NewOrder(100.0, 10.0)
	order2, _ := entity.NewOrder(200.0, 20.0)
	order3, _ := entity.NewOrder(300.0, 30.0)

	mockRepo.Save(order1)
	mockRepo.Save(order2)
	mockRepo.Save(order3)

	uc := usecase.NewDeleteOrderUseCase(mockRepo)

	input := usecase.DeleteOrderInputDTO{
		ID: order2.ID,
	}

	err := uc.Execute(input)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	orders, _ := mockRepo.FindAll()

	if len(orders) != 2 {
		t.Errorf("Expected 2 orders remaining, got %d", len(orders))
	}

	for _, order := range orders {
		if order.ID == order2.ID {
			t.Error("Expected order2 to be deleted")
		}
	}
}

func TestDeleteOrderUseCase_DeleteAllOrders(t *testing.T) {
	mockRepo := mocks.NewOrderRepositoryMock()

	order1, _ := entity.NewOrder(100.0, 10.0)
	order2, _ := entity.NewOrder(200.0, 20.0)

	mockRepo.Save(order1)
	mockRepo.Save(order2)

	uc := usecase.NewDeleteOrderUseCase(mockRepo)

	uc.Execute(usecase.DeleteOrderInputDTO{ID: order1.ID})
	uc.Execute(usecase.DeleteOrderInputDTO{ID: order2.ID})

	orders, _ := mockRepo.FindAll()

	if len(orders) != 0 {
		t.Errorf("Expected 0 orders, got %d", len(orders))
	}
}

func TestDeleteOrderUseCase_DeleteSameOrderTwice(t *testing.T) {
	mockRepo := mocks.NewOrderRepositoryMock()

	order, _ := entity.NewOrder(100.0, 10.0)
	mockRepo.Save(order)

	uc := usecase.NewDeleteOrderUseCase(mockRepo)

	input := usecase.DeleteOrderInputDTO{
		ID: order.ID,
	}

	err := uc.Execute(input)
	if err != nil {
		t.Fatalf("First delete failed: %v", err)
	}

	err = uc.Execute(input)
	if err != entity.ErrOrderNotFound {
		t.Errorf("Expected ErrOrderNotFound on second delete, got %v", err)
	}
}

func TestDeleteOrderUseCase_RepositoryFindError(t *testing.T) {
	mockRepo := mocks.NewOrderRepositoryMock()
	mockRepo.SetFindError(mocks.ErrDatabaseConnection)

	uc := usecase.NewDeleteOrderUseCase(mockRepo)

	input := usecase.DeleteOrderInputDTO{
		ID: "some-id",
	}

	err := uc.Execute(input)

	if err != mocks.ErrDatabaseConnection {
		t.Errorf("Expected ErrDatabaseConnection, got %v", err)
	}
}

func TestDeleteOrderUseCase_RepositoryDeleteError(t *testing.T) {
	mockRepo := mocks.NewOrderRepositoryMock()

	order, _ := entity.NewOrder(100.0, 10.0)
	mockRepo.Save(order)

	mockRepo.SetDeleteError(mocks.ErrRepositoryFailure)

	uc := usecase.NewDeleteOrderUseCase(mockRepo)

	input := usecase.DeleteOrderInputDTO{
		ID: order.ID,
	}

	err := uc.Execute(input)

	if err != mocks.ErrRepositoryFailure {
		t.Errorf("Expected ErrRepositoryFailure, got %v", err)
	}
}

func TestDeleteOrderUseCase_EmptyID(t *testing.T) {
	mockRepo := mocks.NewOrderRepositoryMock()
	uc := usecase.NewDeleteOrderUseCase(mockRepo)

	input := usecase.DeleteOrderInputDTO{
		ID: "",
	}

	err := uc.Execute(input)

	if err != entity.ErrOrderNotFound {
		t.Errorf("Expected ErrOrderNotFound, got %v", err)
	}
}

func TestDeleteOrderUseCase_NilRepository(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic with nil repository")
		}
	}()

	uc := usecase.NewDeleteOrderUseCase(nil)

	input := usecase.DeleteOrderInputDTO{
		ID: "some-id",
	}

	uc.Execute(input)
}
