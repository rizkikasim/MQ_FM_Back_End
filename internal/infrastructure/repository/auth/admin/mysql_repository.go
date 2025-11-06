package authadminrepository

import (
	"database/sql"
	"errors"

	authadminentity "mqfm_backend/internal/domain/entities/auth/admin"

)

// MySQLAuthAdminRepository implementasi repository untuk MySQL
type MySQLAuthAdminRepository struct {
	db *sql.DB
}

// NewMySQLAuthAdminRepository constructor-nya
func NewMySQLAuthAdminRepository(db *sql.DB) *MySQLAuthAdminRepository {
	return &MySQLAuthAdminRepository{db: db}
}

// ---------------- SAVE ----------------
func (r *MySQLAuthAdminRepository) Save(admin *authadminentity.Admin) error {
	query := `
		INSERT INTO admins (
			email, username, phone, password,
			profile_initial, profile_color, token,
			created_at, updated_at
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`

	res, err := r.db.Exec(query,
		admin.Email,
		admin.Username,
		admin.Phone,
		admin.Password,
		admin.ProfileTheme.Initial,
		admin.ProfileTheme.PrimaryColor,
		admin.Token,
	)
	if err != nil {
		return err
	}

	// ambil ID auto increment dari hasil insert
	id, err := res.LastInsertId()
	if err != nil {
		admin.AdminID = 0 // fallback aman
	} else {
		admin.AdminID = int(id)
	}

	return nil
}

// ---------------- FIND BY EMAIL ----------------
func (r *MySQLAuthAdminRepository) FindByEmail(email string) (*authadminentity.Admin, bool) {
	query := `
		SELECT admin_id, email, username, phone, password,
		       profile_initial, profile_color, token,
		       created_at, updated_at
		FROM admins
		WHERE email = ? LIMIT 1
	`

	row := r.db.QueryRow(query, email)
	admin := &authadminentity.Admin{}
	var initial, color sql.NullString

	err := row.Scan(
		&admin.AdminID,
		&admin.Email,
		&admin.Username,
		&admin.Phone,
		&admin.Password,
		&initial,
		&color,
		&admin.Token,
		&admin.CreatedAt,
		&admin.UpdatedAt,
	)
	if err != nil {
		return nil, false
	}

	admin.ProfileTheme.Initial = initial.String
	admin.ProfileTheme.PrimaryColor = color.String
	return admin, true
}

// ---------------- FIND BY IDENTIFIER ----------------
func (r *MySQLAuthAdminRepository) FindByIdentifier(identifier string) (*authadminentity.Admin, bool) {
	query := `
		SELECT admin_id, email, username, phone, password,
		       profile_initial, profile_color, token,
		       created_at, updated_at
		FROM admins
		WHERE email = ? OR username = ? OR phone = ?
		LIMIT 1
	`

	row := r.db.QueryRow(query, identifier, identifier, identifier)
	admin := &authadminentity.Admin{}
	var initial, color sql.NullString

	err := row.Scan(
		&admin.AdminID,
		&admin.Email,
		&admin.Username,
		&admin.Phone,
		&admin.Password,
		&initial,
		&color,
		&admin.Token,
		&admin.CreatedAt,
		&admin.UpdatedAt,
	)
	if err != nil {
		return nil, false
	}

	admin.ProfileTheme.Initial = initial.String
	admin.ProfileTheme.PrimaryColor = color.String
	return admin, true
}

// ---------------- UPDATE ----------------
func (r *MySQLAuthAdminRepository) Update(admin *authadminentity.Admin) error {
	query := `
		UPDATE admins
		SET username = ?, phone = ?, password = ?,
		    profile_initial = ?, profile_color = ?,
		    token = ?, updated_at = NOW()
		WHERE email = ?
	`

	_, err := r.db.Exec(query,
		admin.Username,
		admin.Phone,
		admin.Password,
		admin.ProfileTheme.Initial,
		admin.ProfileTheme.PrimaryColor,
		admin.Token,
		admin.Email,
	)
	return err
}

// ---------------- DELETE ----------------
func (r *MySQLAuthAdminRepository) Delete(email string) error {
	query := `DELETE FROM admins WHERE email = ?`
	res, err := r.db.Exec(query, email)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New("akun tidak ditemukan")
	}
	return nil
}

// ---------------- GET ALL ----------------
func (r *MySQLAuthAdminRepository) GetAll() []*authadminentity.Admin {
	query := `
		SELECT admin_id, email, username, phone,
		       profile_initial, profile_color,
		       created_at, updated_at
		FROM admins
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return []*authadminentity.Admin{}
	}
	defer rows.Close()

	var admins []*authadminentity.Admin
	for rows.Next() {
		a := &authadminentity.Admin{}
		var initial, color sql.NullString
		_ = rows.Scan(
			&a.AdminID,
			&a.Email,
			&a.Username,
			&a.Phone,
			&initial,
			&color,
			&a.CreatedAt,
			&a.UpdatedAt,
		)
		a.ProfileTheme.Initial = initial.String
		a.ProfileTheme.PrimaryColor = color.String
		admins = append(admins, a)
	}

	return admins
}

// ---------------- COUNT ----------------
func (r *MySQLAuthAdminRepository) Count() int {
	query := `SELECT COUNT(*) FROM admins`
	var count int
	_ = r.db.QueryRow(query).Scan(&count)
	return count
}

// ---------------- EXISTS ----------------
func (r *MySQLAuthAdminRepository) Exists(email string) bool {
	query := `SELECT COUNT(*) FROM admins WHERE email = ?`
	var count int
	_ = r.db.QueryRow(query, email).Scan(&count)
	return count > 0
}

// ---------------- CLEAR ----------------
func (r *MySQLAuthAdminRepository) Clear() {
	_, _ = r.db.Exec("DELETE FROM admins")
}
