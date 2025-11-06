package domainrequest

import "mime/multipart"

// ---------------- CREATE CATEGORY REQUEST ----------------
type CreateCategoryAdminRequest struct {
	Name        string                `form:"name" validate:"required"`
	Description string                `form:"description,omitempty"`
	Image       *multipart.FileHeader `form:"image,omitempty"`
}

// ---------------- UPDATE CATEGORY REQUEST ----------------
type UpdateCategoryAdminRequest struct {
	ID          int                   `form:"id" validate:"required"`
	NewName     string                `form:"new_name,omitempty"`
	NewDesc     string                `form:"new_description,omitempty"`
	NewImage    *multipart.FileHeader `form:"new_image,omitempty"`
}

// ---------------- DELETE CATEGORY REQUEST ----------------
type DeleteCategoryAdminRequest struct {
	ID int `json:"id" validate:"required"`
}

// ---------------- GET CATEGORY REQUEST ----------------
type GetCategoryAdminRequest struct {
	ID int `json:"id" validate:"required"`
}
