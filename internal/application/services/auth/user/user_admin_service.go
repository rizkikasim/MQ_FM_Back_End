package authuserservice

import (
	"errors"
	"time"

	authuserentity "mqfm_backend/internal/domain/entities/auth/user"
	authuserrepository "mqfm_backend/internal/infrastructure/repository/auth/user"
	domainutils "mqfm_backend/internal/domain/utils"

)

type AuthUserService struct {
	repo   authuserrepository.AuthUserRepository
	jwtKey string
}

func NewAuthUserService(repo authuserrepository.AuthUserRepository, jwtSecret string) *AuthUserService {
	return &AuthUserService{
		repo:   repo,
		jwtKey: jwtSecret,
	}
}

// ---------------- REGISTER ----------------
func (s *AuthUserService) Register(email, username, phone, password, customInitial string) (*authuserentity.User, error) {
	if existing, found := s.repo.FindByEmail(email); found && existing != nil {
		return nil, errors.New("email sudah digunakan")
	}
	if existing, found := s.repo.FindByIdentifier(username); found && existing != nil {
		return nil, errors.New("username sudah digunakan")
	}
	if existing, found := s.repo.FindByIdentifier(phone); found && existing != nil {
		return nil, errors.New("nomor telepon sudah digunakan")
	}

	user, err := authuserentity.NewUser(email, username, phone, password, customInitial)
	if err != nil {
		return nil, err
	}

	token, err := domainutils.GenerateJWT(user.Email, s.jwtKey, 24*time.Hour)
	if err != nil {
		return nil, errors.New("gagal membuat token JWT")
	}
	user.Token = token

	if err := s.repo.Save(user); err != nil {
		return nil, err
	}

	if err := s.repo.Update(user); err != nil {
		return nil, err
	}

	return sanitizeUser(user), nil
}

// ---------------- LOGIN ----------------
func (s *AuthUserService) Login(identifier, password string) (*authuserentity.User, error) {
	user, found := s.repo.FindByIdentifier(identifier)
	if !found {
		return nil, errors.New("akun tidak ditemukan")
	}

	if !user.VerifyPassword(password) {
		return nil, errors.New("password salah")
	}

	identifierForTheme := user.ProfileTheme.Initial
	if identifierForTheme == "" {
		identifierForTheme = user.Email
	}
	theme := domainutils.GenerateDefaultProfileTheme(identifierForTheme)
	user.ProfileTheme = authuserentity.ProfileTheme{
		Initial:        theme.Initial,
		PrimaryColor:   theme.PrimaryColor,
		BackgroundTop:  theme.BackgroundTop,
		BackgroundDown: theme.BackgroundDown,
		TextColor:      theme.TextColor,
	}

	token, err := domainutils.GenerateJWT(user.Email, s.jwtKey, 24*time.Hour)
	if err != nil {
		return nil, errors.New("gagal membuat token JWT")
	}

	user.Token = token
	if err := s.repo.Update(user); err != nil {
		return nil, err
	}

	return sanitizeUser(user), nil
}

// ---------------- ME (pakai Bearer Token) ----------------
func (s *AuthUserService) Me(token string) (*authuserentity.User, error) {
	email, err := domainutils.ParseJWT(token, s.jwtKey)
	if err != nil {
		return nil, errors.New("token tidak valid atau kadaluarsa")
	}

	user, found := s.repo.FindByEmail(email)
	if !found {
		return nil, errors.New("akun tidak ditemukan")
	}

	if user.Token == "" || user.Token != token {
		return nil, errors.New("token sudah tidak aktif, silakan login ulang")
	}

	identifierForTheme := user.ProfileTheme.Initial
	if identifierForTheme == "" {
		identifierForTheme = user.Email
	}
	theme := domainutils.GenerateDefaultProfileTheme(identifierForTheme)
	user.ProfileTheme = authuserentity.ProfileTheme{
		Initial:        theme.Initial,
		PrimaryColor:   theme.PrimaryColor,
		BackgroundTop:  theme.BackgroundTop,
		BackgroundDown: theme.BackgroundDown,
		TextColor:      theme.TextColor,
	}

	return sanitizeUser(user), nil
}

// ---------------- UPDATE ----------------
func (s *AuthUserService) Update(email, newUsername, newPhone, newPassword, newInitial, newProfileImage string) (*authuserentity.User, error) {
	user, found := s.repo.FindByEmail(email)
	if !found {
		return nil, errors.New("akun tidak ditemukan")
	}

	if newProfileImage != "" {
		newInitial = ""
	}

	if err := user.UpdateProfile(newUsername, newPhone, newPassword, newInitial, newProfileImage); err != nil {
		return nil, err
	}

	if user.ProfileImage == "" {
		identifierForTheme := user.ProfileTheme.Initial
		if identifierForTheme == "" {
			identifierForTheme = email
		}
		newTheme := domainutils.GenerateDefaultProfileTheme(identifierForTheme)
		user.ProfileTheme = authuserentity.ProfileTheme{
			Initial:        newTheme.Initial,
			PrimaryColor:   newTheme.PrimaryColor,
			BackgroundTop:  newTheme.BackgroundTop,
			BackgroundDown: newTheme.BackgroundDown,
			TextColor:      newTheme.TextColor,
		}
	}

	if err := s.repo.Update(user); err != nil {
		return nil, err
	}

	return sanitizeUser(user), nil
}

// ---------------- DELETE ACCOUNT ----------------
func (s *AuthUserService) DeleteAccount(email string) error {
	if _, found := s.repo.FindByEmail(email); !found {
		return errors.New("akun tidak ditemukan")
	}
	return s.repo.Delete(email)
}

// ---------------- LOGOUT ----------------
func (s *AuthUserService) Logout(email string) error {
	user, found := s.repo.FindByEmail(email)
	if !found {
		return errors.New("akun tidak ditemukan")
	}

	user.Token = ""
	return s.repo.Update(user)
}

// ---------------- Helper ----------------
func sanitizeUser(u *authuserentity.User) *authuserentity.User {
	clone := *u
	clone.Password = ""
	return &clone
}
