package authadminentity

import (
	"time"

	"golang.org/x/crypto/bcrypt"

	"mqfm_backend/internal/helper"

)

// ---------------- ENTITY ----------------
type Admin struct {
	AdminID      int
	Email        string
	Username     string
	Phone        string
	Password     string
	ProfileImage string        // 🆕 untuk menyimpan nama/path gambar profil
	ProfileTheme ProfileTheme
	Token        string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// ProfileTheme berisi warna dan inisial profil admin
type ProfileTheme struct {
	Initial        string `json:"initial"`
	PrimaryColor   string `json:"primary_color"`
	BackgroundTop  string `json:"background_top"`
	BackgroundDown string `json:"background_down"`
	TextColor      string `json:"text_color"`
}

// ---------------- FACTORY FUNCTION ----------------
// digunakan saat register admin baru
func NewAdmin(email, username, phone, password, customInitial string) (*Admin, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// hasilkan tema profil default berdasarkan email
	defaultTheme := helper.GenerateDefaultProfileTheme(email)

	// kalau user punya custom inisial, ganti
	if customInitial != "" {
		defaultTheme.Initial = customInitial
	}

	return &Admin{
		Email:        email,
		Username:     username,
		Phone:        phone,
		Password:     string(hashed),
		ProfileTheme: ProfileTheme(defaultTheme),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}, nil
}

// ---------------- DOMAIN METHODS ----------------

// VerifyPassword memastikan password cocok saat login
func (a *Admin) VerifyPassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(a.Password), []byte(password)) == nil
}

// UpdateProfile memperbarui data profil admin (termasuk gambar)
func (a *Admin) UpdateProfile(newUsername, newPhone, newPassword, newInitial, newProfileImage string) error {
	if newUsername != "" {
		a.Username = newUsername
	}
	if newPhone != "" {
		a.Phone = newPhone
	}
	if newPassword != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		a.Password = string(hashed)
	}

	if newProfileImage != "" {
		// 🧠 kalau user upload gambar, kosongkan inisial dan simpan path image
		a.ProfileImage = newProfileImage
		a.ProfileTheme.Initial = ""
	} else if newInitial != "" {
		// kalau nggak upload gambar, pakai inisial baru
		a.ProfileTheme.Initial = newInitial
	}

	// ⚙️ regenerasi warna tapi JANGAN restore inisial lama
	newTheme := helper.GenerateDefaultProfileTheme(a.Email)
	newTheme.Initial = a.ProfileTheme.Initial // bisa kosong kalau udah punya foto
	a.ProfileTheme = ProfileTheme(newTheme)

	a.UpdatedAt = time.Now()
	return nil
}

