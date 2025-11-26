package models

import (
	"time"
)

// BarangMasuk adalah model untuk transaksi barang masuk
type BarangMasuk struct {
	ID                    uint       `gorm:"primaryKey;column:id" json:"id"`
	NomorTransaksi        string     `gorm:"uniqueIndex;not null;column:nomor_transaksi" json:"nomor_transaksi"`
	IDPemasok             *uint      `gorm:"index;column:id_supplier" json:"id_pemasok"`
	Pemasok               *Pemasok   `gorm:"foreignKey:IDPemasok" json:"pemasok,omitempty"`
	IDPO                  *uint      `gorm:"index;column:id_po" json:"id_po"` // Link ke Purchase Order
	DiterimaOleh          uint       `gorm:"not null;column:diterima_oleh" json:"diterima_oleh"`
	DiterimaOlehPengguna  Pengguna   `gorm:"foreignKey:DiterimaOleh" json:"diterima_oleh_pengguna,omitempty"`
	DiterimaPada          time.Time  `gorm:"column:diterima_pada" json:"diterima_pada"`
	DisetujuiOleh         *uint      `gorm:"index;column:disetujui_oleh" json:"disetujui_oleh"`
	DisetujuiOlehPengguna *Pengguna  `gorm:"foreignKey:DisetujuiOleh" json:"disetujui_oleh_pengguna,omitempty"`
	DisetujuiPada         *time.Time `gorm:"column:disetujui_pada" json:"disetujui_pada"`
	Status                string     `gorm:"type:varchar(20);default:'pending';column:status" json:"status"` // pending, approved, rejected
	Keterangan            string     `gorm:"type:text;column:keterangan" json:"keterangan"`
	DibuatPada            time.Time  `gorm:"column:dibuat_pada" json:"dibuat_pada"`
	DiperbaruiPada        time.Time  `gorm:"column:diperbarui_pada" json:"diperbarui_pada"`

	// Relationship
	Items []ItemBarangMasuk `gorm:"foreignKey:IDBarangMasuk;constraint:OnDelete:CASCADE" json:"items,omitempty"`
}

// TableName mengembalikan nama tabel untuk model BarangMasuk
func (BarangMasuk) TableName() string {
	return "barang_masuk"
}

// StockIn adalah alias untuk backward compatibility (akan dihapus nanti)
type StockIn = BarangMasuk

// ItemBarangMasuk adalah model untuk detail item barang masuk
type ItemBarangMasuk struct {
	ID             uint        `gorm:"primaryKey;column:id" json:"id"`
	IDBarangMasuk  uint        `gorm:"index;not null;column:id_stock_in" json:"id_barang_masuk"`
	BarangMasuk    BarangMasuk `gorm:"foreignKey:IDBarangMasuk" json:"barang_masuk,omitempty"`
	IDProduk       uint        `gorm:"index;not null;column:id_produk" json:"id_produk"`
	Produk         Produk      `gorm:"foreignKey:IDProduk" json:"produk,omitempty"`
	Jumlah         int         `gorm:"not null;column:jumlah" json:"jumlah"`
	HargaSatuan    float64     `gorm:"type:decimal(15,2);not null;column:harga_satuan" json:"harga_satuan"`
	HargaPO        *float64    `gorm:"type:decimal(15,2);column:harga_po" json:"harga_po"` // Harga dari PO untuk perbandingan
	IDGudang       uint        `gorm:"index;not null;column:id_gudang" json:"id_gudang"`
	Gudang         Gudang      `gorm:"foreignKey:IDGudang" json:"gudang,omitempty"`
	Lokasi         string      `gorm:"type:text;column:lokasi" json:"lokasi"` // "Rak A, Slot B" - bisa jadi FK nanti
	DibuatPada     time.Time   `gorm:"column:dibuat_pada" json:"dibuat_pada"`
	DiperbaruiPada time.Time   `gorm:"column:diperbarui_pada" json:"diperbarui_pada"`
}

// TableName mengembalikan nama tabel untuk model ItemBarangMasuk
func (ItemBarangMasuk) TableName() string {
	return "item_barang_masuk"
}

// StockInItem adalah alias untuk backward compatibility (akan dihapus nanti)
type StockInItem = ItemBarangMasuk

// BarangKeluar adalah model untuk transaksi barang keluar
type BarangKeluar struct {
	ID                 uint      `gorm:"primaryKey;column:id" json:"id"`
	NomorTransaksi     string    `gorm:"uniqueIndex;not null;column:nomor_transaksi" json:"nomor_transaksi"`
	Alasan             string    `gorm:"type:varchar(50);not null;column:alasan" json:"alasan"`        // penjualan, mutasi, produksi, rusak, adjustment
	IDReferensi        *uint     `gorm:"index;column:id_referensi" json:"id_referensi"`                // Link ke sales, transfer, dll
	TipeReferensi      string    `gorm:"type:varchar(50);column:tipe_referensi" json:"tipe_referensi"` // sales, transfer, dll
	DibuatOleh         uint      `gorm:"not null;column:dibuat_oleh" json:"dibuat_oleh"`
	DibuatOlehPengguna Pengguna  `gorm:"foreignKey:DibuatOleh" json:"dibuat_oleh_pengguna,omitempty"`
	DibuatPada         time.Time `gorm:"column:dibuat_pada" json:"dibuat_pada"`
	DiperbaruiPada     time.Time `gorm:"column:diperbarui_pada" json:"diperbarui_pada"`

	// Relationship
	Items []ItemBarangKeluar `gorm:"foreignKey:IDBarangKeluar;constraint:OnDelete:CASCADE" json:"items,omitempty"`
}

