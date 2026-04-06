package services

import (
	"fmt"
	"math"
	"real-erp-mebel/be/internal/dto"
	"real-erp-mebel/be/internal/models"
	"time"

	"gorm.io/gorm"
)

type ReportService interface {
	GetSalesReportByPeriod(req *dto.SalesReportRequest) (*dto.SalesReportByPeriodResponse, error)
	GetSalesReportByProduct(req *dto.SalesReportRequest) (*dto.SalesReportByProductResponse, error)
	GetSalesReportByCustomer(req *dto.SalesReportRequest) (*dto.SalesReportByCustomerResponse, error)
	GetReturnReport(req *dto.ReturnReportRequest) (*dto.ReturnReportResponse, error)
	GetStockReport(req *dto.StockReportRequest) (*dto.StockReportResponse, error)
}

type reportService struct {
	db *gorm.DB
}

func NewReportService(db *gorm.DB) ReportService {
	return &reportService{db: db}
}

// buildBaseQuery membangun query dasar penjualan dalam rentang tanggal & gudang yang diminta.
func (s *reportService) buildBaseQuery(req *dto.SalesReportRequest) *gorm.DB {
	startOfDay := time.Date(req.TanggalDari.Year(), req.TanggalDari.Month(), req.TanggalDari.Day(), 0, 0, 0, 0, time.Local)
	endOfDay := time.Date(req.TanggalSampai.Year(), req.TanggalSampai.Month(), req.TanggalSampai.Day(), 23, 59, 59, 999999999, time.Local)

	q := s.db.Model(&models.Penjualan{}).
		Where("status = ?", "completed").
		Where("dibuat_pada BETWEEN ? AND ?", startOfDay, endOfDay)

	if req.IDGudang != nil {
		q = q.Where("id_gudang = ?", *req.IDGudang)
	}
	return q
}

// GetSalesReportByPeriod mengembalikan ringkasan penjualan + breakdown per hari + per metode bayar.
func (s *reportService) GetSalesReportByPeriod(req *dto.SalesReportRequest) (*dto.SalesReportByPeriodResponse, error) {
	q := s.buildBaseQuery(req)

	// Ringkasan total
	type Summary struct {
		TotalTrx     int64
		TotalRevenue float64
		TotalCOGS    float64
		TotalDiskon  float64
	}
	var summary Summary
	if err := q.Select(
		"COUNT(*) as total_trx, " +
			"SUM(total) as total_revenue, " +
			"SUM(total_harga_modal) as total_cogs, " +
			"SUM(jumlah_diskon) as total_diskon",
	).Scan(&summary).Error; err != nil {
		return nil, err
	}

	laba := summary.TotalRevenue - summary.TotalCOGS
	margin := 0.0
	if summary.TotalRevenue > 0 {
		margin = math.Round(laba/summary.TotalRevenue*10000) / 100 // persen 2 desimal
	}
	rataRata := 0.0
	if summary.TotalTrx > 0 {
		rataRata = math.Round(summary.TotalRevenue/float64(summary.TotalTrx)*100) / 100
	}

	// Breakdown per hari
	type DailyRow struct {
		Tanggal      string
		TotalTrx     int64
		TotalRevenue float64
		TotalCOGS    float64
	}
	var dailyRows []DailyRow
	if err := s.buildBaseQuery(req).Select(
		"TO_CHAR(dibuat_pada, 'YYYY-MM-DD') as tanggal, " +
			"COUNT(*) as total_trx, " +
			"SUM(total) as total_revenue, " +
			"SUM(total_harga_modal) as total_cogs",
	).Group("tanggal").Order("tanggal ASC").Scan(&dailyRows).Error; err != nil {
		return nil, err
	}

	var perHari []dto.SalesPerHariResponse
	for _, row := range dailyRows {
		l := math.Round((row.TotalRevenue-row.TotalCOGS)*100) / 100
		perHari = append(perHari, dto.SalesPerHariResponse{
			Tanggal:      row.Tanggal,
			TotalTrx:     row.TotalTrx,
			TotalRevenue: row.TotalRevenue,
			TotalCOGS:    row.TotalCOGS,
			TotalLaba:    l,
		})
	}

	// Breakdown per metode bayar
	type MetodeRow struct {
		Metode       string
		TotalTrx     int64
		TotalRevenue float64
	}
	var metodeRows []MetodeRow
	if err := s.buildBaseQuery(req).Select(
		"metode_pembayaran as metode, COUNT(*) as total_trx, SUM(total) as total_revenue",
	).Group("metode").Scan(&metodeRows).Error; err != nil {
		return nil, err
	}

	var perMetode []dto.MetodeBayarSummary
	for _, row := range metodeRows {
		perMetode = append(perMetode, dto.MetodeBayarSummary{
			Metode:       row.Metode,
			TotalTrx:     row.TotalTrx,
			TotalRevenue: row.TotalRevenue,
		})
	}

	// Nama gudang jika difilter
	namaGudang := ""
	if req.IDGudang != nil {
		var g models.Gudang
		if err := s.db.First(&g, *req.IDGudang).Error; err == nil {
			namaGudang = g.Nama
		}
	}

	return &dto.SalesReportByPeriodResponse{
		TanggalDari:    req.TanggalDari,
		TanggalSampai:  req.TanggalSampai,
		NamaGudang:     namaGudang,
		TotalTransaksi: summary.TotalTrx,
		TotalRevenue:   math.Round(summary.TotalRevenue*100) / 100,
		TotalCOGS:      math.Round(summary.TotalCOGS*100) / 100,
		TotalLaba:      math.Round(laba*100) / 100,
		MarginPersen:   margin,
		TotalDiskon:    math.Round(summary.TotalDiskon*100) / 100,
		RataRataPerTrx: rataRata,
		PerHari:        perHari,
		PerMetodeBayar: perMetode,
	}, nil
}

