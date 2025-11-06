package categoriesadminservice

import (
	"errors"
	"time"

	categoriesadminentity "mqfm_backend/internal/domain/entities/podcast/categories/admin"
	categoriesadminrepository "mqfm_backend/internal/infrastructure/repository/podcast/categories/admin"

)

// ---------------- SERVICE STRUCT ----------------
type CategoriesAdminService struct {
	repo categoriesadminrepository.CategoriesAdminRepository
}

// ---------------- CONSTRUCTOR ----------------
func NewCategoriesAdminService(repo categoriesadminrepository.CategoriesAdminRepository) *CategoriesAdminService {
	return &CategoriesAdminService{repo: repo}
}

// ---------------- CREATE ----------------
func (s *CategoriesAdminService) Create(name, description, image string) (*categoriesadminentity.CategoryAdmin, error) {
	if existing, found := s.repo.FindByName(name); found && existing != nil {
		return nil, errors.New("nama kategori sudah digunakan")
	}

	category := categoriesadminentity.NewCategoryAdmin(name, description, image)
	category.CreatedAt = time.Now()
	category.UpdatedAt = time.Now()

	if err := s.repo.Save(category); err != nil {
		return nil, err
	}

	return category, nil
}

// ---------------- GET ALL ----------------
func (s *CategoriesAdminService) GetAll() ([]*categoriesadminentity.CategoryAdmin, error) {
	categories := s.repo.GetAll()
	if len(categories) == 0 {
		return nil, errors.New("belum ada kategori")
	}
	return categories, nil
}

// ---------------- GET BY ID ----------------
func (s *CategoriesAdminService) GetByID(id int) (*categoriesadminentity.CategoryAdmin, error) {
	category, found := s.repo.FindByID(id)
	if !found {
		return nil, errors.New("kategori tidak ditemukan")
	}
	return category, nil
}

// ---------------- UPDATE ----------------
func (s *CategoriesAdminService) Update(id int, newName, newDescription, newImage string) (*categoriesadminentity.CategoryAdmin, error) {
	category, found := s.repo.FindByID(id)
	if !found {
		return nil, errors.New("kategori tidak ditemukan")
	}

	if newName != "" {
		category.Name = newName
	}
	if newDescription != "" {
		category.Description = newDescription
	}
	if newImage != "" {
		category.Image = newImage
	}

	category.UpdatedAt = time.Now()

	if err := s.repo.Update(category); err != nil {
		return nil, err
	}

	return category, nil
}

// ---------------- DELETE ----------------
func (s *CategoriesAdminService) Delete(id int) error {
	category, found := s.repo.FindByID(id)
	if !found {
		return errors.New("kategori tidak ditemukan")
	}
	return s.repo.Delete(category.ID)
}
