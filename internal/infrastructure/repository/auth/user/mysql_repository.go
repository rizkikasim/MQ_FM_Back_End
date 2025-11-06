package authuserrepository

import (
	"database/sql"
	"errors"

	authuserentity "mqfm_backend/internal/domain/entities/auth/user"

)

// MySQLAuthUserRepository implementasi repository untuk MySQL
type MySQLAuthUserRepository struct {
	db *sql.DB
}

// NewMySQLAuthUserRepository constructor-nya
func NewMySQLAuthUserRepository(db *sql.DB) *MySQLAuthUserRepository {
	return &MySQLAuthUserRepository{db: db}
}

// ---------------- SAVE ----------------
func (r *MySQLAuthUserRepository) Save(user *authuserentity.User) error {
	query := `
		INSERT INTO users (
			email, username, phone, password,
			profile_initial, profile_color, token,
			created_at, updated_at
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`

	res, err := r.db.Exec(query,
		user.Email,
		user.Username,
		user.Phone,
		user.Password,
		user.ProfileTheme.Initial,
		user.ProfileTheme.PrimaryColor,
		user.Token,
	)
	if err != nil {
		return err
	}

	// ambil ID auto increment dari hasil insert
	id, err := res.LastInsertId()
	if err != nil {
		user.UserID = 0 // fallback aman
	} else {
		user.UserID = int(id)
	}

	return nil
}

// ---------------- FIND BY EMAIL ----------------
func (r *MySQLAuthUserRepository) FindByEmail(email string) (*authuserentity.User, bool) {
	query := `
		SELECT user_id, email, username, phone, password,
		       profile_initial, profile_color, token,
		       created_at, updated_at
		FROM users
		WHERE email = ? LIMIT 1
	`

	row := r.db.QueryRow(query, email)
	user := &authuserentity.User{}
	var initial, color sql.NullString

	err := row.Scan(
		&user.UserID,
		&user.Email,
		&user.Username,
		&user.Phone,
		&user.Password,
		&initial,
		&color,
		&user.Token,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, false
	}

	user.ProfileTheme.Initial = initial.String
	user.ProfileTheme.PrimaryColor = color.String
	return user, true
}

// ---------------- FIND BY IDENTIFIER ----------------
func (r *MySQLAuthUserRepository) FindByIdentifier(identifier string) (*authuserentity.User, bool) {
	query := `
		SELECT user_id, email, username, phone, password,
		       profile_initial, profile_color, token,
		       created_at, updated_at
		FROM users
		WHERE email = ? OR username = ? OR phone = ?
		LIMIT 1
	`

	row := r.db.QueryRow(query, identifier, identifier, identifier)
	user := &authuserentity.User{}
	var initial, color sql.NullString

	err := row.Scan(
		&user.UserID,
		&user.Email,
		&user.Username,
		&user.Phone,
		&user.Password,
		&initial,
		&color,
		&user.Token,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, false
	}

	user.ProfileTheme.Initial = initial.String
	user.ProfileTheme.PrimaryColor = color.String
	return user, true
}

// ---------------- UPDATE ----------------
func (r *MySQLAuthUserRepository) Update(user *authuserentity.User) error {
	query := `
		UPDATE users
		SET username = ?, phone = ?, password = ?,
		    profile_initial = ?, profile_color = ?,
		    token = ?, updated_at = NOW()
		WHERE email = ?
	`

	_, err := r.db.Exec(query,
		user.Username,
		user.Phone,
		user.Password,
		user.ProfileTheme.Initial,
		user.ProfileTheme.PrimaryColor,
		user.Token,
		user.Email,
	)
	return err
}

// ---------------- DELETE ----------------
func (r *MySQLAuthUserRepository) Delete(email string) error {
	query := `DELETE FROM users WHERE email = ?`
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
func (r *MySQLAuthUserRepository) GetAll() []*authuserentity.User {
	query := `
		SELECT user_id, email, username, phone,
		       profile_initial, profile_color,
		       created_at, updated_at
		FROM users
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return []*authuserentity.User{}
	}
	defer rows.Close()

	var users []*authuserentity.User
	for rows.Next() {
		u := &authuserentity.User{}
		var initial, color sql.NullString
		_ = rows.Scan(
			&u.UserID,
			&u.Email,
			&u.Username,
			&u.Phone,
			&initial,
			&color,
			&u.CreatedAt,
			&u.UpdatedAt,
		)
		u.ProfileTheme.Initial = initial.String
		u.ProfileTheme.PrimaryColor = color.String
		users = append(users, u)
	}

	return users
}

// ---------------- COUNT ----------------
func (r *MySQLAuthUserRepository) Count() int {
	query := `SELECT COUNT(*) FROM users`
	var count int
	_ = r.db.QueryRow(query).Scan(&count)
	return count
}

// ---------------- EXISTS ----------------
func (r *MySQLAuthUserRepository) Exists(email string) bool {
	query := `SELECT COUNT(*) FROM users WHERE email = ?`
	var count int
	_ = r.db.QueryRow(query, email).Scan(&count)
	return count > 0
}

// ---------------- CLEAR ----------------
func (r *MySQLAuthUserRepository) Clear() {
	_, _ = r.db.Exec("DELETE FROM users")
}
