package models

import (
	"time"
)

// ReturPenjualan adalah model untuk retur penjualan (customer return barang)
type ReturPenjualan struct {
	ID                    uint       `gorm:"primaryKey;column:id" json:"id"`
	NomorRetur            string     `gorm:"uniqueIndex;not null;column:nomor_retur" json:"nomor_retur"`
	IDPenjualan           uint       `gorm:"index;not null;column:id_penjualan" json:"id_penjualan"`
	Penjualan             Penjualan  `gorm:"foreignKey:IDPenjualan" json:"penjualan,omitempty"`
	NamaPelanggan         string     `gorm:"type:varchar(100);column:nama_pelanggan" json:"nama_pelanggan"`
	KontakPelanggan       string     `gorm:"type:varchar(50);column:kontak_pelanggan" json:"kontak_pelanggan"`
	Alasan                string     `gorm:"type:varchar(100);not null;column:alasan" json:"alasan"` // rusak, tidak sesuai, cacat, dll
	Subtotal              float64    `gorm:"type:decimal(15,2);not null;column:subtotal" json:"subtotal"`
	Total                 float64    `gorm:"type:decimal(15,2);not null;column:total" json:"total"`
	MetodePengembalian    string     `gorm:"type:varchar(20);not null;column:metode_pengembalian" json:"metode_pengembalian"` // cash, transfer, tukar_barang
	JumlahPengembalian    float64    `gorm:"type:decimal(15,2);not null;column:jumlah_pengembalian" json:"jumlah_pengembalian"`
	Status                string     `gorm:"type:varchar(20);default:'pending';column:status" json:"status"` // pending, approved, rejected, completed
	Keterangan            string     `gorm:"type:text;column:keterangan" json:"keterangan"`
	DiprosesOleh          uint       `gorm:"index;not null;column:diproses_oleh" json:"diproses_oleh"`
	DiprosesOlehPengguna  Pengguna   `gorm:"foreignKey:DiprosesOleh" json:"diproses_oleh_pengguna,omitempty"`
	DisetujuiOleh         *uint      `gorm:"index;column:disetujui_oleh" json:"disetujui_oleh"`
	DisetujuiOlehPengguna *Pengguna  `gorm:"foreignKey:DisetujuiOleh" json:"disetujui_oleh_pengguna,omitempty"`
	DisetujuiPada         *time.Time `gorm:"column:disetujui_pada" json:"disetujui_pada"`
	DiprosesPada          time.Time  `gorm:"column:diproses_pada" json:"diproses_pada"`
	DibuatPada            time.Time  `gorm:"column:dibuat_pada" json:"dibuat_pada"`
	DiperbaruiPada        time.Time  `gorm:"column:diperbarui_pada" json:"diperbarui_pada"`

	// Relationship
	Items []ItemReturPenjualan `gorm:"foreignKey:IDReturPenjualan;constraint:OnDelete:CASCADE" json:"items,omitempty"`
}

// TableName mengembalikan nama tabel untuk model ReturPenjualan
func (ReturPenjualan) TableName() string {
	return "retur_penjualan"
}

// SalesReturn adalah alias untuk backward compatibility (akan dihapus nanti)
type SalesReturn = ReturPenjualan

// ItemReturPenjualan adalah model untuk detail item retur penjualan
type ItemReturPenjualan struct {
	ID               uint           `gorm:"primaryKey;column:id" json:"id"`
	IDReturPenjualan uint           `gorm:"index;not null;column:id_retur_penjualan" json:"id_retur_penjualan"`
	ReturPenjualan   ReturPenjualan `gorm:"foreignKey:IDReturPenjualan" json:"retur_penjualan,omitempty"`
	IDItemPenjualan  uint           `gorm:"index;not null;column:id_item_penjualan" json:"id_item_penjualan"`
	ItemPenjualan    ItemPenjualan  `gorm:"foreignKey:IDItemPenjualan" json:"item_penjualan,omitempty"`
	IDProduk         uint           `gorm:"index;not null;column:id_produk" json:"id_produk"`
	Produk           Produk         `gorm:"foreignKey:IDProduk" json:"produk,omitempty"`
	Jumlah           int            `gorm:"not null;column:jumlah" json:"jumlah"`
	HargaSatuan      float64        `gorm:"type:decimal(15,2);not null;column:harga_satuan" json:"harga_satuan"`
	Subtotal         float64        `gorm:"type:decimal(15,2);not null;column:subtotal" json:"subtotal"`
	IDGudang         uint           `gorm:"index;not null;column:id_gudang" json:"id_gudang"`
	Gudang           Gudang         `gorm:"foreignKey:IDGudang" json:"gudang,omitempty"`
	DibuatPada       time.Time      `gorm:"column:dibuat_pada" json:"dibuat_pada"`
	DiperbaruiPada   time.Time      `gorm:"column:diperbarui_pada" json:"diperbarui_pada"`
}

// TableName mengembalikan nama tabel untuk model ItemReturPenjualan
func (ItemReturPenjualan) TableName() string {
	return "item_retur_penjualan"
}

// SalesReturnItem adalah alias untuk backward compatibility (akan dihapus nanti)
type SalesReturnItem = ItemReturPenjualan