// GetSalesReportByProduct mengembalikan ringkasan penjualan per produk.
func (s *reportService) GetSalesReportByProduct(req *dto.SalesReportRequest) (*dto.SalesReportByProductResponse, error) {
	startOfDay := time.Date(req.TanggalDari.Year(), req.TanggalDari.Month(), req.TanggalDari.Day(), 0, 0, 0, 0, time.Local)
	endOfDay := time.Date(req.TanggalSampai.Year(), req.TanggalSampai.Month(), req.TanggalSampai.Day(), 23, 59, 59, 999999999, time.Local)

	// Join item_penjualan → penjualan → produk
	baseQ := s.db.Table("item_penjualan ip").
		Joins("JOIN penjualan p ON p.id = ip.id_penjualan").
		Joins("JOIN produk pr ON pr.id = ip.id_produk").
		Where("p.status = ?", "completed").
		Where("p.dibuat_pada BETWEEN ? AND ?", startOfDay, endOfDay)

	if req.IDGudang != nil {
		baseQ = baseQ.Where("p.id_gudang = ?", *req.IDGudang)
	}

	type ProductRow struct {
		IDProduk      uint
		SKU           string
		NamaProduk    string
		Kategori      string
		JumlahTerjual int64
		TotalRevenue  float64
		TotalCOGS     float64
	}
	var rows []ProductRow

	err := baseQ.Select(fmt.Sprintf(
		"ip.id_produk, pr.sku, pr.nama as nama_produk, pr.kategori, " +
			"SUM(ip.jumlah) as jumlah_terjual, " +
			"SUM(ip.subtotal) as total_revenue, " +
			"SUM(ip.total_modal) as total_cogs",
	)).Group("ip.id_produk, pr.sku, pr.nama, pr.kategori").
		Order("jumlah_terjual DESC").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	var products []dto.ProductSalesSummary
	for _, row := range rows {
		l := math.Round((row.TotalRevenue-row.TotalCOGS)*100) / 100
		m := 0.0
		if row.TotalRevenue > 0 {
			m = math.Round(l/row.TotalRevenue*10000) / 100
		}
		products = append(products, dto.ProductSalesSummary{
			IDProduk:      row.IDProduk,
			SKU:           row.SKU,
			NamaProduk:    row.NamaProduk,
			Kategori:      row.Kategori,
			JumlahTerjual: row.JumlahTerjual,
			TotalRevenue:  row.TotalRevenue,
			TotalCOGS:     row.TotalCOGS,
			TotalLaba:     l,
			MarginPersen:  m,
		})
	}

	namaGudang := ""
	if req.IDGudang != nil {
		var g models.Gudang
		if err := s.db.First(&g, *req.IDGudang).Error; err == nil {
			namaGudang = g.Nama
		}
	}

	return &dto.SalesReportByProductResponse{
		TanggalDari:   req.TanggalDari,
		TanggalSampai: req.TanggalSampai,
		NamaGudang:    namaGudang,
		Products:      products,
	}, nil
}

