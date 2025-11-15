package audioadminrepository

import (
	"database/sql"
	"errors"

	audioentities "mqfm_backend/internal/domain/entities/podcast/audio/admin"

)

// ------------------------------------------------------------
// 🔥 MYSQL AUDIO PODCAST ADMIN REPOSITORY
// ------------------------------------------------------------
type MySQLAudioAdminPodcastRepository struct {
	db *sql.DB
}

// ------------------------------------------------------------
// 🔥 CONSTRUCTOR
// ------------------------------------------------------------
func NewMySQLAudioAdminPodcastRepository(db *sql.DB) *MySQLAudioAdminPodcastRepository {
	return &MySQLAudioAdminPodcastRepository{db: db}
}

// ------------------------------------------------------------
// 🔥 SAVE (CREATE)
// ------------------------------------------------------------
func (r *MySQLAudioAdminPodcastRepository) Save(p *audioentities.AudioAdminPodcast) error {
	query := `
		INSERT INTO audio_admin_podcast (
			title, description, audio_url, duration,
			category_id, thumbnail,
			created_at, updated_at
		)
		VALUES (?, ?, ?, ?, ?, ?, NOW(), NOW())
	`

	res, err := r.db.Exec(query,
		p.Title,
		p.Description,
		p.AudioURL,
		p.Duration,
		p.CategoryID,
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

// ------------------------------------------------------------
// 🔥 FIND BY ID
// ------------------------------------------------------------
func (r *MySQLAudioAdminPodcastRepository) FindByID(id int) (*audioentities.AudioAdminPodcast, bool) {
	query := `
		SELECT id, title, description, audio_url, duration,
			   category_id, thumbnail,
			   created_at, updated_at
		FROM audio_admin_podcast
		WHERE id = ? LIMIT 1
	`

	row := r.db.QueryRow(query, id)
	p := &audioentities.AudioAdminPodcast{}

	err := row.Scan(
		&p.ID,
		&p.Title,
		&p.Description,
		&p.AudioURL,
		&p.Duration,
		&p.CategoryID,
		&p.Thumbnail,
		&p.CreatedAt,
		&p.UpdatedAt,
	)

	if err != nil {
		return nil, false
	}

	return p, true
}

// ------------------------------------------------------------
// 🔥 FIND BY TITLE
// ------------------------------------------------------------
func (r *MySQLAudioAdminPodcastRepository) FindByTitle(title string) (*audioentities.AudioAdminPodcast, bool) {
	query := `
		SELECT id, title, description, audio_url, duration,
			   category_id, thumbnail,
			   created_at, updated_at
		FROM audio_admin_podcast
		WHERE title = ? LIMIT 1
	`

	row := r.db.QueryRow(query, title)
	p := &audioentities.AudioAdminPodcast{}

	err := row.Scan(
		&p.ID,
		&p.Title,
		&p.Description,
		&p.AudioURL,
		&p.Duration,
		&p.CategoryID,
		&p.Thumbnail,
		&p.CreatedAt,
		&p.UpdatedAt,
	)

	if err != nil {
		return nil, false
	}

	return p, true
}

// ------------------------------------------------------------
// 🔥 GET ALL
// ------------------------------------------------------------
func (r *MySQLAudioAdminPodcastRepository) GetAll() []*audioentities.AudioAdminPodcast {
	query := `
		SELECT id, title, description, audio_url, duration,
			   category_id, thumbnail,
			   created_at, updated_at
		FROM audio_admin_podcast
		ORDER BY id DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return []*audioentities.AudioAdminPodcast{}
	}
	defer rows.Close()

	var list []*audioentities.AudioAdminPodcast

	for rows.Next() {
		p := &audioentities.AudioAdminPodcast{}
		_ = rows.Scan(
			&p.ID,
			&p.Title,
			&p.Description,
			&p.AudioURL,
			&p.Duration,
			&p.CategoryID,
			&p.Thumbnail,
			&p.CreatedAt,
			&p.UpdatedAt,
		)
		list = append(list, p)
	}

	return list
}

// ------------------------------------------------------------
// 🔥 UPDATE
// ------------------------------------------------------------
func (r *MySQLAudioAdminPodcastRepository) Update(p *audioentities.AudioAdminPodcast) error {
	query := `
		UPDATE audio_admin_podcast
		SET title = ?, description = ?, audio_url = ?,
			duration = ?, category_id = ?, thumbnail = ?,
			updated_at = NOW()
		WHERE id = ?
	`

	_, err := r.db.Exec(query,
		p.Title,
		p.Description,
		p.AudioURL,
		p.Duration,
		p.CategoryID,
		p.Thumbnail,
		p.ID,
	)

	return err
}

// ------------------------------------------------------------
// 🔥 DELETE
// ------------------------------------------------------------
func (r *MySQLAudioAdminPodcastRepository) Delete(id int) error {
	query := `DELETE FROM audio_admin_podcast WHERE id = ?`
	res, err := r.db.Exec(query, id)

	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New("podcast audio tidak ditemukan")
	}

	return nil
}

// ------------------------------------------------------------
// 🔥 COUNT
// ------------------------------------------------------------
func (r *MySQLAudioAdminPodcastRepository) Count() int {
	query := `SELECT COUNT(*) FROM audio_admin_podcast`
	var count int
	_ = r.db.QueryRow(query).Scan(&count)
	return count
}

// ------------------------------------------------------------
// 🔥 CLEAR (Testing)
// ------------------------------------------------------------
func (r *MySQLAudioAdminPodcastRepository) Clear() {
	_, _ = r.db.Exec("DELETE FROM audio_admin_podcast")
}
