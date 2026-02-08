package repository

import (
	"database/sql"
	"kasir-api/internal/models"
)

type TransactionRepository interface {
	CreateTransaction(transaction *models.Transaction) error
}

type postgresTransactionRepository struct {
	db *sql.DB
}

func NewPostgresTransactionRepository(db *sql.DB) TransactionRepository {
	return &postgresTransactionRepository{db: db}
}

func (r *postgresTransactionRepository) CreateTransaction(t *models.Transaction) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert Transaction
	query := `INSERT INTO transactions (total_amount, created_at) VALUES ($1, NOW()) RETURNING id, created_at`
	err = tx.QueryRow(query, t.TotalAmount).Scan(&t.ID, &t.CreatedAt)
	if err != nil {
		return err
	}

	// Insert Details
	detailQuery := `INSERT INTO transaction_details (transaction_id, product_id, quantity, subtotal) VALUES ($1, $2, $3, $4)`
	stmt, err := tx.Prepare(detailQuery)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, detail := range t.Details {
		_, err = stmt.Exec(t.ID, detail.ProductID, detail.Quantity, detail.Subtotal)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
