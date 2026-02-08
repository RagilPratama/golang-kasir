package models

type ProductBestSeller struct {
	Name       string `json:"nama"`
	QtyTerjual int    `json:"qty_terjual"`
}

type SalesReport struct {
	TotalRevenue   int               `json:"total_revenue"`
	TotalTransaksi int               `json:"total_transaksi"`
	ProdukTerlaris ProductBestSeller `json:"produk_terlaris"`
}
