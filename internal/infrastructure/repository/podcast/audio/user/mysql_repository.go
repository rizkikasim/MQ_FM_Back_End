package audiouserrepository

import (
	"database/sql"

	audioentities "mqfm_backend/internal/domain/entities/podcast/audio/user"

)

// ------------------------------------------------------------
// 🔥 MYSQL AUDIO PODCAST USER REPOSITORY (READ ONLY)
// ------------------------------------------------------------
type MySQLAudioUserPodcastRepository struct {
	db *sql.DB
}

// ------------------------------------------------------------
// 🔥 CONSTRUCTOR
// ------------------------------------------------------------
func NewMySQLAudioUserPodcastRepository(db *sql.DB) *MySQLAudioUserPodcastRepository {
	return &MySQLAudioUserPodcastRepository{db: db}
}

// ------------------------------------------------------------
// 🔥 FIND BY ID
// ------------------------------------------------------------
func (r *MySQLAudioUserPodcastRepository) FindByID(id int) (*audioentities.AudioUserPodcast, bool) {
	query := `
		SELECT id, title, description, audio_url, duration,
			   category_id, thumbnail,
			   created_at, updated_at
		FROM audio_admin_podcast
		WHERE id = ? LIMIT 1
	`

	row := r.db.QueryRow(query, id)
	p := &audioentities.AudioUserPodcast{}

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
// 🔥 GET ALL (semua podcast yang available untuk user)
// ------------------------------------------------------------
func (r *MySQLAudioUserPodcastRepository) GetAll() []*audioentities.AudioUserPodcast {
	query := `
		SELECT id, title, description, audio_url, duration,
			   category_id, thumbnail,
			   created_at, updated_at
		FROM audio_admin_podcast
		ORDER BY id DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return []*audioentities.AudioUserPodcast{}
	}
	defer rows.Close()

	var list []*audioentities.AudioUserPodcast

	for rows.Next() {
		p := &audioentities.AudioUserPodcast{}
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
