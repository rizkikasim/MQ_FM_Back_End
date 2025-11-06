package categoriesuserservice

import (
	"errors"
	"time"

	categoriesuserentity "mqfm_backend/internal/domain/entities/podcast/categories/user"
	categoriesuserrepository "mqfm_backend/internal/infrastructure/repository/podcast/categories/user"

)

// ---------------- SERVICE STRUCT ----------------
type CategoriesUserService struct {
	repo categoriesuserrepository.CategoriesUserRepository
}

// ---------------- CONSTRUCTOR ----------------
func NewCategoriesUserService(repo categoriesuserrepository.CategoriesUserRepository) *CategoriesUserService {
	return &CategoriesUserService{repo: repo}
}

// ---------------- GET ALL ----------------
func (s *CategoriesUserService) GetAll() ([]*categoriesuserentity.CategoryUser, error) {
	categories := s.repo.GetAll()
	if len(categories) == 0 {
		return nil, errors.New("belum ada kategori yang tersedia")
	}
	return categories, nil
}

// ---------------- GET BY ID ----------------
func (s *CategoriesUserService) GetByID(id int) (*categoriesuserentity.CategoryUser, error) {
	category, found := s.repo.FindByID(id)
	if !found {
		return nil, errors.New("kategori tidak ditemukan")
	}
	return category, nil
}

// ---------------- SEARCH BY NAME ----------------
func (s *CategoriesUserService) FindByName(name string) (*categoriesuserentity.CategoryUser, error) {
	category, found := s.repo.FindByName(name)
	if !found {
		return nil, errors.New("kategori tidak ditemukan")
	}
	return category, nil
}

// ---------------- UPDATE LAST VIEWED ----------------
// (opsional untuk tracking interaksi user)
func (s *CategoriesUserService) UpdateLastViewed(id int) error {
	category, found := s.repo.FindByID(id)
	if !found {
		return errors.New("kategori tidak ditemukan")
	}

	category.LastViewedAt = time.Now()
	return s.repo.Update(category)
}
