package categoriesuserrepository

import (
	"database/sql"

	categoriesuserentity "mqfm_backend/internal/domain/entities/podcast/categories/user"

)

// ---------------- MySQL CATEGORIES USER REPOSITORY ----------------
type MySQLCategoriesUserRepository struct {
	db *sql.DB
}

// ---------------- CONSTRUCTOR ----------------
func NewMySQLCategoriesUserRepository(db *sql.DB) *MySQLCategoriesUserRepository {
	return &MySQLCategoriesUserRepository{db: db}
}

// ---------------- GET ALL ----------------
func (r *MySQLCategoriesUserRepository) GetAll() []*categoriesuserentity.CategoryUser {
	query := `
		SELECT id, name, description, image, created_at, updated_at
		FROM categories_admin
		ORDER BY id DESC
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return []*categoriesuserentity.CategoryUser{}
	}
	defer rows.Close()

	var categories []*categoriesuserentity.CategoryUser
	for rows.Next() {
		c := &categoriesuserentity.CategoryUser{}
		_ = rows.Scan(&c.ID, &c.Name, &c.Description, &c.Image, &c.CreatedAt, &c.UpdatedAt)
		categories = append(categories, c)
	}
	return categories
}

// ---------------- FIND BY ID ----------------
func (r *MySQLCategoriesUserRepository) FindByID(id int) (*categoriesuserentity.CategoryUser, bool) {
	query := `
		SELECT id, name, description, image, created_at, updated_at
		FROM categories_admin
		WHERE id = ? LIMIT 1
	`
	row := r.db.QueryRow(query, id)
	c := &categoriesuserentity.CategoryUser{}

	err := row.Scan(&c.ID, &c.Name, &c.Description, &c.Image, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, false
	}
	return c, true
}

// ---------------- FIND BY NAME ----------------
func (r *MySQLCategoriesUserRepository) FindByName(name string) (*categoriesuserentity.CategoryUser, bool) {
	query := `
		SELECT id, name, description, image, created_at, updated_at
		FROM categories_admin
		WHERE name = ? LIMIT 1
	`
	row := r.db.QueryRow(query, name)
	c := &categoriesuserentity.CategoryUser{}

	err := row.Scan(&c.ID, &c.Name, &c.Description, &c.Image, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, false
	}
	return c, true
}

// ---------------- UPDATE (dummy) ----------------
// digunakan hanya agar repository memenuhi interface CategoriesUserRepository
func (r *MySQLCategoriesUserRepository) Update(category *categoriesuserentity.CategoryUser) error {
	// user tidak boleh mengubah kategori, jadi tidak ada aksi update
	return nil
}