// GetSalesReportByCustomer mengembalikan ringkasan penjualan per nama pelanggan.
func (s *reportService) GetSalesReportByCustomer(req *dto.SalesReportRequest) (*dto.SalesReportByCustomerResponse, error) {
	q := s.buildBaseQuery(req)

	type CustomerRow struct {
		NamaPelanggan   string
		KontakPelanggan string
		TotalTrx        int64
		TotalBelanja    float64
		TotalLaba       float64
	}
	var rows []CustomerRow

	err := q.Select(
		"COALESCE(NULLIF(nama_pelanggan,''), 'Umum') as nama_pelanggan, " +
			"kontak_pelanggan, " +
			"COUNT(*) as total_trx, " +
			"SUM(total) as total_belanja, " +
			"SUM(total - total_harga_modal) as total_laba",
	).Group("nama_pelanggan, kontak_pelanggan").
		Order("total_belanja DESC").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	var customers []dto.CustomerSalesSummary
	for _, row := range rows {
		customers = append(customers, dto.CustomerSalesSummary{
			NamaPelanggan:   row.NamaPelanggan,
			KontakPelanggan: row.KontakPelanggan,
			TotalTransaksi:  row.TotalTrx,
			TotalBelanja:    math.Round(row.TotalBelanja*100) / 100,
			TotalLaba:       math.Round(row.TotalLaba*100) / 100,
		})
	}

	return &dto.SalesReportByCustomerResponse{
		TanggalDari:   req.TanggalDari,
		TanggalSampai: req.TanggalSampai,
		Customers:     customers,
	}, nil
}

