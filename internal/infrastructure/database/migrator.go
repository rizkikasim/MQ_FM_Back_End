package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"

)

// RunMigrations menjalankan otomatis migration SQL di folder /internal/migrations
func RunMigrations() {
	_ = godotenv.Load("env_staging.env")

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")

	// pastikan path absolut supaya bisa ditemukan
	root, _ := os.Getwd()
	migrationsPath := fmt.Sprintf("file://%s/internal/migrations", root)

	// encode password biar aman dari karakter khusus
	dsn := fmt.Sprintf("mysql://%s:%s@tcp(%s)/%s?multiStatements=true",
		dbUser, dbPass, dbHost, dbName)

	m, err := migrate.New(migrationsPath, dsn)
	if err != nil {
		log.Fatalf("❌ Error inisialisasi migrasi: %v", err)
	}

	if err := m.Up(); err != nil && err.Error() != "no change" {
		log.Fatalf("❌ Gagal menjalankan migrasi: %v", err)
	}

	log.Println("✅ Migrasi database berhasil atau sudah up-to-date")
}


// RefreshMigrations menghapus semua tabel & migrasi ulang dari awal
func RefreshMigrations() {
	_ = godotenv.Load("env_staging.env")

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", dbUser, dbPass, dbHost, dbName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("❌ Gagal konek DB untuk refresh: %v", err)
	}
	defer db.Close()

	_, _ = db.Exec("DROP TABLE IF EXISTS admins, schema_migrations")

	log.Println("🗑️  Semua tabel dihapus, menjalankan migrasi ulang...")

	RunMigrations()
}

