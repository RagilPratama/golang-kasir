package repository

import (
	"database/sql"
	"errors"
	"kasir-api/internal/models"
)

var (
	ErrCategoryNotFound = errors.New("category not found")
)

type CategoriesRepository interface {
	GetAll() []models.Category
	GetByID(id int) (*models.Category, error)
	Create(c models.Category) models.Category
	Update(id int, c models.Category) (*models.Category, error)
	Delete(id int) error
}

type postgresCategoryRepository struct {
	db *sql.DB
}

func NewPostgresCategoryRepository(db *sql.DB) CategoriesRepository {
	return &postgresCategoryRepository{db: db}
}

func (r *postgresCategoryRepository) GetAll() []models.Category {
	rows, err := r.db.Query("SELECT id, name, description FROM categories")
	if err != nil {
		return []models.Category{}
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var c models.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.Description); err != nil {
			continue
		}
		categories = append(categories, c)
	}
	return categories
}

func (r *postgresCategoryRepository) GetByID(id int) (*models.Category, error) {
	var c models.Category
	err := r.db.QueryRow("SELECT id, name, description FROM categories WHERE id = $1", id).Scan(&c.ID, &c.Name, &c.Description)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}
	return &c, nil
}

func (r *postgresCategoryRepository) Create(c models.Category) models.Category {
	err := r.db.QueryRow(
		"INSERT INTO categories (name, description) VALUES ($1, $2) RETURNING id",
		c.Name, c.Description,
	).Scan(&c.ID)
	if err != nil {
		return models.Category{}
	}
	return c
}

func (r *postgresCategoryRepository) Update(id int, c models.Category) (*models.Category, error) {
	result, err := r.db.Exec(
		"UPDATE categories SET name = $1, description = $2 WHERE id = $3",
		c.Name, c.Description, id,
	)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, ErrCategoryNotFound
	}

	c.ID = id
	return &c, nil
}

func (r *postgresCategoryRepository) Delete(id int) error {
	result, err := r.db.Exec("DELETE FROM categories WHERE id = $1", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrCategoryNotFound
	}

	return nil
}

type MemoryCategoryRepository struct {
	categories []models.Category
}

func NewMemoryCategoryRepository() CategoriesRepository {
	return &MemoryCategoryRepository{
		categories: []models.Category{
			{
				ID:          1,
				Name:        "Makanan",
				Description: "Kategori makanan",
			},
			{
				ID:          2,
				Name:        "Minuman",
				Description: "Kategori minuman",
			},
		},
	}
}

func (r *MemoryCategoryRepository) GetAll() []models.Category {
	return r.categories
}

func (r *MemoryCategoryRepository) GetByID(id int) (*models.Category, error) {
	for i := range r.categories {
		if r.categories[i].ID == id {
			return &r.categories[i], nil
		}
	}
	return nil, ErrCategoryNotFound
}

func (r *MemoryCategoryRepository) Create(c models.Category) models.Category {
	c.ID = len(r.categories) + 1
	r.categories = append(r.categories, c)
	return c
}

func (r *MemoryCategoryRepository) Update(id int, updateData models.Category) (*models.Category, error) {
	for i := range r.categories {
		if r.categories[i].ID == id {
			r.categories[i].Name = updateData.Name
			r.categories[i].Description = updateData.Description
			return &r.categories[i], nil
		}
	}
	return nil, ErrCategoryNotFound
}

func (r *MemoryCategoryRepository) Delete(id int) error {
	for i := range r.categories {
		if r.categories[i].ID == id {
			r.categories = append(r.categories[:i], r.categories[i+1:]...)
			return nil
		}
	}
	return ErrCategoryNotFound
}
