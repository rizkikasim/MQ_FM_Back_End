package audiouserentity

import "time"

// ----------------------------------------------------
// 🔥 ENTITY: Audio Podcast (Versi User / Read-Only)
// ----------------------------------------------------
type AudioUserPodcast struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	AudioURL    string    `json:"audio_url"`
	Duration    int       `json:"duration"`      // detik
	CategoryID  int       `json:"category_id"`
	Thumbnail   string    `json:"thumbnail"`     // cover
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
