package playlistrepository

import (
	"database/sql"
	"errors"

	playlistentity "mqfm_backend/internal/domain/entities/playlist"

)

//
// ============================================================
// 🔥 MYSQL PLAYLIST REPOSITORY
// ============================================================
type MySQLPlaylistRepository struct {
	db *sql.DB
}

//
// ============================================================
// 🔥 CONSTRUCTOR
// ============================================================
func NewMySQLPlaylistRepository(db *sql.DB) *MySQLPlaylistRepository {
	return &MySQLPlaylistRepository{db: db}
}

//
// ============================================================
// 🔥 SAVE (CREATE PLAYLIST)
// ============================================================
func (r *MySQLPlaylistRepository) Save(p *playlistentity.Playlist) error {
	query := `
		INSERT INTO playlists (
			user_id, name, description, thumbnail,
			created_at, updated_at
		)
		VALUES (?, ?, ?, ?, NOW(), NOW())
	`

	res, err := r.db.Exec(query,
		p.UserID,
		p.Name,
		p.Description,
		p.Thumbnail,
	)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err == nil {
		p.ID = int(id)
	}
	return nil
}

//
// ============================================================
// 🔥 FIND BY ID
// ============================================================
func (r *MySQLPlaylistRepository) FindByID(id int) (*playlistentity.Playlist, bool) {
	query := `
		SELECT id, user_id, name, description, thumbnail,
		       created_at, updated_at
		FROM playlists
		WHERE id = ? LIMIT 1
	`

	row := r.db.QueryRow(query, id)
	p := &playlistentity.Playlist{}

	err := row.Scan(
		&p.ID,
		&p.UserID,
		&p.Name,
		&p.Description,
		&p.Thumbnail,
		&p.CreatedAt,
		&p.UpdatedAt,
	)

	if err != nil {
		return nil, false
	}

	return p, true
}

//
// ============================================================
// 🔥 FIND BY USER ID (list playlist milik user)
// ============================================================
func (r *MySQLPlaylistRepository) FindByUserID(userID int) ([]*playlistentity.Playlist, error) {
	query := `
		SELECT id, user_id, name, description, thumbnail,
		       created_at, updated_at
		FROM playlists
		WHERE user_id = ?
		ORDER BY id DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*playlistentity.Playlist

	for rows.Next() {
		p := &playlistentity.Playlist{}
		_ = rows.Scan(
			&p.ID,
			&p.UserID,
			&p.Name,
			&p.Description,
			&p.Thumbnail,
			&p.CreatedAt,
			&p.UpdatedAt,
		)
		list = append(list, p)
	}

	return list, nil
}

//
// ============================================================
// 🔥 GET ALL PLAYLIST (GLOBAL)
// ============================================================
func (r *MySQLPlaylistRepository) GetAll() []*playlistentity.Playlist {
	query := `
		SELECT id, user_id, name, description, thumbnail,
		       created_at, updated_at
		FROM playlists
		ORDER BY id DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return []*playlistentity.Playlist{}
	}
	defer rows.Close()

	var list []*playlistentity.Playlist

	for rows.Next() {
		p := &playlistentity.Playlist{}
		_ = rows.Scan(
			&p.ID,
			&p.UserID,
			&p.Name,
			&p.Description,
			&p.Thumbnail,
			&p.CreatedAt,
			&p.UpdatedAt,
		)
		list = append(list, p)
	}

	return list
}

//
// ============================================================
// 🔥 UPDATE PLAYLIST
// ============================================================
func (r *MySQLPlaylistRepository) Update(p *playlistentity.Playlist) error {
	query := `
		UPDATE playlists
		SET name = ?, description = ?, thumbnail = ?,
		    updated_at = NOW()
		WHERE id = ?
	`

	_, err := r.db.Exec(query,
		p.Name,
		p.Description,
		p.Thumbnail,
		p.ID,
	)

	return err
}

//
// ============================================================
// 🔥 DELETE PLAYLIST
// ============================================================
func (r *MySQLPlaylistRepository) Delete(id int) error {
	query := `DELETE FROM playlists WHERE id = ?`
	res, err := r.db.Exec(query, id)

	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New("playlist tidak ditemukan")
	}

	return nil
}

//
// ============================================================
// 🔥 COUNT
// ============================================================
func (r *MySQLPlaylistRepository) Count() int {
	query := `SELECT COUNT(*) FROM playlists`
	var count int
	_ = r.db.QueryRow(query).Scan(&count)
	return count
}

//
// ============================================================
// 🔥 CLEAR (Testing Only)
// ============================================================
func (r *MySQLPlaylistRepository) Clear() {
	_, _ = r.db.Exec("DELETE FROM playlists")
}
