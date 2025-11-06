package categoriesadminrepository

import "mqfm_backend/internal/domain/entities/podcast/categories/admin"

// ---------------- INTERFACE ----------------
type CategoriesAdminRepository interface {
	// Simpan kategori baru
	Save(category *categoriesadminentity.CategoryAdmin) error

	// Ambil kategori berdasarkan ID
	FindByID(id int) (*categoriesadminentity.CategoryAdmin, bool)

		FindByName(name string) (*categoriesadminentity.CategoryAdmin, bool) // ✅ tambahkan ini


	// Ambil semua kategori
	GetAll() []*categoriesadminentity.CategoryAdmin

	// Update kategori (nama, deskripsi, gambar)
	Update(category *categoriesadminentity.CategoryAdmin) error

	// Hapus kategori berdasarkan ID
	Delete(id int) error

	// Cek apakah kategori sudah ada berdasarkan nama
	ExistsByName(name string) bool

	// Hitung jumlah total kategori
	Count() int

	// Hapus semua kategori (opsional, hanya untuk testing)
	Clear()
}
