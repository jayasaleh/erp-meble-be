package models

import (
	"time"
)

// BarangMasuk adalah model untuk transaksi barang masuk
type BarangMasuk struct {
	ID                uint       `gorm:"primaryKey;column:id" json:"id"`
	TransactionNumber string     `gorm:"uniqueIndex;not null;column:nomor_transaksi" json:"transaction_number"`
	SupplierID        *uint      `gorm:"index;column:id_supplier" json:"supplier_id"`
	Supplier          *Pemasok   `gorm:"foreignKey:SupplierID" json:"supplier,omitempty"`
	POID              *uint      `gorm:"index;column:id_po" json:"po_id"` // Link ke Purchase Order
	ReceivedBy        uint       `gorm:"not null;column:diterima_oleh" json:"received_by"`
	ReceivedByUser    Pengguna   `gorm:"foreignKey:ReceivedBy" json:"received_by_user,omitempty"`
	ReceivedAt        time.Time  `gorm:"column:diterima_pada" json:"received_at"`
	ApprovedBy        *uint      `gorm:"index;column:disetujui_oleh" json:"approved_by"`
	ApprovedByUser    *Pengguna  `gorm:"foreignKey:ApprovedBy" json:"approved_by_user,omitempty"`
	ApprovedAt        *time.Time `gorm:"column:disetujui_pada" json:"approved_at"`
	Status            string     `gorm:"type:varchar(20);default:'pending';column:status" json:"status"` // pending, approved, rejected
	Notes             string     `gorm:"type:text;column:keterangan" json:"notes"`
	CreatedAt         time.Time  `gorm:"column:dibuat_pada" json:"created_at"`
	UpdatedAt         time.Time  `gorm:"column:diperbarui_pada" json:"updated_at"`

	// Relationship
	Items []ItemBarangMasuk `gorm:"foreignKey:StockInID;constraint:OnDelete:CASCADE" json:"items,omitempty"`
}

// TableName mengembalikan nama tabel untuk model BarangMasuk
func (BarangMasuk) TableName() string {
	return "barang_masuk"
}

// StockIn adalah alias untuk backward compatibility (akan dihapus nanti)
type StockIn = BarangMasuk

// ItemBarangMasuk adalah model untuk detail item barang masuk
type ItemBarangMasuk struct {
	ID          uint        `gorm:"primaryKey;column:id" json:"id"`
	StockInID   uint        `gorm:"index;not null;column:id_stock_in" json:"stock_in_id"`
	StockIn     BarangMasuk `gorm:"foreignKey:StockInID" json:"stock_in,omitempty"`
	ProductID   uint        `gorm:"index;not null;column:id_produk" json:"product_id"`
	Product     Produk      `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Quantity    int         `gorm:"not null;column:jumlah" json:"quantity"`
	UnitPrice   float64     `gorm:"type:decimal(15,2);not null;column:harga_satuan" json:"unit_price"`
	POPrice     *float64    `gorm:"type:decimal(15,2);column:harga_po" json:"po_price"` // Harga dari PO untuk perbandingan
	WarehouseID uint        `gorm:"index;not null;column:id_gudang" json:"warehouse_id"`
	Warehouse   Gudang      `gorm:"foreignKey:WarehouseID" json:"warehouse,omitempty"`
	Location    string      `gorm:"type:text;column:lokasi" json:"location"` // "Rak A, Slot B" - bisa jadi FK nanti
	CreatedAt   time.Time   `gorm:"column:dibuat_pada" json:"created_at"`
	UpdatedAt   time.Time   `gorm:"column:diperbarui_pada" json:"updated_at"`
}

// TableName mengembalikan nama tabel untuk model ItemBarangMasuk
func (ItemBarangMasuk) TableName() string {
	return "item_barang_masuk"
}

// StockInItem adalah alias untuk backward compatibility (akan dihapus nanti)
type StockInItem = ItemBarangMasuk

// BarangKeluar adalah model untuk transaksi barang keluar
type BarangKeluar struct {
	ID                uint      `gorm:"primaryKey;column:id" json:"id"`
	TransactionNumber string    `gorm:"uniqueIndex;not null;column:nomor_transaksi" json:"transaction_number"`
	Reason            string    `gorm:"type:varchar(50);not null;column:alasan" json:"reason"`        // penjualan, mutasi, produksi, rusak, adjustment
	ReferenceID       *uint     `gorm:"index;column:id_referensi" json:"reference_id"`                // Link ke sales, transfer, dll
	ReferenceType     string    `gorm:"type:varchar(50);column:tipe_referensi" json:"reference_type"` // sales, transfer, dll
	CreatedBy         uint      `gorm:"not null;column:dibuat_oleh" json:"created_by"`
	CreatedByUser     Pengguna  `gorm:"foreignKey:CreatedBy" json:"created_by_user,omitempty"`
	CreatedAt         time.Time `gorm:"column:dibuat_pada" json:"created_at"`
	UpdatedAt         time.Time `gorm:"column:diperbarui_pada" json:"updated_at"`

	// Relationship
	Items []ItemBarangKeluar `gorm:"foreignKey:StockOutID;constraint:OnDelete:CASCADE" json:"items,omitempty"`
}

// TableName mengembalikan nama tabel untuk model BarangKeluar
func (BarangKeluar) TableName() string {
	return "barang_keluar"
}

// StockOut adalah alias untuk backward compatibility (akan dihapus nanti)
type StockOut = BarangKeluar

// ItemBarangKeluar adalah model untuk detail item barang keluar
type ItemBarangKeluar struct {
	ID            uint             `gorm:"primaryKey;column:id" json:"id"`
	StockOutID    uint             `gorm:"index;not null;column:id_stock_out" json:"stock_out_id"`
	StockOut      BarangKeluar     `gorm:"foreignKey:StockOutID" json:"stock_out,omitempty"`
	ProductID     uint             `gorm:"index;not null;column:id_produk" json:"product_id"`
	Product       Produk           `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Quantity      int              `gorm:"not null;column:jumlah" json:"quantity"`
	WarehouseID   uint             `gorm:"index;not null;column:id_gudang" json:"warehouse_id"`
	Warehouse     Gudang           `gorm:"foreignKey:WarehouseID" json:"warehouse,omitempty"`
	StockInItemID *uint            `gorm:"index;column:id_stock_in_item" json:"stock_in_item_id"` // Untuk FIFO tracking
	StockInItem   *ItemBarangMasuk `gorm:"foreignKey:StockInItemID" json:"stock_in_item,omitempty"`
	CreatedAt     time.Time        `gorm:"column:dibuat_pada" json:"created_at"`
	UpdatedAt     time.Time        `gorm:"column:diperbarui_pada" json:"updated_at"`
}

