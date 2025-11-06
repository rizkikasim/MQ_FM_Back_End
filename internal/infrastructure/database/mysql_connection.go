package database

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"

)

// NewMySQLConnection membuat koneksi ke database MySQL
func NewMySQLConnection(user, pass, host, name string) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", user, pass, host, name)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
