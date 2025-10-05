package database

import (
	"database/sql"

	"github.com/luuan11/clean-architecture/internal/entity"
)

type OrderRepository struct {
	Db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{Db: db}
}

func (r *OrderRepository) Save(order *entity.Order) error {
	stmt, err := r.Db.Prepare("INSERT INTO orders (id, price, tax, final_price, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(order.ID, order.Price, order.Tax, order.FinalPrice, order.CreatedAt, order.UpdatedAt)
	return err
}

func (r *OrderRepository) FindAll() ([]*entity.Order, error) {
	rows, err := r.Db.Query("SELECT id, price, tax, final_price, created_at, updated_at FROM orders ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*entity.Order
	for rows.Next() {
		var order entity.Order
		err := rows.Scan(&order.ID, &order.Price, &order.Tax, &order.FinalPrice, &order.CreatedAt, &order.UpdatedAt)
		if err != nil {
			return nil, err
		}
		orders = append(orders, &order)
	}

	return orders, nil
}

func (r *OrderRepository) FindByID(id string) (*entity.Order, error) {
	var order entity.Order
	err := r.Db.QueryRow("SELECT id, price, tax, final_price, created_at, updated_at FROM orders WHERE id = ?", id).
		Scan(&order.ID, &order.Price, &order.Tax, &order.FinalPrice, &order.CreatedAt, &order.UpdatedAt)
	
	if err == sql.ErrNoRows {
		return nil, entity.ErrOrderNotFound
	}
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (r *OrderRepository) Update(order *entity.Order) error {
	stmt, err := r.Db.Prepare("UPDATE orders SET price = ?, tax = ?, final_price = ?, updated_at = ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(order.Price, order.Tax, order.FinalPrice, order.UpdatedAt, order.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return entity.ErrOrderNotFound
	}

	return nil
}

func (r *OrderRepository) Delete(id string) error {
	stmt, err := r.Db.Prepare("DELETE FROM orders WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return entity.ErrOrderNotFound
	}

	return nil
}
