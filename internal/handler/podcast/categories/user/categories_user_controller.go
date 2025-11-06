package categoriesusercontroller

import (
	"errors"
	"net/http"
	"strconv"

	categoriesuserservice "mqfm_backend/internal/application/services/podcast/categories/user"
	errorinterceptor "mqfm_backend/internal/presentation/interceptor/errors"
	successinterceptor "mqfm_backend/internal/presentation/interceptor/success"

)

// ------------------------------------------------------------
// Controller
// ------------------------------------------------------------
type CategoriesUserController struct {
	service *categoriesuserservice.CategoriesUserService
}

// constructor
func NewCategoriesUserController(service *categoriesuserservice.CategoriesUserService) *CategoriesUserController {
	return &CategoriesUserController{service: service}
}

// ---------------- GET ALL ----------------
func (c *CategoriesUserController) GetAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorinterceptor.ErrorInterceptor(w, r, http.ErrNotSupported, http.StatusMethodNotAllowed, 0)
		return
	}

	categories, err := c.service.GetAll()
	if err != nil {
		errorinterceptor.ErrorInterceptor(w, r, err, http.StatusInternalServerError, 0)
		return
	}

	successinterceptor.SuccessInterceptor(w, r, http.StatusOK, "daftar kategori berhasil diambil", categories, len(categories))
}

// ---------------- GET BY ID ----------------
func (c *CategoriesUserController) GetByID(w http.ResponseWriter, r *http.Request) {
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

	category, err := c.service.GetByID(id)
	if err != nil {
		errorinterceptor.ErrorInterceptor(w, r, err, http.StatusNotFound, id)
		return
	}

	successinterceptor.SuccessInterceptor(w, r, http.StatusOK, "kategori berhasil diambil", category, category.ID)
}
