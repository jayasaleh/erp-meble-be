package seed_test

import (
	"fmt"
	"log"
	"math/rand"
	"testing"
	"time"

	"real-erp-mebel/be/internal/config"
	"real-erp-mebel/be/internal/database"
	"real-erp-mebel/be/internal/models"

	"golang.org/x/crypto/bcrypt"
)

func TestSeedDB(t *testing.T) {
	// 1. Setup
	rand.Seed(time.Now().UnixNano())
	log.Println("Starting Database Seeder (3-month simulation)...")

	config.LoadConfig()
	if err := database.Connect(); err != nil {
		t.Fatalf("Failed to connect to DB: %v", err)
	}
	defer database.Close()
	db := database.DB

	// 2. Clear old data forcefully for pure seeder
	log.Println("Clearing database...")
	db.Exec("TRUNCATE pengguna, gudang, pemasok, produk, gambar_produk, pesanan_pembelian, item_pesanan_pembelian, barang_masuk, item_barang_masuk, hutang_pemasok, stok_inventori, stok_batch, pergerakan_stok, penjualan, item_penjualan, item_penjualan_batch, retur_penjualan, item_retur_penjualan, retur_pembelian, item_retur_pembelian RESTART IDENTITY CASCADE;")

	// 3. Generate Users
	log.Println("Seeding Users...")
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	users := []models.Pengguna{
		{Email: "owner@mebel.com", Nama: "Owner", Peran: "owner", Password: string(hashedPassword), Aktif: true, DibuatPada: time.Now()},
		{Email: "admin@mebel.com", Nama: "Admin Gudang", Peran: "admin_gudang", Password: string(hashedPassword), Aktif: true, DibuatPada: time.Now()},
		{Email: "kasir@mebel.com", Nama: "Kasir", Peran: "kasir", Password: string(hashedPassword), Aktif: true, DibuatPada: time.Now()},
		{Email: "finance@mebel.com", Nama: "Finance", Peran: "finance", Password: string(hashedPassword), Aktif: true, DibuatPada: time.Now()},
	}
	if err := db.Create(&users).Error; err != nil {
		t.Fatalf("Failed to seed users: %v", err)
	}

	// 4. Generate Warehouse
	log.Println("Seeding Warehouses...")
	gudang := models.Gudang{
		Kode: "GDG-001", Nama: "Gudang Utama", Alamat: "Jl. Industri No 1", Keterangan: "Gudang Pusat", Aktif: true, DibuatPada: time.Now(),
	}
	if err := db.Create(&gudang).Error; err != nil {
		t.Fatalf("Failed to seed warehouse: %v", err)
	}

	// 5. Generate Suppliers
	log.Println("Seeding Suppliers...")
	pemasokList := []models.Pemasok{
		{Nama: "PT Jati Kusuma", Kontak: "Budi", Telepon: "081111111", Email: "jati@kusuma.com", Alamat: "Jepara", Aktif: true, DibuatPada: time.Now()},
		{Nama: "CV Mebel Nusantara", Kontak: "Andi", Telepon: "082222222", Email: "mebel@nusantara.com", Alamat: "Surabaya", Aktif: true, DibuatPada: time.Now()},
	}
	if err := db.Create(&pemasokList).Error; err != nil {
		t.Fatalf("Failed to seed suppliers: %v", err)
	}

	// 6. Generate Products
	log.Println("Seeding Products...")
	kategori := []string{"Kursi", "Meja", "Lemari"}
	var produkList []models.Produk
	for i := 1; i <= 15; i++ {
		pemasokID := pemasokList[rand.Intn(len(pemasokList))].ID
		hargaModal := float64(100000 + rand.Intn(900000))
		hargaJual := hargaModal * 1.5 // 50% margin
		bc := fmt.Sprintf("BRC-%04d", i)
		produkList = append(produkList, models.Produk{
			SKU:         fmt.Sprintf("SKU-%04d", i),
			Barcode:     &bc,
			Nama:        fmt.Sprintf("Produk Mebel %d", i),
			Kategori:    kategori[rand.Intn(len(kategori))],
			Merek:       "Mebelku",
			IDPemasok:   &pemasokID,
			HargaModal:  hargaModal,
			HargaJual:   hargaJual,
			StokMinimum: 5,
			IzinDiskon:  true,
			Aktif:       true,
			DibuatOleh:  users[1].ID, // Admin gudang
			DiupdateOleh: users[1].ID,
			DibuatPada:  time.Now(),
		})
	}
	if err := db.Create(&produkList).Error; err != nil {
		t.Fatalf("Failed to seed products: %v", err)
	}

	// Images for products
	for _, p := range produkList {
		db.Create(&models.GambarProduk{
			IDProduk:    p.ID,
			PathGambar:  "https://images.unsplash.com/photo-1592078615290-033ee584e267?q=80&w=640&auto=format&fit=crop",
			GambarUtama: true,
			DibuatPada:  time.Now(),
		})
	}

	// 7. Initialize Stocks (Empty to begin with)
	var stokList []models.StokInventori
	for _, p := range produkList {
		stokList = append(stokList, models.StokInventori{
			IDProduk: p.ID,
			IDGudang: gudang.ID,
			Jumlah:   0,
		})
	}
	if err := db.Create(&stokList).Error; err != nil {
		t.Fatalf("Failed to seed initial stock: %v", err)
	}

	// Simulation Variables
	startDate := time.Now().AddDate(0, -3, 0) // Start 3 months ago
	endDate := time.Now()
	
	poCounter, inCounter, outCounter, salesCounter, returCounter := 1, 1, 1, 1, 1

	log.Println("Simulating Transactions Over 3 Months...")

	// 8. Iterate day by day for 3 months
	for d := startDate; d.Before(endDate); d = d.AddDate(0, 0, 1) {
		
		// 8.a - Pembelian (Restock) every 7 days OR at start
		if d == startDate || d.Day()%7 == 0 {
			sup := pemasokList[rand.Intn(len(pemasokList))]
			po := models.PesananPembelian{
				NomorPO:      fmt.Sprintf("PO-%04d", poCounter),
				IDPemasok:    sup.ID,
				TanggalPesan: d,
				Status:       "completed",
				Total:        0,
				DibuatOleh:   users[1].ID,
				DisetujuiOleh: &users[0].ID,
				DibuatPada:   d,
			}
			poCounter++
			
			var poItems []models.ItemPesananPembelian
			totalPO := 0.0
			for i := 0; i < 5; i++ {
				p := produkList[rand.Intn(len(produkList))]
				qty := 10 + rand.Intn(40)
				subtotal := p.HargaModal * float64(qty)
				totalPO += subtotal
				
				poItems = append(poItems, models.ItemPesananPembelian{
					IDProduk:       p.ID,
					Jumlah:         qty,
					JumlahDiterima: qty,
					HargaSatuan:    p.HargaModal,
					Subtotal:       subtotal,
					DibuatPada:     d,
				})
			}
			po.Total = totalPO
			if err := db.Create(&po).Error; err != nil {
				t.Fatalf("Failed to create PO at %v: %v", d, err)
			}
			for i := range poItems {
				poItems[i].IDPO = po.ID
			}
			if err := db.Create(&poItems).Error; err != nil {
				t.Fatalf("Failed to create PO Items at %v: %v", d, err)
			}

			inbound := models.BarangMasuk{
				NomorTransaksi: fmt.Sprintf("IN-%04d", inCounter),
				IDPemasok:      &sup.ID,
				IDPO:           &po.ID,
				DiterimaOleh:   users[1].ID,
				DiterimaPada:   d,
				DisetujuiOleh:  &users[1].ID,
				Status:         "approved",
				Keterangan:     "Barang dari PO",
				DibuatPada:     d,
			}
			inCounter++
			if err := db.Create(&inbound).Error; err != nil {
				t.Fatalf("Failed to create BarangMasuk at %v: %v", d, err)
			}
			
			var inItems []models.ItemBarangMasuk
			for _, pi := range poItems {
				inItem := models.ItemBarangMasuk{
					IDBarangMasuk: inbound.ID,
					IDProduk:      pi.IDProduk,
					Jumlah:        pi.Jumlah,
					HargaSatuan:   pi.HargaSatuan,
					IDGudang:      gudang.ID,
					Lokasi:        "Rak Utama",
					DibuatPada:    d,
				}
				inItems = append(inItems, inItem)
			}
			if err := db.Create(&inItems).Error; err != nil {
				t.Fatalf("Failed to create ItemBarangMasuk at %v: %v", d, err)
			}
			
			for _, inI := range inItems {
				batch := models.StokBatch{
					IDProduk:       inI.IDProduk,
					IDGudang:       inI.IDGudang,
					TanggalMasuk:   d,
					JumlahAwal:     inI.Jumlah,
					JumlahSaatIni:  inI.Jumlah,
					HargaModal:     inI.HargaSatuan,
					IDReferensi:    &inI.ID,
					TipeReferensi:  "item_barang_masuk",
					Aktif:          true,
					DibuatPada:     d,
				}
				if err := db.Create(&batch).Error; err != nil {
					t.Fatalf("Failed to create StokBatch at %v: %v", d, err)
				}
				
				err := db.Exec("UPDATE stok_inventori SET jumlah = jumlah + ?, pergerakan_terakhir_pada = ? WHERE id_produk = ? AND id_gudang = ?", inI.Jumlah, d, inI.IDProduk, inI.IDGudang).Error
				if err != nil {
					t.Fatalf("Failed to update StokInventori at %v: %v", d, err)
				}

				if err := db.Create(&models.PergerakanStok{
					IDProduk:       inI.IDProduk,
					IDGudang:       inI.IDGudang,
					IDBatch:        &batch.ID,
					TipePergerakan: "in",
					TipeReferensi:  "stock_in",
					IDReferensi:    &inbound.ID,
					Jumlah:         inI.Jumlah,
					IDPengguna:     users[1].ID,
					DibuatPada:     d,
				}).Error; err != nil {
					t.Fatalf("Failed to create PergerakanStok at %v: %v", d, err)
				}
			}
			
			hutang := models.HutangPemasok{
				IDPemasok:     sup.ID,
				IDPO:          &po.ID,
				IDBarangMasuk: &inbound.ID,
				Jumlah:        po.Total,
				SisaHutang:    0,
				Status:        "paid",
				JumlahDibayar: po.Total,
				DibuatPada:    d,
			}
			if err := db.Create(&hutang).Error; err != nil {
				t.Fatalf("Failed to create HutangPemasok at %v: %v", d, err)
			}
		}
		
		// 8.b - Penjualan (Sales) 1-3x a day
		salesCount := rand.Intn(4)
		for s := 0; s < salesCount; s++ {
			p := produkList[rand.Intn(len(produkList))]
			
			var stockItem models.StokInventori
			if err := db.Where("id_produk = ? AND id_gudang = ?", p.ID, gudang.ID).First(&stockItem).Error; err != nil {
				continue // Usually err is record not found if not initialized, but we initialized it. Can be skipped.
			}
			
			qtySell := 1 + rand.Intn(5)
			if stockItem.Jumlah >= qtySell {
				totalHarga := p.HargaJual * float64(qtySell)
				sales := models.Penjualan{
					NomorTransaksi:   fmt.Sprintf("TRX-%04d", salesCounter),
					IDGudang:         gudang.ID,
					NamaPelanggan:    "Pelanggan Umum",
					Subtotal:         totalHarga,
					Total:            totalHarga,
					MetodePembayaran: "cash",
					JumlahPembayaran: totalHarga,
					Status:           "completed",
					IDKasir:          users[2].ID,
					DibuatPada:       d,
				}
				salesCounter++
				if err := db.Create(&sales).Error; err != nil {
					t.Fatalf("Failed to create Penjualan at %v: %v", d, err)
				}

				outbound := models.BarangKeluar{
					NomorTransaksi: fmt.Sprintf("OUT-%04d", outCounter),
					Alasan:         "penjualan",
					IDReferensi:    &sales.ID,
					TipeReferensi:  "sales",
					DibuatOleh:     users[2].ID,
					DibuatPada:     d,
				}
				outCounter++
				if err := db.Create(&outbound).Error; err != nil {
					t.Fatalf("Failed to create BarangKeluar at %v: %v", d, err)
				}
				
				var remainingToSell = qtySell
				var totalCOGS = 0.0
				var activeBatches []models.StokBatch
				
				db.Where("id_produk = ? AND id_gudang = ? AND jumlah_saat_ini > 0 AND aktif = true", p.ID, gudang.ID).Order("tanggal_masuk asc").Find(&activeBatches)
				
				for _, batch := range activeBatches {
					if remainingToSell <= 0 {
						break
					}
					takeQty := remainingToSell
					if batch.JumlahSaatIni < takeQty {
						takeQty = batch.JumlahSaatIni
					}
					
					batch.JumlahSaatIni -= takeQty
					if batch.JumlahSaatIni == 0 {
						batch.Aktif = false
					}
					if err := db.Save(&batch).Error; err != nil {
						t.Fatalf("Failed to update StokBatch at %v: %v", d, err)
					}
					
					cogsForBatch := float64(takeQty) * batch.HargaModal
					totalCOGS += cogsForBatch
					remainingToSell -= takeQty
				}
				
				if remainingToSell == 0 {
					avgCOGSPerUnit := totalCOGS / float64(qtySell)
					
					sItem := models.ItemPenjualan{
						IDPenjualan:  sales.ID,
						IDProduk:     p.ID,
						IDGudang:     gudang.ID,
						Jumlah:       qtySell,
						HargaSatuan:  p.HargaJual,
						HargaModal:   avgCOGSPerUnit,
						Subtotal:     totalHarga,
						TotalModal:   totalCOGS,
						DibuatPada:   d,
					}
					if err := db.Create(&sItem).Error; err != nil {
						t.Fatalf("Failed to create ItemPenjualan at %v: %v", d, err)
					}
					
					sales.TotalHargaModal = totalCOGS
					if err := db.Save(&sales).Error; err != nil {
						t.Fatalf("Failed to update Penjualan at %v: %v", d, err)
					}
					
					if err := db.Create(&models.ItemBarangKeluar{
						IDBarangKeluar: outbound.ID,
						IDProduk:       p.ID,
						Jumlah:         qtySell,
						IDGudang:       gudang.ID,
						DibuatPada:     d,
					}).Error; err != nil {
						t.Fatalf("Failed to create ItemBarangKeluar at %v: %v", d, err)
					}

					if err := db.Exec("UPDATE stok_inventori SET jumlah = jumlah - ?, pergerakan_terakhir_pada = ? WHERE id_produk = ? AND id_gudang = ?", qtySell, d, p.ID, gudang.ID).Error; err != nil {
						t.Fatalf("Failed to update StokInventori out at %v: %v", d, err)
					}
					
					if err := db.Create(&models.PergerakanStok{
						IDProduk:       p.ID,
						IDGudang:       gudang.ID,
						TipePergerakan: "out",
						TipeReferensi:  "sales",
						IDReferensi:    &sales.ID,
						Jumlah:         -qtySell,
						IDPengguna:     users[2].ID,
						DibuatPada:     d,
					}).Error; err != nil {
						t.Fatalf("Failed to create PergerakanStok out at %v: %v", d, err)
					}
				}
			}
		}

		// 8.c - Simulate Retur (Sales Return & Purchase Return) once a month
		if d.Day() == 15 {
			ret := models.ReturPenjualan{
				NomorRetur:         fmt.Sprintf("RET-J-%04d", returCounter),
				IDPenjualan:        1,
				NamaPelanggan:      "Pelanggan Umum",
				Alasan:             "Barang lecet pengiriman",
				Subtotal:           500000,
				Total:              500000,
				MetodePengembalian: "cash",
				JumlahPengembalian: 500000,
				Status:             "completed",
				DiprosesOleh:       users[2].ID,
				DiprosesPada:       d,
				DibuatPada:         d,
			}
			db.Create(&ret)
			returCounter++

			retB := models.ReturPembelian{
				NomorRetur:         fmt.Sprintf("RET-B-%04d", returCounter),
				IDPemasok:          pemasokList[0].ID,
				Alasan:             "Tidak sesuai spesifikasi",
				Subtotal:           100000,
				Total:              100000,
				MetodePengembalian: "potong_hutang",
				Status:             "completed",
				DibuatOleh:         users[1].ID,
				DibuatPada:         d,
			}
			db.Create(&retB)
			returCounter++
		}
	}

	log.Println("✅ 3-Month Database Seeding Completed Successfully!")
}
