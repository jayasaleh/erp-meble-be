package dto

import "time"

// ===========================
// SALES REPORT DTOs
// ===========================

// SalesReportRequest adalah query params untuk semua laporan penjualan
type SalesReportRequest struct {
	TanggalDari   time.Time `form:"tanggal_dari" binding:"required" time_format:"2006-01-02"`
	TanggalSampai time.Time `form:"tanggal_sampai" binding:"required" time_format:"2006-01-02"`
	IDGudang      *uint     `form:"id_gudang"`
}

// ------- Report by Period -------

// SalesReportByPeriodResponse adalah response laporan per periode
type SalesReportByPeriodResponse struct {
	TanggalDari    time.Time              `json:"tanggal_dari"`
	TanggalSampai  time.Time              `json:"tanggal_sampai"`
	NamaGudang     string                 `json:"nama_gudang,omitempty"`
	TotalTransaksi int64                  `json:"total_transaksi"`
	TotalRevenue   float64                `json:"total_revenue"` // Total harga jual
	TotalCOGS      float64                `json:"total_cogs"`    // Total harga modal (FIFO)
	TotalLaba      float64                `json:"total_laba"`    // Revenue - COGS
	MarginPersen   float64                `json:"margin_persen"` // (Laba/Revenue) * 100
	TotalDiskon    float64                `json:"total_diskon"`
	RataRataPerTrx float64                `json:"rata_rata_per_transaksi"`
	PerHari        []SalesPerHariResponse `json:"per_hari"`
	PerMetodeBayar []MetodeBayarSummary   `json:"per_metode_bayar"`
}

// SalesPerHariResponse adalah breakdown penjualan per hari
type SalesPerHariResponse struct {
	Tanggal      string  `json:"tanggal"` // YYYY-MM-DD
	TotalTrx     int64   `json:"total_transaksi"`
	TotalRevenue float64 `json:"total_revenue"`
	TotalCOGS    float64 `json:"total_cogs"`
	TotalLaba    float64 `json:"total_laba"`
}

// MetodeBayarSummary adalah breakdown per metode pembayaran
type MetodeBayarSummary struct {
	Metode       string  `json:"metode"`
	TotalTrx     int64   `json:"total_transaksi"`
	TotalRevenue float64 `json:"total_revenue"`
}

// ------- Report by Product -------

// SalesReportByProductResponse adalah response laporan per produk
type SalesReportByProductResponse struct {
	TanggalDari   time.Time             `json:"tanggal_dari"`
	TanggalSampai time.Time             `json:"tanggal_sampai"`
	NamaGudang    string                `json:"nama_gudang,omitempty"`
	Products      []ProductSalesSummary `json:"products"`
}

// ProductSalesSummary adalah ringkasan penjualan per produk
type ProductSalesSummary struct {
	IDProduk      uint    `json:"id_produk"`
	SKU           string  `json:"sku"`
	NamaProduk    string  `json:"nama_produk"`
	Kategori      string  `json:"kategori"`
	JumlahTerjual int64   `json:"jumlah_terjual"` // Total unit terjual
	TotalRevenue  float64 `json:"total_revenue"`
	TotalCOGS     float64 `json:"total_cogs"`
	TotalLaba     float64 `json:"total_laba"`
	MarginPersen  float64 `json:"margin_persen"`
}

// ------- Report by Customer -------

// SalesReportByCustomerResponse adalah response laporan per pelanggan
type SalesReportByCustomerResponse struct {
	TanggalDari   time.Time              `json:"tanggal_dari"`
	TanggalSampai time.Time              `json:"tanggal_sampai"`
	Customers     []CustomerSalesSummary `json:"customers"`
}

// CustomerSalesSummary adalah ringkasan penjualan per pelanggan
type CustomerSalesSummary struct {
	NamaPelanggan   string  `json:"nama_pelanggan"`
	KontakPelanggan string  `json:"kontak_pelanggan"`
	TotalTransaksi  int64   `json:"total_transaksi"`
	TotalBelanja    float64 `json:"total_belanja"` // Total yang dibayar
	TotalLaba       float64 `json:"total_laba"`    // Laba dari customer ini
}

// ===========================
// RETURN REPORT DTOs
// ===========================

// ReturnReportRequest adalah query params untuk laporan retur (opsional period filter)
type ReturnReportRequest struct {
	TanggalDari   *time.Time `form:"tanggal_dari" time_format:"2006-01-02"`
	TanggalSampai *time.Time `form:"tanggal_sampai" time_format:"2006-01-02"`
}

