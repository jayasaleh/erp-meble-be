package models

import (
	"time"
)

// Penjualan adalah model untuk transaksi penjualan (Mode 1: POS langsung bayar)
type Penjualan struct {
	ID               uint      `gorm:"primaryKey;column:id" json:"id"`
	NomorTransaksi   string    `gorm:"uniqueIndex;not null;column:nomor_transaksi" json:"nomor_transaksi"`
	IDGudang         uint      `gorm:"index;not null;column:id_gudang" json:"id_gudang"`
	Gudang           Gudang    `gorm:"foreignKey:IDGudang" json:"gudang,omitempty"`
	NamaPelanggan    string    `gorm:"type:varchar(100);column:nama_pelanggan" json:"nama_pelanggan"`
	KontakPelanggan  string    `gorm:"type:varchar(50);column:kontak_pelanggan" json:"kontak_pelanggan"`
	Subtotal         float64   `gorm:"type:decimal(15,2);not null;column:subtotal" json:"subtotal"`                             // Total sebelum diskon
	JumlahDiskon     float64   `gorm:"type:decimal(15,2);default:0;column:jumlah_diskon" json:"jumlah_diskon"`                  // Total diskon
	Total            float64   `gorm:"type:decimal(15,2);not null;column:total" json:"total"`                                   // Total setelah diskon (harga jual)
	TotalHargaModal  float64   `gorm:"type:decimal(15,2);not null;default:0;column:total_harga_modal" json:"total_harga_modal"` // Total COGS (dari FIFO batch)
	MetodePembayaran string    `gorm:"type:varchar(20);not null;column:metode_pembayaran" json:"metode_pembayaran"`             // cash, transfer
	JumlahPembayaran float64   `gorm:"type:decimal(15,2);not null;column:jumlah_pembayaran" json:"jumlah_pembayaran"`           // Jumlah yang dibayar
	JumlahKembalian  float64   `gorm:"type:decimal(15,2);default:0;column:jumlah_kembalian" json:"jumlah_kembalian"`            // Kembalian (untuk cash)
	BuktiBayar       *string   `gorm:"type:text;column:bukti_bayar" json:"bukti_bayar,omitempty"`                               // Path foto bukti transfer (nullable)
	Status           string    `gorm:"type:varchar(20);default:'completed';column:status" json:"status"`                        // completed, voided
	CatatanInternal  string    `gorm:"type:text;column:catatan_internal" json:"catatan_internal,omitempty"`                     // Catatan opsional kasir
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

// Sales adalah alias untuk backward compatibility
type Sales = Penjualan

// ItemPenjualan adalah model untuk detail item penjualan
// Setiap item mencatat harga modal (COGS) rata-rata tertimbang dari batch FIFO yang terpakai.
type ItemPenjualan struct {
	ID           uint      `gorm:"primaryKey;column:id" json:"id"`
	IDPenjualan  uint      `gorm:"index;not null;column:id_penjualan" json:"id_penjualan"`
	Penjualan    Penjualan `gorm:"foreignKey:IDPenjualan" json:"penjualan,omitempty"`
	IDProduk     uint      `gorm:"index;not null;column:id_produk" json:"id_produk"`
	Produk       Produk    `gorm:"foreignKey:IDProduk" json:"produk,omitempty"`
	IDGudang     uint      `gorm:"index;not null;column:id_gudang" json:"id_gudang"`
	Gudang       Gudang    `gorm:"foreignKey:IDGudang" json:"gudang,omitempty"`
	Jumlah       int       `gorm:"not null;column:jumlah" json:"jumlah"`
	HargaSatuan  float64   `gorm:"type:decimal(15,2);not null;column:harga_satuan" json:"harga_satuan"`         // Harga jual per unit
	HargaModal   float64   `gorm:"type:decimal(15,2);not null;default:0;column:harga_modal" json:"harga_modal"` // COGS per unit (rata-rata tertimbang dari batch FIFO)
	PersenDiskon *float64  `gorm:"type:decimal(5,2);column:persen_diskon" json:"persen_diskon,omitempty"`       // Persentase diskon per item (%)
	JumlahDiskon float64   `gorm:"type:decimal(15,2);default:0;column:jumlah_diskon" json:"jumlah_diskon"`      // Nominal diskon per item
	Subtotal     float64   `gorm:"type:decimal(15,2);not null;column:subtotal" json:"subtotal"`                 // Total harga jual setelah diskon
	TotalModal   float64   `gorm:"type:decimal(15,2);not null;default:0;column:total_modal" json:"total_modal"` // Total COGS item ini (harga_modal * jumlah)
	DibuatPada   time.Time `gorm:"column:dibuat_pada" json:"dibuat_pada"`

	// Breakdown batch FIFO yang digunakan untuk item ini
	BatchUsage []ItemPenjualanBatch `gorm:"foreignKey:IDItemPenjualan;constraint:OnDelete:CASCADE" json:"batch_usage,omitempty"`
}

// TableName mengembalikan nama tabel untuk model ItemPenjualan
func (ItemPenjualan) TableName() string {
	return "item_penjualan"
}

// SalesItem alias
type SalesItem = ItemPenjualan

// ItemPenjualanBatch adalah model untuk mencatat breakdown batch FIFO per item penjualan.
// Karena 1 item produk bisa memakan stok dari beberapa batch (misal: ambil dari batch A 3 unit,
// batch B 2 unit), tabel ini merekam detail tersebut untuk keperluan akurasi COGS.
type ItemPenjualanBatch struct {
	ID              uint      `gorm:"primaryKey;column:id" json:"id"`
	IDItemPenjualan uint      `gorm:"index;not null;column:id_item_penjualan" json:"id_item_penjualan"`
	IDBatch         uint      `gorm:"index;not null;column:id_batch" json:"id_batch"`
	Batch           StokBatch `gorm:"foreignKey:IDBatch" json:"batch,omitempty"`
	Jumlah          int       `gorm:"not null;column:jumlah" json:"jumlah"`                              // Qty yang diambil dari batch ini
	HargaModal      float64   `gorm:"type:decimal(15,2);not null;column:harga_modal" json:"harga_modal"` // Harga modal batch ini
	TotalModal      float64   `gorm:"type:decimal(15,2);not null;column:total_modal" json:"total_modal"` // Jumlah * HargaModal
	DibuatPada      time.Time `gorm:"column:dibuat_pada" json:"dibuat_pada"`
}

// TableName mengembalikan nama tabel untuk model ItemPenjualanBatch
func (ItemPenjualanBatch) TableName() string {
	return "item_penjualan_batch"
}