// GetReturnReport mengembalikan rekapitulasi data retur penjualan dan pembelian
func (s *reportService) GetReturnReport(req *dto.ReturnReportRequest) (*dto.ReturnReportResponse, error) {
	var startOfDay, endOfDay *time.Time
	if req.TanggalDari != nil && req.TanggalSampai != nil {
		sd := time.Date(req.TanggalDari.Year(), req.TanggalDari.Month(), req.TanggalDari.Day(), 0, 0, 0, 0, time.Local)
		ed := time.Date(req.TanggalSampai.Year(), req.TanggalSampai.Month(), req.TanggalSampai.Day(), 23, 59, 59, 999999999, time.Local)
		startOfDay = &sd
		endOfDay = &ed
	}

	// 1. RETUR PENJUALAN SUMMARY
	qJual := s.db.Model(&models.ReturPenjualan{})
	if startOfDay != nil {
		qJual = qJual.Where("dibuat_pada BETWEEN ? AND ?", *startOfDay, *endOfDay)
	}

	var totalDokJual int64
	var totalNilaiJual float64
	qJual.Select("COUNT(*), COALESCE(SUM(total), 0)").Row().Scan(&totalDokJual, &totalNilaiJual)

	var totalUnitJual int64
	qJualItems := s.db.Table("item_retur_penjualan ir").Joins("JOIN retur_penjualan r ON r.id = ir.id_retur_penjualan")
	if startOfDay != nil {
		qJualItems = qJualItems.Where("r.dibuat_pada BETWEEN ? AND ?", *startOfDay, *endOfDay)
	}
	qJualItems.Select("COALESCE(SUM(ir.jumlah), 0)").Row().Scan(&totalUnitJual)

	var perStatusJual []dto.StatusCount
	qJual.Select("status, COUNT(*) as jumlah").Group("status").Scan(&perStatusJual)

	var perAlasanJual []dto.AlasanCount
	qJual.Select("alasan, COUNT(*) as jumlah").Group("alasan").Scan(&perAlasanJual)

	var perMetodeJual []dto.MetodeCount
	qJual.Select("metode_pengembalian as metode, COUNT(*) as jumlah, COALESCE(SUM(total), 0) as total_nilai").Group("metode_pengembalian").Scan(&perMetodeJual)

	// Hitung unit di karantina dari retur penjualan
	var unitKarantina int64
	qKarantinaTotal := s.db.Model(&models.StokBatch{}).Where("tipe_referensi = ?", "retur_penjualan").Where("aktif = ?", false)
	if startOfDay != nil {
		qKarantinaTotal = qKarantinaTotal.Where("tanggal_masuk BETWEEN ? AND ?", *startOfDay, *endOfDay)
	}
	qKarantinaTotal.Select("COALESCE(SUM(jumlah_saat_ini), 0)").Row().Scan(&unitKarantina)

	resJual := dto.ReturPenjualanSummary{
		TotalDokumen:    totalDokJual,
		TotalUnit:       totalUnitJual,
		TotalNilaiRetur: math.Round(totalNilaiJual*100) / 100,
		TotalNilaiModal: math.Round(totalNilaiJual*100) / 100, // proxy modal = nilai retur
		PerStatus:       perStatusJual,
		PerAlasan:       perAlasanJual,
		PerMetode:       perMetodeJual,
		UnitDiKarantina: unitKarantina,
	}

	// 2. RETUR PEMBELIAN SUMMARY
	qBeli := s.db.Model(&models.ReturPembelian{})
	if startOfDay != nil {
		qBeli = qBeli.Where("dibuat_pada BETWEEN ? AND ?", *startOfDay, *endOfDay)
	}

	var totalDokBeli int64
	var totalNilaiBeli float64
	qBeli.Select("COUNT(*), COALESCE(SUM(total), 0)").Row().Scan(&totalDokBeli, &totalNilaiBeli)

	var totalUnitBeli int64
	qBeliItems := s.db.Table("item_retur_pembelian ir").Joins("JOIN retur_pembelian r ON r.id = ir.id_retur_pembelian")
	if startOfDay != nil {
		qBeliItems = qBeliItems.Where("r.dibuat_pada BETWEEN ? AND ?", *startOfDay, *endOfDay)
	}
	qBeliItems.Select("COALESCE(SUM(ir.jumlah), 0)").Row().Scan(&totalUnitBeli)

	var perStatusBeli []dto.StatusCount
	qBeli.Select("status, COUNT(*) as jumlah").Group("status").Scan(&perStatusBeli)

	var perAlasanBeli []dto.AlasanCount
	qBeli.Select("alasan, COUNT(*) as jumlah").Group("alasan").Scan(&perAlasanBeli)

	var perMetodeBeli []dto.MetodeCount
	qBeli.Select("metode_pengembalian as metode, COUNT(*) as jumlah, COALESCE(SUM(total), 0) as total_nilai").Group("metode_pengembalian").Scan(&perMetodeBeli)

	resBeli := dto.ReturPembelianSummary{
		TotalDokumen:    totalDokBeli,
		TotalUnit:       totalUnitBeli,
		TotalNilaiRetur: math.Round(totalNilaiBeli*100) / 100,
		PerStatus:       perStatusBeli,
		PerAlasan:       perAlasanBeli,
		PerMetode:       perMetodeBeli,
	}

	// 3. STOK KARANTINA LIST
	qKarantina := s.db.Table("stok_batch b").
		Joins("JOIN produk p ON p.id = b.id_produk").
		Joins("JOIN gudang g ON g.id = b.id_gudang").
		Where("b.tipe_referensi = ?", "retur_penjualan").
		Where("b.aktif = ?", false).
		Where("b.jumlah_saat_ini > ?", 0)

	if startOfDay != nil {
		qKarantina = qKarantina.Where("b.tanggal_masuk BETWEEN ? AND ?", *startOfDay, *endOfDay)
	}

	type KarantinaRow struct {
		IDProduk     uint
		SKUProduk    string
		NamaProduk   string
		NamaGudang   string
		JumlahUnit   int
		HargaModal   float64
		Keterangan   string
		TanggalMasuk time.Time
	}
	var kRows []KarantinaRow
	qKarantina.Select("b.id_produk, p.sku as sku_produk, p.nama as nama_produk, g.nama as nama_gudang, b.jumlah_saat_ini as jumlah_unit, b.harga_modal, b.keterangan, b.tanggal_masuk").
		Order("b.tanggal_masuk DESC").
		Scan(&kRows)

	var stokKarantina []dto.StokKarantinaItem
	var totalKarantinaNilai float64
	for _, k := range kRows {
		nilai := k.HargaModal * float64(k.JumlahUnit)
		totalKarantinaNilai += nilai

		stokKarantina = append(stokKarantina, dto.StokKarantinaItem{
			IDProduk:     k.IDProduk,
			SKUProduk:    k.SKUProduk,
			NamaProduk:   k.NamaProduk,
			NamaGudang:   k.NamaGudang,
			JumlahUnit:   k.JumlahUnit,
			NilaiModal:   math.Round(nilai*100) / 100,
			NomorRetur:   k.Keterangan, // Contains: [KARANTINA-RETUR] SKU | Nomor Retur: RETP/...
			TanggalMasuk: k.TanggalMasuk,
		})
	}

	// 4. DAMPAK TOTAL
	dampak := dto.ReturnImpactSummary{
		TotalUangKeluar:        resJual.TotalNilaiRetur, // Yang dibayarkan ke customer (estimasi uang/nilai)
		TotalStokHilang:        resBeli.TotalNilaiRetur, // Nilai barang yang dikembalikan ke supplier
		TotalKarantinaNilai:    math.Round(totalKarantinaNilai*100) / 100,
		EstimasiKerugianBersih: math.Round((resJual.TotalNilaiRetur+totalKarantinaNilai-resBeli.TotalNilaiRetur)*100) / 100,
		Catatan:                "Estimasi kerugian kotor (uang keluar - klaim supplier). Barang karantina yang belum dimusnahkan tetap ada nilainya.",
	}
	if dampak.EstimasiKerugianBersih < 0 {
		dampak.EstimasiKerugianBersih = 0
		dampak.Catatan = "Refund dari supplier menutupi nilai kerugian retur pelanggan."
	}

	var periode *dto.PeriodeInfo
	if req.TanggalDari != nil && req.TanggalSampai != nil {
		periode = &dto.PeriodeInfo{
			TanggalDari:   *startOfDay,
			TanggalSampai: *endOfDay,
		}
	}

	return &dto.ReturnReportResponse{
		Periode:        periode,
		ReturPenjualan: resJual,
		ReturPembelian: resBeli,
		StokKarantina:  stokKarantina,
		DampakTotal:    dampak,
	}, nil
}

