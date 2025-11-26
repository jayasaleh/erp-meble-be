package models

import (
	"time"
)

// PesananPembelian adalah model untuk Purchase Order
type PesananPembelian struct {
	ID                    uint       `gorm:"primaryKey;column:id" json:"id"`
	NomorPO               string     `gorm:"uniqueIndex;not null;column:nomor_po" json:"nomor_po"`
	IDPemasok             uint       `gorm:"index;not null;column:id_supplier" json:"id_pemasok"`
	Pemasok               Pemasok    `gorm:"foreignKey:IDPemasok" json:"pemasok,omitempty"`
	TanggalPesan          time.Time  `gorm:"not null;column:tanggal_pesan" json:"tanggal_pesan"`
	TanggalJatuhTempo     *time.Time `gorm:"column:tanggal_jatuh_tempo" json:"tanggal_jatuh_tempo"`
	Status                string     `gorm:"type:varchar(30);default:'draft';column:status" json:"status"` // draft, sent, approved, partially_received, completed, cancelled
	Total                 float64    `gorm:"type:decimal(15,2);not null;column:total" json:"total"`
	DibuatOleh            uint       `gorm:"index;not null;column:dibuat_oleh" json:"dibuat_oleh"`
	DibuatOlehPengguna    Pengguna   `gorm:"foreignKey:DibuatOleh" json:"dibuat_oleh_pengguna,omitempty"`
	DisetujuiOleh         *uint      `gorm:"index;column:disetujui_oleh" json:"disetujui_oleh"`
	DisetujuiOlehPengguna *Pengguna  `gorm:"foreignKey:DisetujuiOleh" json:"disetujui_oleh_pengguna,omitempty"`
	DisetujuiPada         *time.Time `gorm:"column:disetujui_pada" json:"disetujui_pada"`
	DibuatPada            time.Time  `gorm:"column:dibuat_pada" json:"dibuat_pada"`
	DiperbaruiPada        time.Time  `gorm:"column:diperbarui_pada" json:"diperbarui_pada"`

	// Relationship
	Items []ItemPesananPembelian `gorm:"foreignKey:IDPO;constraint:OnDelete:CASCADE" json:"items,omitempty"`
}

// TableName mengembalikan nama tabel untuk model PesananPembelian
func (PesananPembelian) TableName() string {
	return "pesanan_pembelian"
}

// PurchaseOrder adalah alias untuk backward compatibility (akan dihapus nanti)
type PurchaseOrder = PesananPembelian

// ItemPesananPembelian adalah model untuk detail item Purchase Order
type ItemPesananPembelian struct {
	ID               uint             `gorm:"primaryKey;column:id" json:"id"`
	IDPO             uint             `gorm:"index;not null;column:id_po" json:"id_po"`
	PesananPembelian PesananPembelian `gorm:"foreignKey:IDPO" json:"pesanan_pembelian,omitempty"`
	IDProduk         uint             `gorm:"index;not null;column:id_produk" json:"id_produk"`
	Produk           Produk           `gorm:"foreignKey:IDProduk" json:"produk,omitempty"`
	Jumlah           int              `gorm:"not null;column:jumlah" json:"jumlah"`
	HargaSatuan      float64          `gorm:"type:decimal(15,2);not null;column:harga_satuan" json:"harga_satuan"`
	Subtotal         float64          `gorm:"type:decimal(15,2);not null;column:subtotal" json:"subtotal"` // quantity × unit_price
	JumlahDiterima   int              `gorm:"default:0;column:jumlah_diterima" json:"jumlah_diterima"`     // Jumlah yang sudah diterima
	DibuatPada       time.Time        `gorm:"column:dibuat_pada" json:"dibuat_pada"`
	DiperbaruiPada   time.Time        `gorm:"column:diperbarui_pada" json:"diperbarui_pada"`
}

// TableName mengembalikan nama tabel untuk model ItemPesananPembelian
func (ItemPesananPembelian) TableName() string {
	return "item_pesanan_pembelian"
}

// PurchaseOrderItem adalah alias untuk backward compatibility (akan dihapus nanti)
type PurchaseOrderItem = ItemPesananPembelian
