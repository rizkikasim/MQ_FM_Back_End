package playlistentity

import "time"

// ======================================================
// 🔥 ENTITY: Playlist (Khusus User)
// ======================================================
type Playlist struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Thumbnail   string    `json:"thumbnail"`
	UserID      int       `json:"user_id"` // playlist milik user tertentu
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ======================================================
// 🔥 FACTORY FUNCTION
// ======================================================
func NewPlaylist(
	title string,
	description string,
	thumbnail string,
	userID int,
) *Playlist {

	return &Playlist{
		Title:       title,
		Description: description,
		Thumbnail:   thumbnail,
		UserID:      userID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// ======================================================
// 🔥 DOMAIN METHOD: Update Playlist
// ======================================================
func (p *Playlist) UpdatePlaylist(
	newTitle string,
	newDescription string,
	newThumbnail string,
) {
	if newTitle != "" {
		p.Title = newTitle
	}
	if newDescription != "" {
		p.Description = newDescription
	}
	if newThumbnail != "" {
		p.Thumbnail = newThumbnail
	}

	p.UpdatedAt = time.Now()
}
