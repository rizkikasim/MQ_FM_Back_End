package authuserrepository

import authuserentity "mqfm_backend/internal/domain/entities/auth/user"

type AuthUserRepository interface {
	Save(user *authuserentity.User) error
	FindByEmail(email string) (*authuserentity.User, bool)
	FindByIdentifier(identifier string) (*authuserentity.User, bool)
	Update(user *authuserentity.User) error
	Delete(email string) error
	GetAll() []*authuserentity.User
	Count() int
	Exists(email string) bool
	Clear()
}
