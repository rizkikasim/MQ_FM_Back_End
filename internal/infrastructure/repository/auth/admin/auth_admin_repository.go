package authadminrepository

import "mqfm_backend/internal/domain/entities/auth/admin"

type AuthAdminRepository interface {
	Save(admin *authadminentity.Admin) error
	FindByEmail(email string) (*authadminentity.Admin, bool)
	FindByIdentifier(identifier string) (*authadminentity.Admin, bool)
	Update(admin *authadminentity.Admin) error
	Delete(email string) error
	GetAll() []*authadminentity.Admin
	Count() int
	Exists(email string) bool
	Clear()
}
