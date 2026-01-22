package repository

import (
	"errors"
	"kasir-api/internal/models"
)

var (
	ErrProdukNotFound = errors.New("produk not found")
)

type ProdukRepository interface {
	GetAll() []models.Produk
	GetByID(id int) (*models.Produk, error)
	Create(p models.Produk) models.Produk
	Update(id int, p models.Produk) (*models.Produk, error)
}

type memoryProdukRepository struct {
	produk []models.Produk
}

func NewMemoryProdukRepository() ProdukRepository {
	return &memoryProdukRepository{
		produk: []models.Produk{
			{ID: 1, Nama: "Indomie Bangladesh", Harga: 7500, Stok: 20},
			{ID: 2, Nama: "Teh Tarik", Harga: 3000, Stok: 30},
		},
	}
}

func (r *memoryProdukRepository) GetAll() []models.Produk {
	return r.produk
}

func (r *memoryProdukRepository) GetByID(id int) (*models.Produk, error) {
	for i := range r.produk {
		if r.produk[i].ID == id {
			return &r.produk[i], nil
		}
	}
	return nil, ErrProdukNotFound
}

func (r *memoryProdukRepository) Create(p models.Produk) models.Produk {
	p.ID = len(r.produk) + 1
	r.produk = append(r.produk, p)
	return p
}

func (r *memoryProdukRepository) Update(id int, updateData models.Produk) (*models.Produk, error) {
	for i := range r.produk {
		if r.produk[i].ID == id {
			r.produk[i].Nama = updateData.Nama
			r.produk[i].Harga = updateData.Harga
			r.produk[i].Stok = updateData.Stok
			return &r.produk[i], nil
		}
	}
	return nil, ErrProdukNotFound
}
