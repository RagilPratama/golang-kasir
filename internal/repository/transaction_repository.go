package repository

import (
	"database/sql"
	"kasir-api/internal/models"
	"time"
)

type TransactionRepository interface {
	CreateTransaction(transaction *models.Transaction) error
	GetSalesSummary(startDate, endDate time.Time) (models.SalesReport, error)
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

func (r *postgresTransactionRepository) GetSalesSummary(startDate, endDate time.Time) (models.SalesReport, error) {
	var report models.SalesReport

	// 1. Get Total Revenue & Total Transaction
	querySummary := `
		SELECT COALESCE(SUM(total_amount), 0), COUNT(id)
		FROM transactions
		WHERE created_at BETWEEN $1 AND $2
	`
	err := r.db.QueryRow(querySummary, startDate, endDate).Scan(&report.TotalRevenue, &report.TotalTransaksi)
	if err != nil {
		return report, err
	}

	// 2. Get Best Selling Product
	queryBestSeller := `
		SELECT p.name, COALESCE(SUM(td.quantity), 0) as qty_terjual
		FROM transaction_details td
		JOIN products p ON td.product_id = p.id
		JOIN transactions t ON td.transaction_id = t.id
		WHERE t.created_at BETWEEN $1 AND $2
		GROUP BY p.name
		ORDER BY qty_terjual DESC
		LIMIT 1
	`
	err = r.db.QueryRow(queryBestSeller, startDate, endDate).Scan(&report.ProdukTerlaris.Name, &report.ProdukTerlaris.QtyTerjual)
	if err != nil {
		if err == sql.ErrNoRows {
			// No sales yet
			report.ProdukTerlaris = models.ProductBestSeller{Name: "-", QtyTerjual: 0}
		} else {
			return report, err
		}
	}

	return report, nil
}
