package webserver

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/luuan11/clean-architecture/internal/entity"
	"github.com/luuan11/clean-architecture/internal/usecase"
)

type WebOrderHandler struct {
	CreateOrderUseCase  *usecase.CreateOrderUseCase
	ListOrdersUseCase   *usecase.ListOrdersUseCase
	GetOrderByIDUseCase *usecase.GetOrderByIDUseCase
	UpdateOrderUseCase  *usecase.UpdateOrderUseCase
	DeleteOrderUseCase  *usecase.DeleteOrderUseCase
}

func NewWebOrderHandler(
	createOrderUseCase *usecase.CreateOrderUseCase,
	listOrdersUseCase *usecase.ListOrdersUseCase,
	getOrderByIDUseCase *usecase.GetOrderByIDUseCase,
	updateOrderUseCase *usecase.UpdateOrderUseCase,
	deleteOrderUseCase *usecase.DeleteOrderUseCase,
) *WebOrderHandler {
	return &WebOrderHandler{
		CreateOrderUseCase:  createOrderUseCase,
		ListOrdersUseCase:   listOrdersUseCase,
		GetOrderByIDUseCase: getOrderByIDUseCase,
		UpdateOrderUseCase:  updateOrderUseCase,
		DeleteOrderUseCase:  deleteOrderUseCase,
	}
}

func (h *WebOrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var input usecase.CreateOrderInputDTO
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	output, err := h.CreateOrderUseCase.Execute(input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(output)
}

func (h *WebOrderHandler) ListOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := h.ListOrdersUseCase.Execute()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(orders)
}

func (h *WebOrderHandler) GetOrderByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	input := usecase.GetOrderByIDInputDTO{
		ID: id,
	}

	output, err := h.GetOrderByIDUseCase.Execute(input)
	if err != nil {
		if err == entity.ErrOrderNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(output)
}

func (h *WebOrderHandler) UpdateOrder(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var input usecase.UpdateOrderInputDTO
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	input.ID = id

	output, err := h.UpdateOrderUseCase.Execute(input)
	if err != nil {
		if err == entity.ErrOrderNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		if err == entity.ErrInvalidPrice || err == entity.ErrInvalidTax {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(output)
}

func (h *WebOrderHandler) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	input := usecase.DeleteOrderInputDTO{
		ID: id,
	}

	err := h.DeleteOrderUseCase.Execute(input)
	if err != nil {
		if err == entity.ErrOrderNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
