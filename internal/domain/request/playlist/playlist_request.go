package domainrequest

// =====================================================
// 🔥 CREATE PLAYLIST REQUEST
// =====================================================
type CreatePlaylistRequest struct {
	Title       string `form:"title" validate:"required"`
	Description string `form:"description,omitempty"`
	Thumbnail   string `form:"thumbnail,omitempty"` // URL string (bukan file upload)
	UserID      int    `form:"user_id" validate:"required"`
}

// =====================================================
// 🔥 UPDATE PLAYLIST REQUEST
// =====================================================
type UpdatePlaylistRequest struct {
	ID          int    `form:"id" validate:"required"`
	NewTitle    string `form:"new_title,omitempty"`
	NewDesc     string `form:"new_description,omitempty"`
	NewThumb    string `form:"new_thumbnail,omitempty"`
}

// =====================================================
// 🔥 DELETE PLAYLIST REQUEST
// =====================================================
type DeletePlaylistRequest struct {
	ID int `json:"id" validate:"required"`
}

// =====================================================
// 🔥 GET PLAYLIST BY ID REQUEST
// =====================================================
type GetPlaylistRequest struct {
	ID int `json:"id" validate:"required"`
}
