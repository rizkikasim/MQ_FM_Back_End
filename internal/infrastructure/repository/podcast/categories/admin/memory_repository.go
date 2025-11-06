package categoriesadminrepository

import (
	"sync"

	categoriesadminentity "mqfm_backend/internal/domain/entities/podcast/categories/admin"

)

// ---------------- REPOSITORY (IN-MEMORY) ----------------
type InMemoryCategoriesAdminRepository struct {
	categories map[int]*categoriesadminentity.CategoryAdmin
	mu         sync.RWMutex
	nextID     int
}

// ---------------- CONSTRUCTOR ----------------
func NewInMemoryCategoriesAdminRepository() *InMemoryCategoriesAdminRepository {
	return &InMemoryCategoriesAdminRepository{
		categories: make(map[int]*categoriesadminentity.CategoryAdmin),
		nextID:     1,
	}
}

// ---------------- SAVE ----------------
func (r *InMemoryCategoriesAdminRepository) Save(category *categoriesadminentity.CategoryAdmin) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if category.ID == 0 {
		category.ID = r.nextID
		r.nextID++
	}
	r.categories[category.ID] = category
	return nil
}

// ---------------- FIND BY ID ----------------
func (r *InMemoryCategoriesAdminRepository) FindByID(id int) (*categoriesadminentity.CategoryAdmin, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	c, exists := r.categories[id]
	return c, exists
}

// ---------------- UPDATE ----------------
func (r *InMemoryCategoriesAdminRepository) Update(category *categoriesadminentity.CategoryAdmin) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.categories[category.ID]; !exists {
		return nil
	}
	r.categories[category.ID] = category
	return nil
}

// ---------------- DELETE ----------------
func (r *InMemoryCategoriesAdminRepository) Delete(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.categories, id)
	return nil
}

// ---------------- GET ALL ----------------
func (r *InMemoryCategoriesAdminRepository) GetAll() []*categoriesadminentity.CategoryAdmin {
	r.mu.RLock()
	defer r.mu.RUnlock()
	categoryList := make([]*categoriesadminentity.CategoryAdmin, 0, len(r.categories))
	for _, c := range r.categories {
		categoryList = append(categoryList, c)
	}
	return categoryList
}

// ---------------- EXISTS BY NAME ----------------
func (r *InMemoryCategoriesAdminRepository) ExistsByName(name string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, c := range r.categories {
		if c.Name == name {
			return true
		}
	}
	return false
}

// ---------------- COUNT ----------------
func (r *InMemoryCategoriesAdminRepository) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.categories)
}

// ---------------- CLEAR ----------------
func (r *InMemoryCategoriesAdminRepository) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.categories = make(map[int]*categoriesadminentity.CategoryAdmin)
	r.nextID = 1
}
