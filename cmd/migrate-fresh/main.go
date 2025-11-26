package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"real-erp-mebel/be/internal/config"
	"real-erp-mebel/be/internal/database"
	"real-erp-mebel/be/internal/models"
)

func main() {
	// Load configuration
	config.LoadConfig()

	// Connect to database
	if err := database.Connect(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Warning message
	fmt.Println("⚠️  ⚠️  ⚠️  WARNING: FRESH MIGRATION ⚠️  ⚠️  ⚠️")
	fmt.Println("")
	fmt.Println("This will:")
	fmt.Println("  - DROP ALL TABLES in the database")
	fmt.Println("  - DELETE ALL DATA")
	fmt.Println("  - Recreate all tables from scratch")
	fmt.Println("")
	fmt.Println("This action CANNOT be undone!")
	fmt.Println("")

	// Ask for confirmation
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Type 'yes' to continue, or anything else to cancel: ")
	confirmation, _ := reader.ReadString('\n')
	confirmation = strings.TrimSpace(strings.ToLower(confirmation))

	if confirmation != "yes" {
		fmt.Println("❌ Migration cancelled")
		os.Exit(0)
	}

	fmt.Println("")
	log.Println("Dropping all tables...")

	// Drop all tables using GORM Migrator
	migrator := database.DB.Migrator()

	// Get all tables
	tables, err := migrator.GetTables()
	if err != nil {
		log.Fatalf("Failed to get tables: %v", err)
	}

	// Drop all tables
	for _, table := range tables {
		log.Printf("Dropping table: %s", table)
		if err := migrator.DropTable(table); err != nil {
			log.Printf("Warning: Failed to drop table %s: %v", table, err)
		}
	}

	// Small delay to ensure tables are dropped
	time.Sleep(500 * time.Millisecond)

	log.Println("✅ All tables dropped successfully")
	log.Println("")
	log.Println("Starting fresh migration...")

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

	log.Println("")
	log.Println("✅ Fresh migration completed successfully!")
	log.Println("📊 Total tables migrated: 19 (dengan nama bahasa Indonesia)")
	log.Println("🔄 All tables have been recreated from scratch")
}
