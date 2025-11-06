package categoriesuserrepository

import (
	"sync"

	categoriesuserentity "mqfm_backend/internal/domain/entities/podcast/categories/user"

)

// ---------------- REPOSITORY (IN-MEMORY) ----------------
type InMemoryCategoriesUserRepository struct {
	categories map[int]*categoriesuserentity.CategoryUser
	mu         sync.RWMutex
}

// ---------------- CONSTRUCTOR ----------------
func NewInMemoryCategoriesUserRepository() *InMemoryCategoriesUserRepository {
	return &InMemoryCategoriesUserRepository{
		categories: make(map[int]*categoriesuserentity.CategoryUser),
	}
}

// ---------------- FIND BY ID ----------------
func (r *InMemoryCategoriesUserRepository) FindByID(id int) (*categoriesuserentity.CategoryUser, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	c, exists := r.categories[id]
	return c, exists
}

// ---------------- FIND BY NAME ----------------
func (r *InMemoryCategoriesUserRepository) FindByName(name string) (*categoriesuserentity.CategoryUser, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, c := range r.categories {
		if c.Name == name {
			return c, true
		}
	}
	return nil, false
}

// ---------------- GET ALL ----------------
func (r *InMemoryCategoriesUserRepository) GetAll() []*categoriesuserentity.CategoryUser {
	r.mu.RLock()
	defer r.mu.RUnlock()
	categoryList := make([]*categoriesuserentity.CategoryUser, 0, len(r.categories))
	for _, c := range r.categories {
		categoryList = append(categoryList, c)
	}
	return categoryList
}

// ---------------- UPDATE ----------------
func (r *InMemoryCategoriesUserRepository) Update(category *categoriesuserentity.CategoryUser) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.categories[category.ID]; !exists {
		return nil
	}
	r.categories[category.ID] = category
	return nil
}
