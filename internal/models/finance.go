package models

import (
	"time"
)

// HutangPemasok adalah model untuk hutang supplier
type HutangPemasok struct {
	ID                  uint       `gorm:"primaryKey;column:id" json:"id"`
	IDPemasok           uint       `gorm:"index;not null;column:id_supplier" json:"id_pemasok"`
	Pemasok             Pemasok    `gorm:"foreignKey:IDPemasok" json:"pemasok,omitempty"`
	IDPO                *uint      `gorm:"index;column:id_po" json:"id_po"`                                          // Link ke Purchase Order
	IDBarangMasuk       *uint      `gorm:"index;column:id_stock_in" json:"id_barang_masuk"`                          // Link ke Stock In
	Jumlah              float64    `gorm:"type:decimal(15,2);not null;column:jumlah" json:"jumlah"`                  // Jumlah hutang
	JumlahDibayar       float64    `gorm:"type:decimal(15,2);default:0;column:jumlah_dibayar" json:"jumlah_dibayar"` // Jumlah yang sudah dibayar
	SisaHutang          float64    `gorm:"type:decimal(15,2);not null;column:sisa_hutang" json:"sisa_hutang"`        // Sisa hutang
	JatuhTempo          *time.Time `gorm:"column:jatuh_tempo" json:"jatuh_tempo"`
	Status              string     `gorm:"type:varchar(20);default:'unpaid';column:status" json:"status"` // unpaid, partially_paid, paid
	BuktiPembayaran     string     `gorm:"type:text;column:bukti_pembayaran" json:"bukti_pembayaran"`     // Path ke file bukti bayar
	DibayarPada         *time.Time `gorm:"column:dibayar_pada" json:"dibayar_pada"`
	DibayarOleh         *uint      `gorm:"index;column:dibayar_oleh" json:"dibayar_oleh"`
	DibayarOlehPengguna *Pengguna  `gorm:"foreignKey:DibayarOleh" json:"dibayar_oleh_pengguna,omitempty"`
	DibuatPada          time.Time  `gorm:"column:dibuat_pada" json:"dibuat_pada"`
	DiperbaruiPada      time.Time  `gorm:"column:diperbarui_pada" json:"diperbarui_pada"`
}

// TableName mengembalikan nama tabel untuk model HutangPemasok
func (HutangPemasok) TableName() string {
	return "hutang_pemasok"
}

// SupplierDebt adalah alias untuk backward compatibility (akan dihapus nanti)
type SupplierDebt = HutangPemasok