// ReturnReportResponse adalah response rekap lengkap semua jenis retur
type ReturnReportResponse struct {
	Periode        *PeriodeInfo          `json:"periode,omitempty"` // nil = all-time
	ReturPenjualan ReturPenjualanSummary `json:"retur_penjualan"`   // Customer → Toko
	ReturPembelian ReturPembelianSummary `json:"retur_pembelian"`   // Toko → Supplier
	StokKarantina  []StokKarantinaItem   `json:"stok_karantina"`    // Batch retur yang belum diproses
	DampakTotal    ReturnImpactSummary   `json:"dampak_total"`      // Kalkulasi kerugian/keuntungan
}

// PeriodeInfo adalah info tanggal filter laporan
type PeriodeInfo struct {
	TanggalDari   time.Time `json:"tanggal_dari"`
	TanggalSampai time.Time `json:"tanggal_sampai"`
}

// ReturPenjualanSummary adalah rekap retur penjualan (customer → toko)
type ReturPenjualanSummary struct {
	TotalDokumen    int64         `json:"total_dokumen"`     // Jumlah dokumen retur
	TotalUnit       int64         `json:"total_unit"`        // Total unit barang diretur
	TotalNilaiRetur float64       `json:"total_nilai_retur"` // Uang yang dikembalikan ke customer
	TotalNilaiModal float64       `json:"total_nilai_modal"` // Estimasi nilai modal barang (harga jual yg diretur ≈ proxy)
	PerStatus       []StatusCount `json:"per_status"`        // pending, completed
	PerAlasan       []AlasanCount `json:"per_alasan"`        // rusak, tidak_sesuai, dsb
	PerMetode       []MetodeCount `json:"per_metode"`        // cash, transfer, tukar_barang
	UnitDiKarantina int64         `json:"unit_di_karantina"` // Total unit yang sudah approve tapi belum diproses
}

// ReturPembelianSummary adalah rekap retur pembelian (toko → supplier)
type ReturPembelianSummary struct {
	TotalDokumen    int64         `json:"total_dokumen"`
	TotalUnit       int64         `json:"total_unit"`
	TotalNilaiRetur float64       `json:"total_nilai_retur"` // Total nilai yang diklaim ke supplier
	PerStatus       []StatusCount `json:"per_status"`
	PerAlasan       []AlasanCount `json:"per_alasan"`
	PerMetode       []MetodeCount `json:"per_metode"` // potong_hutang, refund, tukar_barang
}

// StokKarantinaItem adalah detail stok retur yang sedang dikarantina
type StokKarantinaItem struct {
	IDProduk     uint      `json:"id_produk"`
	SKUProduk    string    `json:"sku_produk"`
	NamaProduk   string    `json:"nama_produk"`
	NamaGudang   string    `json:"nama_gudang"`
	JumlahUnit   int       `json:"jumlah_unit"`
	NilaiModal   float64   `json:"nilai_modal"` // HargaModal × Jumlah
	NomorRetur   string    `json:"nomor_retur"` // Dari keterangan batch
	TanggalMasuk time.Time `json:"tanggal_masuk"`
}

// ReturnImpactSummary adalah dampak keuangan keseluruhan dari semua retur
type ReturnImpactSummary struct {
	TotalUangKeluar        float64 `json:"total_uang_keluar"`        // Uang dikembalikan ke customer (retur penjualan)
	TotalStokHilang        float64 `json:"total_stok_hilang"`        // Nilai stok yang dikirim balik ke supplier (retur pembelian, COGS keluar)
	TotalKarantinaNilai    float64 `json:"total_karantina_nilai"`    // Nilai stok yang sedang dikarantina
	EstimasiKerugianBersih float64 `json:"estimasi_kerugian_bersih"` // UangKeluar + KarantinaNilai (belum bisa dijual) - StokHilang yg direclaim ke supplier
	Catatan                string  `json:"catatan"`
}

// StatusCount untuk breakdown per status
type StatusCount struct {
	Status string `json:"status"`
	Jumlah int64  `json:"jumlah"`
}

// AlasanCount untuk breakdown per alasan
type AlasanCount struct {
	Alasan string `json:"alasan"`
	Jumlah int64  `json:"jumlah"`
}

// MetodeCount untuk breakdown per metode pengembalian
type MetodeCount struct {
	Metode     string  `json:"metode"`
	Jumlah     int64   `json:"jumlah"`
	TotalNilai float64 `json:"total_nilai"`
}
