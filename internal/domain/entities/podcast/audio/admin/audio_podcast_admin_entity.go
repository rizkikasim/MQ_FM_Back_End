package audioadminentity

import "time"

// ---------------- ENTITY ----------------
type AudioAdminPodcast struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	AudioURL    string    `json:"audio_url"`
	Duration    int       `json:"duration"`      // dalam detik
	CategoryID  int       `json:"category_id"`
	Thumbnail   string    `json:"thumbnail"`     // url gambar
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ---------------- FACTORY FUNCTION ----------------
// digunakan saat admin mengupload podcast audio baru
func NewAudioAdminPodcast(
	title string,
	description string,
	audioURL string,
	duration int,
	categoryID int,
	thumbnail string,
) *AudioAdminPodcast {
	return &AudioAdminPodcast{
		Title:       title,
		Description: description,
		AudioURL:    audioURL,
		Duration:    duration,
		CategoryID:  categoryID,
		Thumbnail:   thumbnail,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// ---------------- DOMAIN METHODS ----------------

// UpdatePodcast memperbarui data podcast audio
func (p *AudioAdminPodcast) UpdatePodcast(
	newTitle string,
	newDescription string,
	newAudioURL string,
	newDuration int,
	newCategoryID int,
	newThumbnail string,
) {
	if newTitle != "" {
		p.Title = newTitle
	}
	if newDescription != "" {
		p.Description = newDescription
	}
	if newAudioURL != "" {
		p.AudioURL = newAudioURL
	}
	if newDuration > 0 {
		p.Duration = newDuration
	}
	if newCategoryID > 0 {
		p.CategoryID = newCategoryID
	}
	if newThumbnail != "" {
		p.Thumbnail = newThumbnail
	}

	p.UpdatedAt = time.Now()
}
