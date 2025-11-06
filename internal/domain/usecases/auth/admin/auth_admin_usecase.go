package authadminusecase

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	authadminentity "mqfm_backend/internal/domain/entities/auth/admin"
	authadminservice "mqfm_backend/internal/application/services/auth/admin"
	req "mqfm_backend/internal/domain/request"

)

type AuthAdminUsecase struct {
	service *authadminservice.AuthAdminService
}

func NewAuthAdminUsecase(service *authadminservice.AuthAdminService) *AuthAdminUsecase {
	return &AuthAdminUsecase{service: service}
}

// ---------------- REGISTER ----------------
func (u *AuthAdminUsecase) Register(request req.RegisterAdminRequest) (*authadminentity.Admin, error) {
	admin, err := u.service.Register(
		request.Email,
		request.Username,
		request.Phone,
		request.Password,
		request.ProfileImage,
	)
	if err != nil {
		return nil, err
	}
	return admin, nil
}

// ---------------- LOGIN ----------------
func (u *AuthAdminUsecase) Login(request req.LoginAdminRequest) (*authadminentity.Admin, error) {
	admin, err := u.service.Login(
		request.Identifier,
		request.Password,
	)
	if err != nil {
		return nil, err
	}
	return admin, nil
}

// ---------------- ME ----------------
func (u *AuthAdminUsecase) Me(request req.MeAdminRequest) (*authadminentity.Admin, error) {
	return u.service.Me(request.Token)
}

// ---------------- UPDATE ----------------
func (u *AuthAdminUsecase) Update(request req.UpdateAdminRequest) (*authadminentity.Admin, error) {
	admin, err := u.service.Me(request.Token)
	if err != nil {
		return nil, err
	}

	newInitial := admin.ProfileTheme.Initial // default pakai yang lama
	newProfileImage := admin.ProfileImage     // 🆕 default ke gambar lama

	// ✅ handle upload gambar baru
	if request.NewProfileImage != nil {
		uploadDir := "storage/profile"
		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
			return nil, fmt.Errorf("gagal membuat folder upload: %v", err)
		}

		filename := fmt.Sprintf("admin_%d_%s", admin.AdminID, request.NewProfileImage.Filename)
		filePath := filepath.Join(uploadDir, filename)

		src, err := request.NewProfileImage.Open()
		if err != nil {
			return nil, fmt.Errorf("gagal membuka file upload: %v", err)
		}
		defer src.Close()

		dst, err := os.Create(filePath)
		if err != nil {
			return nil, fmt.Errorf("gagal menyimpan file upload: %v", err)
		}
		defer dst.Close()

		if _, err := io.Copy(dst, src); err != nil {
			return nil, fmt.Errorf("gagal menulis file upload: %v", err)
		}

		// 🆕 simpan hanya nama file (tanpa path)
		newProfileImage = filename
	}

	// ✅ kirim semua ke service (6 argumen)
	updatedAdmin, err := u.service.Update(
		admin.Email,
		request.NewUsername,
		request.NewPhone,
		request.NewPassword,
		newInitial,
		newProfileImage, // 🆕 tambahkan ini
	)
	if err != nil {
		return nil, err
	}

	return updatedAdmin, nil
}


// ---------------- DELETE ACCOUNT ----------------
func (u *AuthAdminUsecase) DeleteAccount(request req.DeleteAdminRequest) error {
	return u.service.DeleteAccount(request.Email)
}

// ---------------- LOGOUT ----------------
func (u *AuthAdminUsecase) Logout(request req.LogoutAdminRequest) error {
	admin, err := u.service.Me(request.Token)
	if err != nil {
		return err
	}

	return u.service.Logout(admin.Email)
}