// GetStockReport mengembalikan laporan stok (keseluruhan maupun peringatan stok menipis) lengkap dengan pagination
func (s *reportService) GetStockReport(req *dto.StockReportRequest) (*dto.StockReportResponse, error) {
	threshold := req.Threshold
	if threshold <= 0 {
		threshold = 5
	}

	// Query dasar dari tabel produk join stok_batch untuk agregasi
	buildBase := func() *gorm.DB {
		q := s.db.Table("produk p").
			Select("p.id as id_produk, p.sku, p.nama as nama_produk, p.kategori, " +
				"COALESCE(SUM(b.jumlah_saat_ini), 0) as total_stok, " +
				"COALESCE(SUM(b.jumlah_saat_ini * b.harga_modal), 0) as valuasi_modal").
			Joins("LEFT JOIN stok_batch b ON p.id = b.id_produk AND b.aktif = true").
			Where("p.dihapus_pada IS NULL").
			Group("p.id, p.sku, p.nama, p.kategori")

		if req.Search != "" {
			searchLike := "%" + req.Search + "%"
			q = q.Where("p.nama ILIKE ? OR p.sku ILIKE ?", searchLike, searchLike)
		}

		if req.IDProduk != nil {
			q = q.Where("p.id = ?", *req.IDProduk)
		}

		if req.LowStockOnly {
			q = q.Having("COALESCE(SUM(b.jumlah_saat_ini), 0) <= ?", threshold)
		}

		return q
	}

	// Count total rows using a raw count over the subquery
	var count int64
	var totalValuasi float64

	type countRow struct {
		Total int64
	}
	var countResults []dto.StockReportItem
	if err := buildBase().Scan(&countResults).Error; err != nil {
		return nil, err
	}
	count = int64(len(countResults))
	for _, r := range countResults {
		totalValuasi += r.ValuasiModal
	}

	// Kalkulasi Pagination
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 {
		req.Limit = 10
	}

	offset := (req.Page - 1) * req.Limit
	totalPage := int(math.Ceil(float64(count) / float64(req.Limit)))

	// Tentukan Order & paginate
	baseQ := buildBase()
	if req.LowStockOnly {
		baseQ = baseQ.Order("total_stok ASC") // Menipis paling atas
	} else {
		baseQ = baseQ.Order("p.nama ASC")
	}

	var rows []dto.StockReportItem
	if err := baseQ.Offset(offset).Limit(req.Limit).Scan(&rows).Error; err != nil {
		return nil, err
	}

	// Pembulatan desimal valuasi agar rapi
	for i := range rows {
		rows[i].ValuasiModal = math.Round(rows[i].ValuasiModal*100) / 100
	}

	return &dto.StockReportResponse{
		TotalData:    count,
		TotalPage:    totalPage,
		Page:         req.Page,
		Limit:        req.Limit,
		TotalValuasi: math.Round(totalValuasi*100) / 100,
		Data:         rows,
	}, nil
}
