package categoriesuserentity

import "time"

// ---------------- ENTITY ----------------
type CategoryUser struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Image       string    `json:"image"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	LastViewedAt time.Time `json:"last_viewed_at,omitempty"` // opsional untuk tracking interaksi user
}

// ---------------- FACTORY FUNCTION ----------------
// digunakan saat sistem menampilkan kategori ke user (bukan membuat baru)
func NewCategoryUser(id int, name, description, image string, createdAt, updatedAt time.Time) *CategoryUser {
	return &CategoryUser{
		ID:          id,
		Name:        name,
		Description: description,
		Image:       image,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}

// ---------------- DOMAIN METHODS ----------------

// TouchLastViewed memperbarui waktu terakhir kategori dilihat user
func (c *CategoryUser) TouchLastViewed() {
	c.LastViewedAt = time.Now()
}
