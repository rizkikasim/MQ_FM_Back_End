package authuserusecase

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	authuserentity "mqfm_backend/internal/domain/entities/auth/user"
	authuserservice "mqfm_backend/internal/application/services/auth/user"
	req "mqfm_backend/internal/domain/request/auth/user"

)

type AuthUserUsecase struct {
	service *authuserservice.AuthUserService
}

func NewAuthUserUsecase(service *authuserservice.AuthUserService) *AuthUserUsecase {
	return &AuthUserUsecase{service: service}
}

// ---------------- REGISTER ----------------
func (u *AuthUserUsecase) Register(request req.RegisterUserRequest) (*authuserentity.User, error) {
	user, err := u.service.Register(
		request.Email,
		request.Username,
		request.Phone,
		request.Password,
		request.ProfileImage,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// ---------------- LOGIN ----------------
func (u *AuthUserUsecase) Login(request req.LoginUserRequest) (*authuserentity.User, error) {
	user, err := u.service.Login(
		request.Identifier,
		request.Password,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// ---------------- ME ----------------
func (u *AuthUserUsecase) Me(request req.MeUserRequest) (*authuserentity.User, error) {
	return u.service.Me(request.Token)
}

// ---------------- UPDATE ----------------
func (u *AuthUserUsecase) Update(request req.UpdateUserRequest) (*authuserentity.User, error) {
	user, err := u.service.Me(request.Token)
	if err != nil {
		return nil, err
	}

	newInitial := user.ProfileTheme.Initial
	newProfileImage := user.ProfileImage

	if request.NewProfileImage != nil {
		uploadDir := "storage/profile"
		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
			return nil, fmt.Errorf("gagal membuat folder upload: %v", err)
		}

		filename := fmt.Sprintf("user_%d_%s", user.UserID, request.NewProfileImage.Filename)
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

		newProfileImage = filename
	}

	updatedUser, err := u.service.Update(
		user.Email,
		request.NewUsername,
		request.NewPhone,
		request.NewPassword,
		newInitial,
		newProfileImage,
	)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}

// ---------------- DELETE ACCOUNT ----------------
func (u *AuthUserUsecase) DeleteAccount(request req.DeleteUserRequest) error {
	return u.service.DeleteAccount(request.Email)
}

// ---------------- LOGOUT ----------------
func (u *AuthUserUsecase) Logout(request req.LogoutUserRequest) error {
	user, err := u.service.Me(request.Token)
	if err != nil {
		return err
	}

	return u.service.Logout(user.Email)
}
