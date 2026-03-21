package models

import (
	"time"
)

// StokBatch adalah model untuk tracking batch stok (FIFO)
// Setiap batch merepresentasikan kelompok item yang masuk bersamaan
type StokBatch struct {
	ID          uint      `gorm:"primaryKey;column:id" json:"id"`
	IDProduk    uint      `gorm:"index:idx_batch_product_warehouse;not null;column:id_produk" json:"id_produk"`
	Produk      Produk    `gorm:"foreignKey:IDProduk" json:"produk,omitempty"`
	IDGudang    uint      `gorm:"index:idx_batch_product_warehouse;not null;column:id_gudang" json:"id_gudang"`
	Gudang      Gudang    `gorm:"foreignKey:IDGudang" json:"gudang,omitempty"`
	
	// Field kunci untuk FIFO
	TanggalMasuk   time.Time  `gorm:"index;not null;column:tanggal_masuk" json:"tanggal_masuk"` // Kunci sorting FIFO
	TanggalKadaluarsa *time.Time `gorm:"column:tanggal_kadaluarsa" json:"tanggal_kadaluarsa,omitempty"`
	
	// Manajemen kuantitas
	JumlahAwal     int       `gorm:"not null;column:jumlah_awal" json:"jumlah_awal"`         // Qty saat batch dibuat
	JumlahSaatIni  int       `gorm:"not null;default:0;column:jumlah_saat_ini" json:"jumlah_saat_ini"` // Qty tersisa
	
	// Costing
	HargaModal     float64   `gorm:"type:decimal(15,2);not null;column:harga_modal" json:"harga_modal"` // HPP spesifik batch ini
	
	// Referensi ke transaksi sumber
	IDReferensi    *uint     `gorm:"index;column:id_referensi" json:"id_referensi"` // ID BarangMasuk
	TipeReferensi  string    `gorm:"type:varchar(50);column:tipe_referensi" json:"tipe_referensi"` // "stock_in", "adjustment", dll
	
	// Metadata
	Aktif          bool      `gorm:"default:true;column:aktif" json:"aktif"` // False jika JumlahSaatIni == 0
	Keterangan     string    `gorm:"type:text;column:keterangan" json:"keterangan"`
	DibuatPada     time.Time `gorm:"column:dibuat_pada" json:"dibuat_pada"`
	DiperbaruiPada time.Time `gorm:"column:diperbarui_pada" json:"diperbarui_pada"`
}

// TableName mengembalikan nama tabel untuk model StokBatch
func (StokBatch) TableName() string {
	return "stok_batch"
}

// StockBatch adalah alias untuk backward compatibility
type StockBatch = StokBatch
