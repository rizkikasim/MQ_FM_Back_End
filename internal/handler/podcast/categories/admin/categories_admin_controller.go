package categoriesadmincontroller

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	categoriesadminusecase "mqfm_backend/internal/domain/usecases/podcast/categories/admin"
	errorinterceptor "mqfm_backend/internal/presentation/interceptor/errors"
	successinterceptor "mqfm_backend/internal/presentation/interceptor/success"
	validator "mqfm_backend/internal/presentation/validator"
	domainrequest "mqfm_backend/internal/domain/request/podcast/categories/admin"

)

// ------------------------------------------------------------
// Controller
// ------------------------------------------------------------
type CategoriesAdminController struct {
	usecase *categoriesadminusecase.CategoriesAdminUsecase
}

func NewCategoriesAdminController(usecase *categoriesadminusecase.CategoriesAdminUsecase) *CategoriesAdminController {
	return &CategoriesAdminController{usecase: usecase}
}

// ---------------- CREATE ----------------
func (c *CategoriesAdminController) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errorinterceptor.ErrorInterceptor(w, r, http.ErrNotSupported, http.StatusMethodNotAllowed, 0)
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		errorinterceptor.ErrorInterceptor(w, r, errors.New("gagal membaca form data"), http.StatusBadRequest, 0)
		return
	}

	file, header, _ := r.FormFile("image")

	req := domainrequest.CreateCategoryAdminRequest{
		Name:        r.FormValue("name"),
		Description: r.FormValue("description"),
	}

	if header != nil {
		req.Image = header
	}
	if file != nil {
		file.Close()
	}

	if err := validator.ValidateRequiredFields(map[string]string{
		"name": req.Name,
	}); err != nil {
		errorinterceptor.ErrorInterceptor(w, r, err, http.StatusBadRequest, 0)
		return
	}

	category, err := c.usecase.Create(req)
	if err != nil {
		errorinterceptor.ErrorInterceptor(w, r, err, http.StatusBadRequest, 0)
		return
	}

	successinterceptor.SuccessInterceptor(w, r, http.StatusCreated, "kategori berhasil dibuat", category, category.ID)
}

// ---------------- GET ALL ----------------
func (c *CategoriesAdminController) GetAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorinterceptor.ErrorInterceptor(w, r, http.ErrNotSupported, http.StatusMethodNotAllowed, 0)
		return
	}

	categories, err := c.usecase.GetAll()
	if err != nil {
		errorinterceptor.ErrorInterceptor(w, r, err, http.StatusInternalServerError, 0)
		return
	}

	successinterceptor.SuccessInterceptor(w, r, http.StatusOK, "daftar kategori berhasil diambil", categories, len(categories))
}

// ---------------- GET BY ID ----------------
func (c *CategoriesAdminController) GetByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorinterceptor.ErrorInterceptor(w, r, http.ErrNotSupported, http.StatusMethodNotAllowed, 0)
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		errorinterceptor.ErrorInterceptor(w, r, errors.New("parameter id wajib diisi"), http.StatusBadRequest, 0)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		errorinterceptor.ErrorInterceptor(w, r, errors.New("id tidak valid"), http.StatusBadRequest, 0)
		return
	}

	category, err := c.usecase.GetByID(domainrequest.GetCategoryAdminRequest{ID: id})
	if err != nil {
		errorinterceptor.ErrorInterceptor(w, r, err, http.StatusNotFound, id)
		return
	}

	successinterceptor.SuccessInterceptor(w, r, http.StatusOK, "kategori berhasil diambil", category, category.ID)
}

// ---------------- UPDATE ----------------
func (c *CategoriesAdminController) Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errorinterceptor.ErrorInterceptor(w, r, http.ErrNotSupported, http.StatusMethodNotAllowed, 0)
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		errorinterceptor.ErrorInterceptor(w, r, errors.New("gagal membaca form data"), http.StatusBadRequest, 0)
		return
	}

	file, header, _ := r.FormFile("new_image")

	idStr := r.FormValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		errorinterceptor.ErrorInterceptor(w, r, errors.New("id kategori tidak valid"), http.StatusBadRequest, 0)
		return
	}

	req := domainrequest.UpdateCategoryAdminRequest{
		ID:       id,
		NewName:  r.FormValue("new_name"),
		NewDesc:  r.FormValue("new_description"),
		NewImage: header,
	}

	if file != nil {
		file.Close()
	}

	category, err := c.usecase.Update(req)
	if err != nil {
		errorinterceptor.ErrorInterceptor(w, r, err, http.StatusBadRequest, req.ID)
		return
	}

	successinterceptor.SuccessInterceptor(w, r, http.StatusOK, "kategori berhasil diperbarui", category, category.ID)
}

// ---------------- DELETE ----------------
func (c *CategoriesAdminController) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		errorinterceptor.ErrorInterceptor(w, r, http.ErrNotSupported, http.StatusMethodNotAllowed, 0)
		return
	}

	defer r.Body.Close()

	var body struct {
		ID int `json:"id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		errorinterceptor.ErrorInterceptor(w, r, errors.New("body tidak valid, pastikan format JSON benar"), http.StatusBadRequest, 0)
		return
	}

	if body.ID == 0 {
		errorinterceptor.ErrorInterceptor(w, r, errors.New("id wajib diisi dan harus berupa angka"), http.StatusBadRequest, 0)
		return
	}

	req := domainrequest.DeleteCategoryAdminRequest{ID: body.ID}

	if err := c.usecase.Delete(req); err != nil {
		errorinterceptor.ErrorInterceptor(w, r, err, http.StatusBadRequest, body.ID)
		return
	}

	successinterceptor.SuccessInterceptor(w, r, http.StatusOK, "kategori berhasil dihapus", nil, body.ID)
}
