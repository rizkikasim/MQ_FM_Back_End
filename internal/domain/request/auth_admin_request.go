package domainrequest

import "mime/multipart"

// ---------------- REGISTER REQUEST ----------------
type RegisterAdminRequest struct {
	Email        string `json:"email" validate:"required,email"`
	Username     string `json:"username" validate:"required"`
	Phone        string `json:"phone" validate:"required"`
	Password     string `json:"password" validate:"required,min=6"`
	ProfileImage string `json:"profile_image,omitempty"`
}

// ---------------- LOGIN REQUEST ----------------
type LoginAdminRequest struct {
	Identifier string `json:"identifier" validate:"required"`
	Password   string `json:"password" validate:"required"`
}

type MeAdminRequest struct {
	Token string `json:"token"`
}

// ---------------- UPDATE REQUEST ----------------
// Versi baru — pakai token & file upload
type UpdateAdminRequest struct {
	Token           string                `form:"token" validate:"required"`
	NewUsername     string                `form:"new_username,omitempty"`
	NewPhone        string                `form:"new_phone,omitempty"`
	NewPassword     string                `form:"new_password,omitempty"`
	NewProfileImage *multipart.FileHeader `form:"new_profile_image,omitempty"`
}

// ---------------- DELETE ACCOUNT REQUEST ----------------
type DeleteAdminRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// ---------------- LOGOUT REQUEST ----------------
type LogoutAdminRequest struct {
	Token string `json:"token" validate:"required"`
}
