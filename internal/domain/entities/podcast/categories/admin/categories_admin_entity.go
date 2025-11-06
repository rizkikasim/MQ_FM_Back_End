package categoriesadminentity

import "time"

// ---------------- ENTITY ----------------
type CategoryAdmin struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Image       string    `json:"image"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ---------------- FACTORY FUNCTION ----------------
// digunakan saat admin membuat kategori baru
func NewCategoryAdmin(name, description, image string) *CategoryAdmin {
	return &CategoryAdmin{
		Name:        name,
		Description: description,
		Image:       image,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// ---------------- DOMAIN METHODS ----------------

// UpdateCategory memperbarui data kategori (nama, deskripsi, atau gambar)
func (c *CategoryAdmin) UpdateCategory(newName, newDescription, newImage string) {
	if newName != "" {
		c.Name = newName
	}
	if newDescription != "" {
		c.Description = newDescription
	}
	if newImage != "" {
		c.Image = newImage
	}
	c.UpdatedAt = time.Now()
}