// TableName mengembalikan nama tabel untuk model ItemBarangKeluar
func (ItemBarangKeluar) TableName() string {
	return "item_barang_keluar"
}

// StockOutItem adalah alias untuk backward compatibility (akan dihapus nanti)
type StockOutItem = ItemBarangKeluar

// StokInventori adalah model untuk stok saat ini per produk per gudang
type StokInventori struct {
	ID             uint       `gorm:"primaryKey;column:id" json:"id"`
	ProductID      uint       `gorm:"uniqueIndex:idx_product_warehouse;not null;column:id_produk" json:"product_id"`
	Product        Produk     `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	WarehouseID    uint       `gorm:"uniqueIndex:idx_product_warehouse;not null;column:id_gudang" json:"warehouse_id"`
	Warehouse      Gudang     `gorm:"foreignKey:WarehouseID" json:"warehouse,omitempty"`
	Quantity       int        `gorm:"not null;default:0;column:jumlah" json:"quantity"` // Stok saat ini
	LastMovementAt *time.Time `gorm:"column:pergerakan_terakhir_pada" json:"last_movement_at"`
	UpdatedAt      time.Time  `gorm:"column:diperbarui_pada" json:"updated_at"`
}

// TableName mengembalikan nama tabel untuk model StokInventori
func (StokInventori) TableName() string {
	return "stok_inventori"
}

// InventoryStock adalah alias untuk backward compatibility (akan dihapus nanti)
type InventoryStock = StokInventori

// PergerakanStok adalah model untuk kartu stok (semua pergerakan)
type PergerakanStok struct {
	ID            uint      `gorm:"primaryKey;column:id" json:"id"`
	ProductID     uint      `gorm:"index;not null;column:id_produk" json:"product_id"`
	Product       Produk    `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	WarehouseID   uint      `gorm:"index;not null;column:id_gudang" json:"warehouse_id"`
	Warehouse     Gudang    `gorm:"foreignKey:WarehouseID" json:"warehouse,omitempty"`
	MovementType  string    `gorm:"type:varchar(20);not null;column:tipe_pergerakan" json:"movement_type"` // in, out, transfer_in, transfer_out, adjustment
	ReferenceType string    `gorm:"type:varchar(50);not null;column:tipe_referensi" json:"reference_type"` // stock_in, stock_out, sales, transfer, opname
	ReferenceID   *uint     `gorm:"index;column:id_referensi" json:"reference_id"`
	Quantity      int       `gorm:"not null;column:jumlah" json:"quantity"`             // Positif untuk in, negatif untuk out
	BalanceAfter  int       `gorm:"not null;column:saldo_setelah" json:"balance_after"` // Running balance
	UserID        uint      `gorm:"index;not null;column:id_pengguna" json:"user_id"`
	User          Pengguna  `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Notes         string    `gorm:"type:text;column:keterangan" json:"notes"`
	CreatedAt     time.Time `gorm:"column:dibuat_pada" json:"created_at"`
}

// TableName mengembalikan nama tabel untuk model PergerakanStok
func (PergerakanStok) TableName() string {
	return "pergerakan_stok"
}

// StockMovement adalah alias untuk backward compatibility (akan dihapus nanti)
type StockMovement = PergerakanStok
