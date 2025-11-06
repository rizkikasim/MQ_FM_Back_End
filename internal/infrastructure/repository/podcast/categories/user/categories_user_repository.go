package categoriesuserrepository

import categoriesuserentity "mqfm_backend/internal/domain/entities/podcast/categories/user"

// ---------------- INTERFACE ----------------
type CategoriesUserRepository interface {
	// Ambil semua kategori yang dapat dilihat user
	GetAll() []*categoriesuserentity.CategoryUser

	// Ambil kategori berdasarkan ID
	FindByID(id int) (*categoriesuserentity.CategoryUser, bool)

	// Cari kategori berdasarkan nama
	FindByName(name string) (*categoriesuserentity.CategoryUser, bool)

	// Update data kategori (misal waktu terakhir dilihat user)
	Update(category *categoriesuserentity.CategoryUser) error
}
