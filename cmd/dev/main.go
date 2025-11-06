package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"mqfm_backend/internal/infrastructure/database"
	"mqfm_backend/internal/presentation/router"

)

func main() {
	// 🔹 load env file
	if err := godotenv.Load("env_staging.env"); err != nil {
		log.Fatalf("❌ gagal load env: %v", err)
	}

	// 🔹 cek apakah perlu refresh tabel
	if os.Getenv("DB_REFRESH") == "true" {
		log.Println("🗑️  Mode refresh aktif — drop & re-run semua tabel...")
		database.RefreshMigrations()
	} else {
		// 🔹 jalankan migrasi biasa
		database.RunMigrations()
	}

	// 🔹 setup router
	mux := router.SetupAppRouter()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	log.Printf("✅ MQFM backend running at :%s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("💀 Server mati: %v", err)
	}
}
