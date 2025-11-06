package authadminrepository

import (
	"sync"

	authadminentity "mqfm_backend/internal/domain/entities/auth/admin"

)

// ---------------- REPOSITORY (IN-MEMORY) ----------------
type InMemoryAuthAdminRepository struct {
	admins map[string]*authadminentity.Admin
	mu     sync.RWMutex
	nextID int
}

// ---------------- CONSTRUCTOR ----------------
func NewInMemoryAuthAdminRepository() *InMemoryAuthAdminRepository {
	return &InMemoryAuthAdminRepository{
		admins: make(map[string]*authadminentity.Admin),
		nextID: 1,
	}
}

// ---------------- SAVE ----------------
func (r *InMemoryAuthAdminRepository) Save(admin *authadminentity.Admin) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if admin.AdminID == 0 {
		admin.AdminID = r.nextID
		r.nextID++
	}
	r.admins[admin.Email] = admin
	return nil
}

// ---------------- FIND BY EMAIL ----------------
func (r *InMemoryAuthAdminRepository) FindByEmail(email string) (*authadminentity.Admin, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	admin, exists := r.admins[email]
	return admin, exists
}

// ---------------- FIND BY IDENTIFIER ----------------
func (r *InMemoryAuthAdminRepository) FindByIdentifier(identifier string) (*authadminentity.Admin, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, a := range r.admins {
		if a.Email == identifier || a.Username == identifier || a.Phone == identifier {
			return a, true
		}
	}
	return nil, false
}

// ---------------- UPDATE ----------------
func (r *InMemoryAuthAdminRepository) Update(admin *authadminentity.Admin) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.admins[admin.Email] = admin
	return nil
}

// ---------------- DELETE ----------------
func (r *InMemoryAuthAdminRepository) Delete(email string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.admins, email)
	return nil
}

// ---------------- GET ALL ----------------
func (r *InMemoryAuthAdminRepository) GetAll() []*authadminentity.Admin {
	r.mu.RLock()
	defer r.mu.RUnlock()
	adminList := make([]*authadminentity.Admin, 0, len(r.admins))
	for _, a := range r.admins {
		adminList = append(adminList, a)
	}
	return adminList
}

// ---------------- COUNT ----------------
func (r *InMemoryAuthAdminRepository) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.admins)
}

// ---------------- EXISTS ----------------
func (r *InMemoryAuthAdminRepository) Exists(email string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, exists := r.admins[email]
	return exists
}

// ---------------- CLEAR ----------------
func (r *InMemoryAuthAdminRepository) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.admins = make(map[string]*authadminentity.Admin)
	r.nextID = 1
}

// ---------------- ME ----------------
func (r *InMemoryAuthAdminRepository) Me(email string) (*authadminentity.Admin, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	admin, exists := r.admins[email]
	return admin, exists
}
