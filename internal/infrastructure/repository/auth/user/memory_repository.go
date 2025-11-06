package authuserrepository

import (
	"sync"

	authuserentity "mqfm_backend/internal/domain/entities/auth/user"

)

// ---------------- REPOSITORY (IN-MEMORY) ----------------
type InMemoryAuthUserRepository struct {
	users  map[string]*authuserentity.User
	mu     sync.RWMutex
	nextID int
}

// ---------------- CONSTRUCTOR ----------------
func NewInMemoryAuthUserRepository() *InMemoryAuthUserRepository {
	return &InMemoryAuthUserRepository{
		users:  make(map[string]*authuserentity.User),
		nextID: 1,
	}
}

// ---------------- SAVE ----------------
func (r *InMemoryAuthUserRepository) Save(user *authuserentity.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if user.UserID == 0 {
		user.UserID = r.nextID
		r.nextID++
	}
	r.users[user.Email] = user
	return nil
}

// ---------------- FIND BY EMAIL ----------------
func (r *InMemoryAuthUserRepository) FindByEmail(email string) (*authuserentity.User, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	user, exists := r.users[email]
	return user, exists
}

// ---------------- FIND BY IDENTIFIER ----------------
func (r *InMemoryAuthUserRepository) FindByIdentifier(identifier string) (*authuserentity.User, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, u := range r.users {
		if u.Email == identifier || u.Username == identifier || u.Phone == identifier {
			return u, true
		}
	}
	return nil, false
}

// ---------------- UPDATE ----------------
func (r *InMemoryAuthUserRepository) Update(user *authuserentity.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.users[user.Email] = user
	return nil
}

// ---------------- DELETE ----------------
func (r *InMemoryAuthUserRepository) Delete(email string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.users, email)
	return nil
}

// ---------------- GET ALL ----------------
func (r *InMemoryAuthUserRepository) GetAll() []*authuserentity.User {
	r.mu.RLock()
	defer r.mu.RUnlock()
	userList := make([]*authuserentity.User, 0, len(r.users))
	for _, u := range r.users {
		userList = append(userList, u)
	}
	return userList
}

// ---------------- COUNT ----------------
func (r *InMemoryAuthUserRepository) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.users)
}

// ---------------- EXISTS ----------------
func (r *InMemoryAuthUserRepository) Exists(email string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, exists := r.users[email]
	return exists
}

// ---------------- CLEAR ----------------
func (r *InMemoryAuthUserRepository) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.users = make(map[string]*authuserentity.User)
	r.nextID = 1
}

// ---------------- ME ----------------
func (r *InMemoryAuthUserRepository) Me(email string) (*authuserentity.User, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	user, exists := r.users[email]
	return user, exists
}
