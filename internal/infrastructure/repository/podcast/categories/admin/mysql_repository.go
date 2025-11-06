package categoriesadminrepository

import (
	"database/sql"
	"errors"

	categoriesadminentity "mqfm_backend/internal/domain/entities/podcast/categories/admin"

)

// ---------------- MySQL CATEGORIES ADMIN REPOSITORY ----------------
type MySQLCategoriesAdminRepository struct {
	db *sql.DB
}

// ---------------- CONSTRUCTOR ----------------
func NewMySQLCategoriesAdminRepository(db *sql.DB) *MySQLCategoriesAdminRepository {
	return &MySQLCategoriesAdminRepository{db: db}
}

// ---------------- SAVE ----------------
func (r *MySQLCategoriesAdminRepository) Save(category *categoriesadminentity.CategoryAdmin) error {
	query := `
		INSERT INTO categories_admin (
			name, description, image,
			created_at, updated_at
		)
		VALUES (?, ?, ?, NOW(), NOW())
	`

	res, err := r.db.Exec(query,
		category.Name,
		category.Description,
		category.Image,
	)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		category.ID = 0
	} else {
		category.ID = int(id)
	}
	return nil
}

// ---------------- FIND BY ID ----------------
func (r *MySQLCategoriesAdminRepository) FindByID(id int) (*categoriesadminentity.CategoryAdmin, bool) {
	query := `
		SELECT id, name, description, image,
		       created_at, updated_at
		FROM categories_admin
		WHERE id = ? LIMIT 1
	`

	row := r.db.QueryRow(query, id)
	c := &categoriesadminentity.CategoryAdmin{}

	err := row.Scan(
		&c.ID,
		&c.Name,
		&c.Description,
		&c.Image,
		&c.CreatedAt,
		&c.UpdatedAt,
	)
	if err != nil {
		return nil, false
	}
	return c, true
}

// ---------------- FIND BY NAME ----------------
func (r *MySQLCategoriesAdminRepository) FindByName(name string) (*categoriesadminentity.CategoryAdmin, bool) {
	query := `
		SELECT id, name, description, image,
		       created_at, updated_at
		FROM categories_admin
		WHERE name = ? LIMIT 1
	`

	row := r.db.QueryRow(query, name)
	c := &categoriesadminentity.CategoryAdmin{}

	err := row.Scan(
		&c.ID,
		&c.Name,
		&c.Description,
		&c.Image,
		&c.CreatedAt,
		&c.UpdatedAt,
	)
	if err != nil {
		return nil, false
	}
	return c, true
}

// ---------------- GET ALL ----------------
func (r *MySQLCategoriesAdminRepository) GetAll() []*categoriesadminentity.CategoryAdmin {
	query := `
		SELECT id, name, description, image,
		       created_at, updated_at
		FROM categories_admin
		ORDER BY id DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return []*categoriesadminentity.CategoryAdmin{}
	}
	defer rows.Close()

	var categories []*categoriesadminentity.CategoryAdmin
	for rows.Next() {
		c := &categoriesadminentity.CategoryAdmin{}
		_ = rows.Scan(
			&c.ID,
			&c.Name,
			&c.Description,
			&c.Image,
			&c.CreatedAt,
			&c.UpdatedAt,
		)
		categories = append(categories, c)
	}
	return categories
}

// ---------------- UPDATE ----------------
func (r *MySQLCategoriesAdminRepository) Update(category *categoriesadminentity.CategoryAdmin) error {
	query := `
		UPDATE categories_admin
		SET name = ?, description = ?, image = ?, updated_at = NOW()
		WHERE id = ?
	`

	_, err := r.db.Exec(query,
		category.Name,
		category.Description,
		category.Image,
		category.ID,
	)
	return err
}

// ---------------- DELETE ----------------
func (r *MySQLCategoriesAdminRepository) Delete(id int) error {
	query := `DELETE FROM categories_admin WHERE id = ?`
	res, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New("kategori tidak ditemukan")
	}
	return nil
}

// ---------------- EXISTS BY NAME ----------------
func (r *MySQLCategoriesAdminRepository) ExistsByName(name string) bool {
	query := `SELECT COUNT(*) FROM categories_admin WHERE name = ?`
	var count int
	_ = r.db.QueryRow(query, name).Scan(&count)
	return count > 0
}

// ---------------- COUNT ----------------
func (r *MySQLCategoriesAdminRepository) Count() int {
	query := `SELECT COUNT(*) FROM categories_admin`
	var count int
	_ = r.db.QueryRow(query).Scan(&count)
	return count
}

// ---------------- CLEAR ----------------
func (r *MySQLCategoriesAdminRepository) Clear() {
	_, _ = r.db.Exec("DELETE FROM categories_admin")
}
