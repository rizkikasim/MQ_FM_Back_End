package authuserentity

import (
	"time"

	"golang.org/x/crypto/bcrypt"

	"mqfm_backend/internal/helper"

)

// ---------------- ENTITY ----------------
type User struct {
	UserID       int
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

// ProfileTheme berisi warna dan inisial profil user
type ProfileTheme struct {
	Initial        string `json:"initial"`
	PrimaryColor   string `json:"primary_color"`
	BackgroundTop  string `json:"background_top"`
	BackgroundDown string `json:"background_down"`
	TextColor      string `json:"text_color"`
}

// ---------------- FACTORY FUNCTION ----------------
// digunakan saat register user baru
func NewUser(email, username, phone, password, customInitial string) (*User, error) {
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

	return &User{
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
func (u *User) VerifyPassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)) == nil
}

// UpdateProfile memperbarui data profil user (termasuk gambar)
func (u *User) UpdateProfile(newUsername, newPhone, newPassword, newInitial, newProfileImage string) error {
	if newUsername != "" {
		u.Username = newUsername
	}
	if newPhone != "" {
		u.Phone = newPhone
	}
	if newPassword != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = string(hashed)
	}

	if newProfileImage != "" {
		// 🧠 kalau user upload gambar, kosongkan inisial dan simpan path image
		u.ProfileImage = newProfileImage
		u.ProfileTheme.Initial = ""
	} else if newInitial != "" {
		// kalau nggak upload gambar, pakai inisial baru
		u.ProfileTheme.Initial = newInitial
	}

	// ⚙️ regenerasi warna tapi JANGAN restore inisial lama
	newTheme := helper.GenerateDefaultProfileTheme(u.Email)
	newTheme.Initial = u.ProfileTheme.Initial // bisa kosong kalau udah punya foto
	u.ProfileTheme = ProfileTheme(newTheme)

	u.UpdatedAt = time.Now()
	return nil
}
