package domainrequest

import "mime/multipart"

// ---------------- CREATE AUDIO PODCAST REQUEST ----------------
type CreateAudioAdminPodcastRequest struct {
	Title       string                `form:"title" validate:"required"`
	Description string                `form:"description,omitempty"`
	AudioFile   *multipart.FileHeader `form:"audio_file" validate:"required"`  // file audio
	CategoryID  int                   `form:"category_id" validate:"required"` // id kategori
	Thumbnail   *multipart.FileHeader `form:"thumbnail,omitempty"`             // optional
}

// ---------------- UPDATE AUDIO PODCAST REQUEST ----------------
type UpdateAudioAdminPodcastRequest struct {
	ID          int                   `form:"id" validate:"required"`
	NewTitle    string                `form:"new_title,omitempty"`
	NewDesc     string                `form:"new_description,omitempty"`
	NewAudio    *multipart.FileHeader `form:"new_audio_file,omitempty"`
	NewCategory int                   `form:"new_category_id,omitempty"`
	NewThumb    *multipart.FileHeader `form:"new_thumbnail,omitempty"`
}

// ---------------- DELETE AUDIO PODCAST REQUEST ----------------
type DeleteAudioAdminPodcastRequest struct {
	ID int `json:"id" validate:"required"`
}

// ---------------- GET AUDIO PODCAST REQUEST ----------------
type GetAudioAdminPodcastRequest struct {
	ID int `json:"id" validate:"required"`
}
