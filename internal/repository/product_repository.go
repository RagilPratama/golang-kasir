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
	GetAll(nameFilter string) []models.Product
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

func (r *postgresProductRepository) GetAll(nameFilter string) []models.Product {
	query := `
		SELECT p.id, p.name, p.price, p.stock, p.category_id, 
		       c.id, c.name, c.description
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
	`

	args := []interface{}{}
	if nameFilter != "" {
		query += " WHERE p.name ILIKE $1"
		args = append(args, "%"+nameFilter+"%")
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return []models.Product{}
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		var categoryID sql.NullInt64
		var cID sql.NullInt64
		var cName sql.NullString
		var cDesc sql.NullString

		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &categoryID,
			&cID, &cName, &cDesc); err != nil {
			continue
		}

		if categoryID.Valid {
			p.CategoryID = int(categoryID.Int64)
			if cID.Valid {
				p.Category = &models.Category{
					ID:          int(cID.Int64),
					Name:        cName.String,
					Description: cDesc.String,
				}
			}
		}
		products = append(products, p)
	}
	return products
}

func (r *postgresProductRepository) GetByID(id int) (*models.Product, error) {
	var p models.Product
	var categoryID sql.NullInt64
	var cID sql.NullInt64
	var cName sql.NullString
	var cDesc sql.NullString

	query := `
		SELECT p.id, p.name, p.price, p.stock, p.category_id, 
		       c.id, c.name, c.description
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
		WHERE p.id = $1
	`
	err := r.db.QueryRow(query, id).Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &categoryID,
		&cID, &cName, &cDesc)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrProductNotFound
		}
		return nil, err
	}

	if categoryID.Valid {
		p.CategoryID = int(categoryID.Int64)
		if cID.Valid {
			p.Category = &models.Category{
				ID:          int(cID.Int64),
				Name:        cName.String,
				Description: cDesc.String,
			}
		}
	}
	return &p, nil
}

func (r *postgresProductRepository) Create(p models.Product) models.Product {
	err := r.db.QueryRow(
		"INSERT INTO products (name, price, stock, category_id) VALUES ($1, $2, $3, $4) RETURNING id",
		p.Name, p.Price, p.Stock, p.CategoryID,
	).Scan(&p.ID)
	if err != nil {
		return models.Product{}
	}
	return p
}

func (r *postgresProductRepository) Update(id int, p models.Product) (*models.Product, error) {
	result, err := r.db.Exec(
		"UPDATE products SET name = $1, price = $2, stock = $3, category_id = $4 WHERE id = $5",
		p.Name, p.Price, p.Stock, p.CategoryID, id,
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
