package repository

import (
	"database/sql"
	"errors"
	"kasir-api/internal/models"
)

var (
	ErrProductNotFound = errors.New("product not found")
)

type ProductRepository interface {
	GetAll() []models.Product
	GetByID(id int) (*models.Product, error)
	Create(p models.Product) models.Product
	Update(id int, p models.Product) (*models.Product, error)
	Delete(id int) error
}

type postgresProductRepository struct {
	db *sql.DB
}

func NewPostgresProductRepository(db *sql.DB) ProductRepository {
	return &postgresProductRepository{db: db}
}

func (r *postgresProductRepository) GetAll() []models.Product {
	rows, err := r.db.Query("SELECT id, name, price, stock FROM products")
	if err != nil {
		return []models.Product{}
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock); err != nil {
			continue
		}
		products = append(products, p)
	}
	return products
}

func (r *postgresProductRepository) GetByID(id int) (*models.Product, error) {
	var p models.Product
	err := r.db.QueryRow("SELECT id, name, price, stock FROM products WHERE id = $1", id).Scan(&p.ID, &p.Name, &p.Price, &p.Stock)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrProductNotFound
		}
		return nil, err
	}
	return &p, nil
}

func (r *postgresProductRepository) Create(p models.Product) models.Product {
	err := r.db.QueryRow(
		"INSERT INTO products (name, price, stock) VALUES ($1, $2, $3) RETURNING id",
		p.Name, p.Price, p.Stock,
	).Scan(&p.ID)
	if err != nil {
		return models.Product{}
	}
	return p
}

func (r *postgresProductRepository) Update(id int, p models.Product) (*models.Product, error) {
	result, err := r.db.Exec(
		"UPDATE products SET name = $1, price = $2, stock = $3 WHERE id = $4",
		p.Name, p.Price, p.Stock, id,
	)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, ErrProductNotFound
	}

	p.ID = id
	return &p, nil
}

func (r *postgresProductRepository) Delete(id int) error {
	result, err := r.db.Exec("DELETE FROM products WHERE id = $1", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrProductNotFound
	}
	return nil
}
