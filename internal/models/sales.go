package models

import (
	"time"
)

// Penjualan adalah model untuk transaksi penjualan (POS)
type Penjualan struct {
	ID               uint      `gorm:"primaryKey;column:id" json:"id"`
	NomorTransaksi   string    `gorm:"uniqueIndex;not null;column:nomor_transaksi" json:"nomor_transaksi"`
	NamaPelanggan    string    `gorm:"type:varchar(100);column:nama_pelanggan" json:"nama_pelanggan"`
	KontakPelanggan  string    `gorm:"type:varchar(50);column:kontak_pelanggan" json:"kontak_pelanggan"`
	Subtotal         float64   `gorm:"type:decimal(15,2);not null;column:subtotal" json:"subtotal"`                   // Subtotal sebelum diskon
	JumlahDiskon     float64   `gorm:"type:decimal(15,2);default:0;column:jumlah_diskon" json:"jumlah_diskon"`        // Jumlah diskon
	Total            float64   `gorm:"type:decimal(15,2);not null;column:total" json:"total"`                         // Total setelah diskon
	MetodePembayaran string    `gorm:"type:varchar(20);not null;column:metode_pembayaran" json:"metode_pembayaran"`   // cash, transfer, qris
	JumlahPembayaran float64   `gorm:"type:decimal(15,2);not null;column:jumlah_pembayaran" json:"jumlah_pembayaran"` // Jumlah yang dibayar
	JumlahKembalian  float64   `gorm:"type:decimal(15,2);default:0;column:jumlah_kembalian" json:"jumlah_kembalian"`  // Kembalian
	IDKasir          uint      `gorm:"index;not null;column:id_kasir" json:"id_kasir"`
	Kasir            Pengguna  `gorm:"foreignKey:IDKasir" json:"kasir,omitempty"`
	DibuatPada       time.Time `gorm:"column:dibuat_pada" json:"dibuat_pada"`
	DiperbaruiPada   time.Time `gorm:"column:diperbarui_pada" json:"diperbarui_pada"`

	// Relationship
	Items []ItemPenjualan `gorm:"foreignKey:IDPenjualan;constraint:OnDelete:CASCADE" json:"items,omitempty"`
}

// TableName mengembalikan nama tabel untuk model Penjualan
func (Penjualan) TableName() string {
	return "penjualan"
}

// Sales adalah alias untuk backward compatibility (akan dihapus nanti)
type Sales = Penjualan

// ItemPenjualan adalah model untuk detail item penjualan
type ItemPenjualan struct {
	ID           uint      `gorm:"primaryKey;column:id" json:"id"`
	IDPenjualan  uint      `gorm:"index;not null;column:id_penjualan" json:"id_penjualan"`
	Penjualan    Penjualan `gorm:"foreignKey:IDPenjualan" json:"penjualan,omitempty"`
	IDProduk     uint      `gorm:"index;not null;column:id_produk" json:"id_produk"`
	Produk       Produk    `gorm:"foreignKey:IDProduk" json:"produk,omitempty"`
	Jumlah       int       `gorm:"not null;column:jumlah" json:"jumlah"`
	HargaSatuan  float64   `gorm:"type:decimal(15,2);not null;column:harga_satuan" json:"harga_satuan"`
	PersenDiskon *float64  `gorm:"type:decimal(5,2);column:persen_diskon" json:"persen_diskon"`  // Persentase diskon
	JumlahDiskon *float64  `gorm:"type:decimal(15,2);column:jumlah_diskon" json:"jumlah_diskon"` // Jumlah diskon
	Subtotal     float64   `gorm:"type:decimal(15,2);not null;column:subtotal" json:"subtotal"`  // Total setelah diskon
	DibuatPada   time.Time `gorm:"column:dibuat_pada" json:"dibuat_pada"`
}

// TableName mengembalikan nama tabel untuk model ItemPenjualan
func (ItemPenjualan) TableName() string {
	return "item_penjualan"
}

// SalesItem adalah alias untuk backward compatibility (akan dihapus nanti)
type SalesItem = ItemPenjualan