// TableName mengembalikan nama tabel untuk model BarangKeluar
func (BarangKeluar) TableName() string {
	return "barang_keluar"
}

// StockOut adalah alias untuk backward compatibility (akan dihapus nanti)
type StockOut = BarangKeluar

// ItemBarangKeluar adalah model untuk detail item barang keluar
type ItemBarangKeluar struct {
	ID                uint             `gorm:"primaryKey;column:id" json:"id"`
	IDBarangKeluar    uint             `gorm:"index;not null;column:id_stock_out" json:"id_barang_keluar"`
	BarangKeluar      BarangKeluar     `gorm:"foreignKey:IDBarangKeluar" json:"barang_keluar,omitempty"`
	IDProduk          uint             `gorm:"index;not null;column:id_produk" json:"id_produk"`
	Produk            Produk           `gorm:"foreignKey:IDProduk" json:"produk,omitempty"`
	Jumlah            int              `gorm:"not null;column:jumlah" json:"jumlah"`
	IDGudang          uint             `gorm:"index;not null;column:id_gudang" json:"id_gudang"`
	Gudang            Gudang           `gorm:"foreignKey:IDGudang" json:"gudang,omitempty"`
	IDItemBarangMasuk *uint            `gorm:"index;column:id_stock_in_item" json:"id_item_barang_masuk"` // Untuk FIFO tracking
	ItemBarangMasuk   *ItemBarangMasuk `gorm:"foreignKey:IDItemBarangMasuk" json:"item_barang_masuk,omitempty"`
	DibuatPada        time.Time        `gorm:"column:dibuat_pada" json:"dibuat_pada"`
	DiperbaruiPada    time.Time        `gorm:"column:diperbarui_pada" json:"diperbarui_pada"`
}

// TableName mengembalikan nama tabel untuk model ItemBarangKeluar
func (ItemBarangKeluar) TableName() string {
	return "item_barang_keluar"
}

// StockOutItem adalah alias untuk backward compatibility (akan dihapus nanti)
type StockOutItem = ItemBarangKeluar

// StokInventori adalah model untuk stok saat ini per produk per gudang
type StokInventori struct {
	ID                     uint       `gorm:"primaryKey;column:id" json:"id"`
	IDProduk               uint       `gorm:"uniqueIndex:idx_product_warehouse;not null;column:id_produk" json:"id_produk"`
	Produk                 Produk     `gorm:"foreignKey:IDProduk" json:"produk,omitempty"`
	IDGudang               uint       `gorm:"uniqueIndex:idx_product_warehouse;not null;column:id_gudang" json:"id_gudang"`
	Gudang                 Gudang     `gorm:"foreignKey:IDGudang" json:"gudang,omitempty"`
	Jumlah                 int        `gorm:"not null;default:0;column:jumlah" json:"jumlah"` // Stok saat ini
	PergerakanTerakhirPada *time.Time `gorm:"column:pergerakan_terakhir_pada" json:"pergerakan_terakhir_pada"`
	DiperbaruiPada         time.Time  `gorm:"column:diperbarui_pada" json:"diperbarui_pada"`
}

// TableName mengembalikan nama tabel untuk model StokInventori
func (StokInventori) TableName() string {
	return "stok_inventori"
}

// InventoryStock adalah alias untuk backward compatibility (akan dihapus nanti)
type InventoryStock = StokInventori

// PergerakanStok adalah model untuk kartu stok (semua pergerakan)
type PergerakanStok struct {
	ID             uint      `gorm:"primaryKey;column:id" json:"id"`
	IDProduk       uint      `gorm:"index;not null;column:id_produk" json:"id_produk"`
	Produk         Produk    `gorm:"foreignKey:IDProduk" json:"produk,omitempty"`
	IDGudang       uint      `gorm:"index;not null;column:id_gudang" json:"id_gudang"`
	Gudang         Gudang    `gorm:"foreignKey:IDGudang" json:"gudang,omitempty"`
	TipePergerakan string    `gorm:"type:varchar(20);not null;column:tipe_pergerakan" json:"tipe_pergerakan"` // in, out, transfer_in, transfer_out, adjustment
	TipeReferensi  string    `gorm:"type:varchar(50);not null;column:tipe_referensi" json:"tipe_referensi"`   // stock_in, stock_out, sales, transfer, opname
	IDReferensi    *uint     `gorm:"index;column:id_referensi" json:"id_referensi"`
	Jumlah         int       `gorm:"not null;column:jumlah" json:"jumlah"`               // Positif untuk in, negatif untuk out
	SaldoSetelah   int       `gorm:"not null;column:saldo_setelah" json:"saldo_setelah"` // Running balance
	IDPengguna     uint      `gorm:"index;not null;column:id_pengguna" json:"id_pengguna"`
	Pengguna       Pengguna  `gorm:"foreignKey:IDPengguna" json:"pengguna,omitempty"`
	Keterangan     string    `gorm:"type:text;column:keterangan" json:"keterangan"`
	DibuatPada     time.Time `gorm:"column:dibuat_pada" json:"dibuat_pada"`
}

// TableName mengembalikan nama tabel untuk model PergerakanStok
func (PergerakanStok) TableName() string {
	return "pergerakan_stok"
}

// StockMovement adalah alias untuk backward compatibility (akan dihapus nanti)
type StockMovement = PergerakanStok
