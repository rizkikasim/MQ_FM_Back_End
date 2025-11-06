package authadminservice

import (
	"errors"
	"time"

	authadminentity "mqfm_backend/internal/domain/entities/auth/admin"
	authadminrepository "mqfm_backend/internal/infrastructure/repository/auth/admin"
	domainutils "mqfm_backend/internal/domain/utils"

)

type AuthAdminService struct {
	repo   authadminrepository.AuthAdminRepository
	jwtKey string
}

func NewAuthAdminService(repo authadminrepository.AuthAdminRepository, jwtSecret string) *AuthAdminService {
	return &AuthAdminService{
		repo:   repo,
		jwtKey: jwtSecret,
	}
}

// ---------------- REGISTER ----------------
func (s *AuthAdminService) Register(email, username, phone, password, customInitial string) (*authadminentity.Admin, error) {
	if existing, found := s.repo.FindByEmail(email); found && existing != nil {
		return nil, errors.New("email sudah digunakan")
	}
	if existing, found := s.repo.FindByIdentifier(username); found && existing != nil {
		return nil, errors.New("username sudah digunakan")
	}
	if existing, found := s.repo.FindByIdentifier(phone); found && existing != nil {
		return nil, errors.New("nomor telepon sudah digunakan")
	}

	admin, err := authadminentity.NewAdmin(email, username, phone, password, customInitial)
	if err != nil {
		return nil, err
	}

	token, err := domainutils.GenerateJWT(admin.Email, s.jwtKey, 24*time.Hour)
	if err != nil {
		return nil, errors.New("gagal membuat token JWT")
	}
	admin.Token = token

	if err := s.repo.Save(admin); err != nil {
		return nil, err
	}

	if err := s.repo.Update(admin); err != nil {
		return nil, err
	}

	return sanitizeAdmin(admin), nil
}

// ---------------- LOGIN ----------------
func (s *AuthAdminService) Login(identifier, password string) (*authadminentity.Admin, error) {
	admin, found := s.repo.FindByIdentifier(identifier)
	if !found {
		return nil, errors.New("akun tidak ditemukan")
	}

	if !admin.VerifyPassword(password) {
		return nil, errors.New("password salah")
	}

	identifierForTheme := admin.ProfileTheme.Initial
	if identifierForTheme == "" {
		identifierForTheme = admin.Email
	}
	theme := domainutils.GenerateDefaultProfileTheme(identifierForTheme)
	admin.ProfileTheme = theme

	token, err := domainutils.GenerateJWT(admin.Email, s.jwtKey, 24*time.Hour)
	if err != nil {
		return nil, errors.New("gagal membuat token JWT")
	}

	admin.Token = token
	if err := s.repo.Update(admin); err != nil {
		return nil, err
	}

	return sanitizeAdmin(admin), nil
}

// ---------------- ME (pakai Bearer Token) ----------------
func (s *AuthAdminService) Me(token string) (*authadminentity.Admin, error) {
	email, err := domainutils.ParseJWT(token, s.jwtKey)
	if err != nil {
		return nil, errors.New("token tidak valid atau kadaluarsa")
	}

	admin, found := s.repo.FindByEmail(email)
	if !found {
		return nil, errors.New("akun tidak ditemukan")
	}

	if admin.Token == "" || admin.Token != token {
		return nil, errors.New("token sudah tidak aktif, silakan login ulang")
	}

	identifierForTheme := admin.ProfileTheme.Initial
	if identifierForTheme == "" {
		identifierForTheme = admin.Email
	}
	theme := domainutils.GenerateDefaultProfileTheme(identifierForTheme)
	admin.ProfileTheme = theme

	return sanitizeAdmin(admin), nil
}

// ---------------- UPDATE ----------------
func (s *AuthAdminService) Update(email, newUsername, newPhone, newPassword, newInitial, newProfileImage string) (*authadminentity.Admin, error) {
	admin, found := s.repo.FindByEmail(email)
	if !found {
		return nil, errors.New("akun tidak ditemukan")
	}

	// kalau ada foto baru → kosongkan inisial
	if newProfileImage != "" {
		newInitial = ""
	}

	if err := admin.UpdateProfile(newUsername, newPhone, newPassword, newInitial, newProfileImage); err != nil {
		return nil, err
	}

	// 🧠 jangan generate ulang theme pakai email kalau udah ada foto profil
	if admin.ProfileImage == "" {
		identifierForTheme := admin.ProfileTheme.Initial
		if identifierForTheme == "" {
			identifierForTheme = email
		}
		newTheme := domainutils.GenerateDefaultProfileTheme(identifierForTheme)
		admin.ProfileTheme = newTheme
	}

	if err := s.repo.Update(admin); err != nil {
		return nil, err
	}

	return sanitizeAdmin(admin), nil
}


// ---------------- DELETE ACCOUNT ----------------
func (s *AuthAdminService) DeleteAccount(email string) error {
	if _, found := s.repo.FindByEmail(email); !found {
		return errors.New("akun tidak ditemukan")
	}
	return s.repo.Delete(email)
}

// ---------------- LOGOUT ----------------
func (s *AuthAdminService) Logout(email string) error {
	admin, found := s.repo.FindByEmail(email)
	if !found {
		return errors.New("akun tidak ditemukan")
	}

	admin.Token = ""
	return s.repo.Update(admin)
}

// ---------------- Helper ----------------
func sanitizeAdmin(a *authadminentity.Admin) *authadminentity.Admin {
	clone := *a
	clone.Password = ""
	return &clone
}
