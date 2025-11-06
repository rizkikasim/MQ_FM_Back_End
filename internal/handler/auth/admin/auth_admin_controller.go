package authadmincontroller

import (
	"errors"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/tidwall/gjson"

	authadminusecase "mqfm_backend/internal/domain/usecases/auth/admin"
	authadminrepository "mqfm_backend/internal/infrastructure/repository/auth/admin"
	errorinterceptor "mqfm_backend/internal/presentation/interceptor/errors"
	successinterceptor "mqfm_backend/internal/presentation/interceptor/success"
	validator "mqfm_backend/internal/presentation/validator"
	domainrequest "mqfm_backend/internal/domain/request"

)

// ------------------------------------------------------------
// Controller
// ------------------------------------------------------------
type AuthAdminController struct {
	usecase *authadminusecase.AuthAdminUsecase
	repo    authadminrepository.AuthAdminRepository
}

func NewAuthAdminController(usecase *authadminusecase.AuthAdminUsecase, repo authadminrepository.AuthAdminRepository) *AuthAdminController {
	return &AuthAdminController{usecase: usecase, repo: repo}
}

// ---------------- REGISTER ----------------
func (c *AuthAdminController) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errorinterceptor.ErrorInterceptor(w, r, http.ErrNotSupported, http.StatusMethodNotAllowed, 0)
		return
	}

	admins := c.repo.GetAll()
	if len(admins) >= 3 {
		logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
		logger.Warn("register ditolak karena batas admin tercapai",
			slog.Int("count", len(admins)),
			slog.String("path", r.URL.Path),
		)
		errorinterceptor.ErrorInterceptor(w, r, http.ErrBodyNotAllowed, http.StatusForbidden, 0)
		return
	}

	bodyBytes, _ := io.ReadAll(r.Body)
	req := domainrequest.RegisterAdminRequest{
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

	admin, err := c.usecase.Register(req)
	if err != nil {
		errorinterceptor.ErrorInterceptor(w, r, err, http.StatusBadRequest, 0)
		return
	}

	successinterceptor.SuccessInterceptor(w, r, http.StatusCreated, "register berhasil", admin, admin.AdminID)
}

// ---------------- LOGIN ----------------
func (c *AuthAdminController) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errorinterceptor.ErrorInterceptor(w, r, http.ErrNotSupported, http.StatusMethodNotAllowed, 0)
		return
	}

	bodyBytes, _ := io.ReadAll(r.Body)
	req := domainrequest.LoginAdminRequest{
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

	admin, err := c.usecase.Login(req)
	if err != nil {
		errorinterceptor.ErrorInterceptor(w, r, err, http.StatusUnauthorized, 0)
		return
	}

	successinterceptor.SuccessInterceptor(w, r, http.StatusOK, "login berhasil", admin, admin.AdminID)
}

// ---------------- UPDATE ----------------
func (c *AuthAdminController) Update(w http.ResponseWriter, r *http.Request) {
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

	// ✅ Parse multipart form (max 10 MB)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		errorinterceptor.ErrorInterceptor(w, r, errors.New("gagal membaca form data"), http.StatusBadRequest, 0)
		return
	}

	// ambil file header (boleh nil)
	_, header, _ := r.FormFile("new_profile_image")

	req := domainrequest.UpdateAdminRequest{
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

	admin, err := c.usecase.Update(req)
	if err != nil {
		adminID := 0
		if admin != nil {
			adminID = admin.AdminID
		}
		errorinterceptor.ErrorInterceptor(w, r, err, http.StatusBadRequest, adminID)
		return
	}

	successinterceptor.SuccessInterceptor(w, r, http.StatusOK, "update berhasil", admin, admin.AdminID)
}

// ---------------- DELETE ACCOUNT ----------------
func (c *AuthAdminController) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		errorinterceptor.ErrorInterceptor(w, r, http.ErrNotSupported, http.StatusMethodNotAllowed, 0)
		return
	}

	bodyBytes, _ := io.ReadAll(r.Body)
	req := domainrequest.DeleteAdminRequest{
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
func (c *AuthAdminController) Logout(w http.ResponseWriter, r *http.Request) {
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

	req := domainrequest.LogoutAdminRequest{
		Token: token,
	}

	if err := c.usecase.Logout(req); err != nil {
		errorinterceptor.ErrorInterceptor(w, r, err, http.StatusBadRequest, 0)
		return
	}

	successinterceptor.SuccessInterceptor(w, r, http.StatusOK, "logout berhasil", nil, 0)
}

// ---------------- ME ----------------
func (c *AuthAdminController) Me(w http.ResponseWriter, r *http.Request) {
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

	admin, err := c.usecase.Me(domainrequest.MeAdminRequest{Token: token})
	if err != nil {
		errorinterceptor.ErrorInterceptor(w, r, err, http.StatusUnauthorized, 0)
		return
	}

	successinterceptor.SuccessInterceptor(w, r, http.StatusOK, "data profil berhasil diambil", admin, admin.AdminID)
}