// ReturPembelian adalah model untuk retur pembelian (return ke supplier)
type ReturPembelian struct {
	ID                    uint              `gorm:"primaryKey;column:id" json:"id"`
	NomorRetur            string            `gorm:"uniqueIndex;not null;column:nomor_retur" json:"nomor_retur"`
	IDPO                  *uint             `gorm:"index;column:id_po" json:"id_po"`
	PesananPembelian      *PesananPembelian `gorm:"foreignKey:IDPO" json:"pesanan_pembelian,omitempty"`
	IDBarangMasuk         *uint             `gorm:"index;column:id_stock_in" json:"id_barang_masuk"`
	BarangMasuk           *BarangMasuk      `gorm:"foreignKey:IDBarangMasuk" json:"barang_masuk,omitempty"`
	IDPemasok             uint              `gorm:"index;not null;column:id_supplier" json:"id_pemasok"`
	Pemasok               Pemasok           `gorm:"foreignKey:IDPemasok" json:"pemasok,omitempty"`
	Alasan                string            `gorm:"type:varchar(100);not null;column:alasan" json:"alasan"` // rusak, tidak sesuai, cacat, dll
	Subtotal              float64           `gorm:"type:decimal(15,2);not null;column:subtotal" json:"subtotal"`
	Total                 float64           `gorm:"type:decimal(15,2);not null;column:total" json:"total"`
	MetodePengembalian    string            `gorm:"type:varchar(20);not null;column:metode_pengembalian" json:"metode_pengembalian"` // potong_hutang, refund, tukar_barang
	Status                string            `gorm:"type:varchar(20);default:'pending';column:status" json:"status"`                  // pending, approved, rejected, completed
	Keterangan            string            `gorm:"type:text;column:keterangan" json:"keterangan"`
	DibuatOleh            uint              `gorm:"index;not null;column:dibuat_oleh" json:"dibuat_oleh"`
	DibuatOlehPengguna    Pengguna          `gorm:"foreignKey:DibuatOleh" json:"dibuat_oleh_pengguna,omitempty"`
	DisetujuiOleh         *uint             `gorm:"index;column:disetujui_oleh" json:"disetujui_oleh"`
	DisetujuiOlehPengguna *Pengguna         `gorm:"foreignKey:DisetujuiOleh" json:"disetujui_oleh_pengguna,omitempty"`
	DisetujuiPada         *time.Time        `gorm:"column:disetujui_pada" json:"disetujui_pada"`
	DibuatPada            time.Time         `gorm:"column:dibuat_pada" json:"dibuat_pada"`
	DiperbaruiPada        time.Time         `gorm:"column:diperbarui_pada" json:"diperbarui_pada"`

	// Relationship
	Items []ItemReturPembelian `gorm:"foreignKey:IDReturPembelian;constraint:OnDelete:CASCADE" json:"items,omitempty"`
}

// TableName mengembalikan nama tabel untuk model ReturPembelian
func (ReturPembelian) TableName() string {
	return "retur_pembelian"
}

// PurchaseReturn adalah alias untuk backward compatibility (akan dihapus nanti)
type PurchaseReturn = ReturPembelian

// ItemReturPembelian adalah model untuk detail item retur pembelian
type ItemReturPembelian struct {
	ID                uint             `gorm:"primaryKey;column:id" json:"id"`
	IDReturPembelian  uint             `gorm:"index;not null;column:id_retur_pembelian" json:"id_retur_pembelian"`
	ReturPembelian    ReturPembelian   `gorm:"foreignKey:IDReturPembelian" json:"retur_pembelian,omitempty"`
	IDItemBarangMasuk *uint            `gorm:"index;column:id_stock_in_item" json:"id_item_barang_masuk"`
	ItemBarangMasuk   *ItemBarangMasuk `gorm:"foreignKey:IDItemBarangMasuk" json:"item_barang_masuk,omitempty"`
	IDProduk          uint             `gorm:"index;not null;column:id_produk" json:"id_produk"`
	Produk            Produk           `gorm:"foreignKey:IDProduk" json:"produk,omitempty"`
	Jumlah            int              `gorm:"not null;column:jumlah" json:"jumlah"`
	HargaSatuan       float64          `gorm:"type:decimal(15,2);not null;column:harga_satuan" json:"harga_satuan"`
	Subtotal          float64          `gorm:"type:decimal(15,2);not null;column:subtotal" json:"subtotal"`
	IDGudang          uint             `gorm:"index;not null;column:id_gudang" json:"id_gudang"`
	Gudang            Gudang           `gorm:"foreignKey:IDGudang" json:"gudang,omitempty"`
	DibuatPada        time.Time        `gorm:"column:dibuat_pada" json:"dibuat_pada"`
	DiperbaruiPada    time.Time        `gorm:"column:diperbarui_pada" json:"diperbarui_pada"`
}

// TableName mengembalikan nama tabel untuk model ItemReturPembelian
func (ItemReturPembelian) TableName() string {
	return "item_retur_pembelian"
}

// PurchaseReturnItem adalah alias untuk backward compatibility (akan dihapus nanti)
type PurchaseReturnItem = ItemReturPembelian
