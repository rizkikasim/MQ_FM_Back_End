package playlistcontroller

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	playlistusecase "mqfm_backend/internal/domain/usecases/playlist"
	domainrequest "mqfm_backend/internal/domain/request/playlist"
	errorinterceptor "mqfm_backend/internal/presentation/interceptor/errors"
	successinterceptor "mqfm_backend/internal/presentation/interceptor/success"
	validator "mqfm_backend/internal/presentation/validator"

)

// ------------------------------------------------------------
// Controller
// ------------------------------------------------------------
type PlaylistController struct {
	usecase *playlistusecase.PlaylistUsecase
}

func NewPlaylistController(usecase *playlistusecase.PlaylistUsecase) *PlaylistController {
	return &PlaylistController{usecase: usecase}
}

// ============================================================
// 🔥 CREATE PLAYLIST
// ============================================================
func (c *PlaylistController) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errorinterceptor.ErrorInterceptor(w, r, http.ErrNotSupported, http.StatusMethodNotAllowed, 0)
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		errorinterceptor.ErrorInterceptor(w, r, errors.New("gagal membaca form-data"), http.StatusBadRequest, 0)
		return
	}

	req := domainrequest.CreatePlaylistRequest{
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
		Thumbnail:   r.FormValue("thumbnail"), // string URL
		UserID:      parseInt(r.FormValue("user_id")),
	}

	// validate required
	if err := validator.ValidateRequiredFields(map[string]string{
		"title":   req.Title,
		"user_id": r.FormValue("user_id"),
	}); err != nil {
		errorinterceptor.ErrorInterceptor(w, r, err, http.StatusBadRequest, 0)
		return
	}

	item, err := c.usecase.Create(req)
	if err != nil {
		errorinterceptor.ErrorInterceptor(w, r, err, http.StatusBadRequest, 0)
		return
	}

	successinterceptor.SuccessInterceptor(
		w, r,
		http.StatusCreated,
		"playlist berhasil dibuat",
		item,
		item.ID,
	)
}

// ============================================================
// 🔥 GET ALL BY USER
// ============================================================
func (c *PlaylistController) GetAllByUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorinterceptor.ErrorInterceptor(w, r, http.ErrNotSupported, http.StatusMethodNotAllowed, 0)
		return
	}

	userID := parseInt(r.URL.Query().Get("user_id"))
	if userID == 0 {
		errorinterceptor.ErrorInterceptor(w, r, errors.New("user_id wajib diisi"), http.StatusBadRequest, 0)
		return
	}

	items, err := c.usecase.GetAllByUser(userID)
	if err != nil {
		errorinterceptor.ErrorInterceptor(w, r, err, http.StatusInternalServerError, 0)
		return
	}

	successinterceptor.SuccessInterceptor(
		w, r,
		http.StatusOK,
		"daftar playlist berhasil diambil",
		items,
		len(items),
	)
}

// ============================================================
// 🔥 GET PLAYLIST BY ID
// ============================================================
func (c *PlaylistController) GetByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorinterceptor.ErrorInterceptor(w, r, http.ErrNotSupported, http.StatusMethodNotAllowed, 0)
		return
	}

	id := parseInt(r.URL.Query().Get("id"))
	if id == 0 {
		errorinterceptor.ErrorInterceptor(w, r, errors.New("id tidak valid"), http.StatusBadRequest, 0)
		return
	}

	item, err := c.usecase.GetByID(domainrequest.GetPlaylistRequest{ID: id})
	if err != nil {
		errorinterceptor.ErrorInterceptor(w, r, err, http.StatusNotFound, id)
		return
	}

	successinterceptor.SuccessInterceptor(
		w, r,
		http.StatusOK,
		"playlist berhasil diambil",
		item,
		item.ID,
	)
}

// ============================================================
// 🔥 UPDATE PLAYLIST
// ============================================================
func (c *PlaylistController) Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errorinterceptor.ErrorInterceptor(w, r, http.ErrNotSupported, http.StatusMethodNotAllowed, 0)
		return
	}

	if err := r.ParseForm(); err != nil {
		errorinterceptor.ErrorInterceptor(w, r, errors.New("gagal membaca form"), http.StatusBadRequest, 0)
		return
	}

	req := domainrequest.UpdatePlaylistRequest{
		ID:       parseInt(r.FormValue("id")),
		NewTitle: r.FormValue("new_title"),
		NewDesc:  r.FormValue("new_description"),
		NewThumb: r.FormValue("new_thumbnail"),
	}

	if req.ID == 0 {
		errorinterceptor.ErrorInterceptor(w, r, errors.New("id wajib diisi"), http.StatusBadRequest, 0)
		return
	}

	item, err := c.usecase.Update(req)
	if err != nil {
		errorinterceptor.ErrorInterceptor(w, r, err, req.ID)
		return
	}

	successinterceptor.SuccessInterceptor(
		w, r,
		http.StatusOK,
		"playlist berhasil diperbarui",
		item,
		item.ID,
	)
}

// ============================================================
// 🔥 DELETE PLAYLIST
// ============================================================
func (c *PlaylistController) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		errorinterceptor.ErrorInterceptor(w, r, http.ErrNotSupported, http.StatusMethodNotAllowed, 0)
		return
	}

	defer r.Body.Close()

	var body struct {
		ID int `json:"id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		errorinterceptor.ErrorInterceptor(
			w, r,
			errors.New("body JSON tidak valid"),
			http.StatusBadRequest,
			0,
		)
		return
	}

	if body.ID == 0 {
		errorinterceptor.ErrorInterceptor(
			w, r,
			errors.New("id wajib diisi"),
			http.StatusBadRequest,
			0,
		)
		return
	}

	if err := c.usecase.Delete(domainrequest.DeletePlaylistRequest{ID: body.ID}); err != nil {
		errorinterceptor.ErrorInterceptor(w, r, err, body.ID)
		return
	}

	successinterceptor.SuccessInterceptor(
		w, r,
		http.StatusOK,
		"playlist berhasil dihapus",
		nil,
		body.ID,
	)
}

// UTIL
func parseInt(str string) int {
	val, _ := strconv.Atoi(str)
	return val
}
