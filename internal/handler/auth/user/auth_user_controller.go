package authusercontroller

import (
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/tidwall/gjson"

	authuserusecase "mqfm_backend/internal/domain/usecases/auth/user"
	authuserrepository "mqfm_backend/internal/infrastructure/repository/auth/user"
	errorinterceptor "mqfm_backend/internal/presentation/interceptor/errors"
	successinterceptor "mqfm_backend/internal/presentation/interceptor/success"
	validator "mqfm_backend/internal/presentation/validator"
	domainrequest "mqfm_backend/internal/domain/request/auth/user"

)

// ------------------------------------------------------------
// Controller
// ------------------------------------------------------------
type AuthUserController struct {
	usecase *authuserusecase.AuthUserUsecase
	repo    authuserrepository.AuthUserRepository
}

func NewAuthUserController(usecase *authuserusecase.AuthUserUsecase, repo authuserrepository.AuthUserRepository) *AuthUserController {
	return &AuthUserController{usecase: usecase, repo: repo}
}

// ---------------- REGISTER ----------------
func (c *AuthUserController) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errorinterceptor.ErrorInterceptor(w, r, http.ErrNotSupported, http.StatusMethodNotAllowed, 0)
		return
	}

	bodyBytes, _ := io.ReadAll(r.Body)
	req := domainrequest.RegisterUserRequest{
		Email:        gjson.GetBytes(bodyBytes, "email").String(),
		Username:     gjson.GetBytes(bodyBytes, "username").String(),
		Phone:        gjson.GetBytes(bodyBytes, "phone").String(),
		Password:     gjson.GetBytes(bodyBytes, "password").String(),
		ProfileImage: gjson.GetBytes(bodyBytes, "profile_image").String(),
	}

	if err := validator.ValidateRequiredFields(map[string]string{
		"email":    req.Email,
		"username": req.Username,
		"phone":    req.Phone,
		"password": req.Password,
	}); err != nil {
		errorinterceptor.ErrorInterceptor(w, r, err, http.StatusBadRequest, 0)
		return
	}

	if err := validator.ValidatePasswordLength(req.Password, 6); err != nil {
		errorinterceptor.ErrorInterceptor(w, r, err, http.StatusBadRequest, 0)
		return
	}

	user, err := c.usecase.Register(req)
	if err != nil {
		errorinterceptor.ErrorInterceptor(w, r, err, http.StatusBadRequest, 0)
		return
	}

	successinterceptor.SuccessInterceptor(w, r, http.StatusCreated, "register berhasil", user, user.UserID)
}

// ---------------- LOGIN ----------------
func (c *AuthUserController) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errorinterceptor.ErrorInterceptor(w, r, http.ErrNotSupported, http.StatusMethodNotAllowed, 0)
		return
	}

	bodyBytes, _ := io.ReadAll(r.Body)
	req := domainrequest.LoginUserRequest{
		Identifier: gjson.GetBytes(bodyBytes, "identifier").String(),
		Password:   gjson.GetBytes(bodyBytes, "password").String(),
	}

	if err := validator.ValidateRequiredFields(map[string]string{
		"identifier": req.Identifier,
		"password":   req.Password,
	}); err != nil {
		errorinterceptor.ErrorInterceptor(w, r, err, http.StatusBadRequest, 0)
		return
	}

	user, err := c.usecase.Login(req)
	if err != nil {
		errorinterceptor.ErrorInterceptor(w, r, err, http.StatusUnauthorized, 0)
		return
	}

	successinterceptor.SuccessInterceptor(w, r, http.StatusOK, "login berhasil", user, user.UserID)
}

// ---------------- UPDATE ----------------
func (c *AuthUserController) Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errorinterceptor.ErrorInterceptor(w, r, http.ErrNotSupported, http.StatusMethodNotAllowed, 0)
		return
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		errorinterceptor.ErrorInterceptor(w, r, errors.New("token tidak ditemukan"), http.StatusUnauthorized, 0)
		return
	}
	token := strings.TrimPrefix(authHeader, "Bearer ")

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		errorinterceptor.ErrorInterceptor(w, r, errors.New("gagal membaca form data"), http.StatusBadRequest, 0)
		return
	}

	_, header, _ := r.FormFile("new_profile_image")

	req := domainrequest.UpdateUserRequest{
		Token:           token,
		NewUsername:     r.FormValue("new_username"),
		NewPhone:        r.FormValue("new_phone"),
		NewPassword:     r.FormValue("new_password"),
		NewProfileImage: header,
	}

	if req.NewPassword != "" {
		if err := validator.ValidatePasswordLength(req.NewPassword, 6); err != nil {
			errorinterceptor.ErrorInterceptor(w, r, err, http.StatusBadRequest, 0)
			return
		}
	}

	user, err := c.usecase.Update(req)
	if err != nil {
		userID := 0
		if user != nil {
			userID = user.UserID
		}
		errorinterceptor.ErrorInterceptor(w, r, err, http.StatusBadRequest, userID)
		return
	}

	successinterceptor.SuccessInterceptor(w, r, http.StatusOK, "update berhasil", user, user.UserID)
}

// ---------------- DELETE ACCOUNT ----------------
func (c *AuthUserController) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		errorinterceptor.ErrorInterceptor(w, r, http.ErrNotSupported, http.StatusMethodNotAllowed, 0)
		return
	}

	bodyBytes, _ := io.ReadAll(r.Body)
	req := domainrequest.DeleteUserRequest{
		Email: gjson.GetBytes(bodyBytes, "email").String(),
	}

	if err := validator.ValidateRequiredFields(map[string]string{"email": req.Email}); err != nil {
		errorinterceptor.ErrorInterceptor(w, r, err, http.StatusBadRequest, 0)
		return
	}

	if err := c.usecase.DeleteAccount(req); err != nil {
		errorinterceptor.ErrorInterceptor(w, r, err, http.StatusBadRequest, 0)
		return
	}

	successinterceptor.SuccessInterceptor(w, r, http.StatusOK, "akun berhasil dihapus", nil, 0)
}

// ---------------- LOGOUT ----------------
func (c *AuthUserController) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errorinterceptor.ErrorInterceptor(w, r, http.ErrNotSupported, http.StatusMethodNotAllowed, 0)
		return
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		errorinterceptor.ErrorInterceptor(w, r, errors.New("token tidak ditemukan"), http.StatusUnauthorized, 0)
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")

	req := domainrequest.LogoutUserRequest{
		Token: token,
	}

	if err := c.usecase.Logout(req); err != nil {
		errorinterceptor.ErrorInterceptor(w, r, err, http.StatusBadRequest, 0)
		return
	}

	successinterceptor.SuccessInterceptor(w, r, http.StatusOK, "logout berhasil", nil, 0)
}

// ---------------- ME ----------------
func (c *AuthUserController) Me(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorinterceptor.ErrorInterceptor(w, r, http.ErrNotSupported, http.StatusMethodNotAllowed, 0)
		return
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		errorinterceptor.ErrorInterceptor(w, r, errors.New("token tidak ditemukan"), http.StatusUnauthorized, 0)
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")

	user, err := c.usecase.Me(domainrequest.MeUserRequest{Token: token})
	if err != nil {
		errorinterceptor.ErrorInterceptor(w, r, err, http.StatusUnauthorized, 0)
		return
	}

	successinterceptor.SuccessInterceptor(w, r, http.StatusOK, "data profil berhasil diambil", user, user.UserID)
}
