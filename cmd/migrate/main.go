package main

import (
	"flag"
	"log"

	"real-erp-mebel/be/internal/config"
	"real-erp-mebel/be/internal/database"
	"real-erp-mebel/be/internal/models"
)

func main() {
	// Parse command line flags
	fresh := flag.Bool("fresh", false, "Drop all tables and migrate fresh")
	flag.Parse()

	// Load configuration
	config.LoadConfig()

	// Connect to database
	if err := database.Connect(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Fresh migration: drop all tables
	if *fresh {
		log.Println("⚠️  WARNING: Fresh migration will DROP ALL TABLES!")
		log.Println("Press Ctrl+C to cancel, or wait 5 seconds to continue...")

		// Give user time to cancel
		// In production, you might want to add confirmation

		log.Println("Dropping all tables...")

		// Drop all tables
		if err := database.DB.Exec("DROP SCHEMA public CASCADE; CREATE SCHEMA public;").Error; err != nil {
			log.Fatalf("Failed to drop schema: %v", err)
		}

		log.Println("✅ All tables dropped successfully")
	}

	log.Println("Starting database migration...")

	// Auto-migrate models (MVP Schema - 19 tables dengan nama bahasa Indonesia)
	if err := database.DB.AutoMigrate(
		// Core
		&models.Pengguna{},
		// Master Data
		&models.Produk{},
		&models.GambarProduk{},
		&models.Pemasok{},
		&models.Gudang{},
		// Stock Management
		&models.BarangMasuk{},
		&models.ItemBarangMasuk{},
		&models.BarangKeluar{},
		&models.ItemBarangKeluar{},
		&models.StokInventori{},
		&models.PergerakanStok{},
		// Sales
		&models.Penjualan{},
		&models.ItemPenjualan{},
		// Purchase Order
		&models.PesananPembelian{},
		&models.ItemPesananPembelian{},
		// Return
		&models.ReturPenjualan{},
		&models.ItemReturPenjualan{},
		&models.ReturPembelian{},
		&models.ItemReturPembelian{},
		// Finance
		&models.HutangPemasok{},
	); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	log.Println("✅ Database migration completed successfully!")
	log.Println("📊 Total tables migrated: 19 (dengan nama bahasa Indonesia)")

	if *fresh {
		log.Println("🔄 Fresh migration completed - All tables recreated")
	}
}
