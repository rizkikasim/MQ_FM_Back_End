package audioadmincontroller

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	audioadminusecase "mqfm_backend/internal/domain/usecases/podcast/audio/admin"
	domainrequest "mqfm_backend/internal/domain/request/podcast/audio/admin"
	errorinterceptor "mqfm_backend/internal/presentation/interceptor/errors"
	successinterceptor "mqfm_backend/internal/presentation/interceptor/success"
	validator "mqfm_backend/internal/presentation/validator"

)

// ------------------------------------------------------------
// Controller
// ------------------------------------------------------------
type AudioAdminPodcastController struct {
	usecase *audioadminusecase.AudioAdminPodcastUsecase
}

func NewAudioAdminPodcastController(usecase *audioadminusecase.AudioAdminPodcastUsecase) *AudioAdminPodcastController {
	return &AudioAdminPodcastController{usecase: usecase}
}

//
// ============================================================
// 🔥 CREATE AUDIO PODCAST
// ============================================================
func (c *AudioAdminPodcastController) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errorinterceptor.ErrorInterceptor(w, r, http.ErrNotSupported, http.StatusMethodNotAllowed, 0)
		return
	}

	if err := r.ParseMultipartForm(50 << 20); err != nil {
		errorinterceptor.ErrorInterceptor(w, r, errors.New("gagal membaca form data"), http.StatusBadRequest, 0)
		return
	}

	// FILE INPUT
	audio, audioHeader, _ := r.FormFile("audio_file")
	thumb, thumbHeader, _ := r.FormFile("thumbnail")

	req := domainrequest.CreateAudioAdminPodcastRequest{
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
		CategoryID:  parseInt(r.FormValue("category_id")),
		AudioFile:   audioHeader,
		Thumbnail:   thumbHeader,
	}

	if audio != nil {
		audio.Close()
	}
	if thumb != nil {
		thumb.Close()
	}

	// REQUIRED VALIDATION
	if err := validator.ValidateRequiredFields(map[string]string{
		"title":       req.Title,
		"category_id": r.FormValue("category_id"),
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
		"podcast audio berhasil dibuat",
		item,
		item.ID,
	)
}

//
// ============================================================
// 🔥 GET ALL
// ============================================================
func (c *AudioAdminPodcastController) GetAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorinterceptor.ErrorInterceptor(w, r, http.ErrNotSupported, http.StatusMethodNotAllowed, 0)
		return
	}

	items, err := c.usecase.GetAll()
	if err != nil {
		errorinterceptor.ErrorInterceptor(w, r, err, http.StatusInternalServerError, 0)
		return
	}

	successinterceptor.SuccessInterceptor(
		w, r,
		http.StatusOK,
		"daftar podcast audio berhasil diambil",
		items,
		len(items),
	)
}

//
// ============================================================
// 🔥 GET BY ID
// ============================================================
func (c *AudioAdminPodcastController) GetByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorinterceptor.ErrorInterceptor(w, r, http.ErrNotSupported, http.StatusMethodNotAllowed, 0)
		return
	}

	id := parseInt(r.URL.Query().Get("id"))
	if id == 0 {
		errorinterceptor.ErrorInterceptor(w, r, errors.New("id tidak valid"), http.StatusBadRequest, 0)
		return
	}

	item, err := c.usecase.GetByID(domainrequest.GetAudioAdminPodcastRequest{ID: id})
	if err != nil {
		errorinterceptor.ErrorInterceptor(w, r, err, http.StatusNotFound, id)
		return
	}

	successinterceptor.SuccessInterceptor(
		w, r,
		http.StatusOK,
		"podcast audio berhasil diambil",
		item,
		item.ID,
	)
}

//
// ============================================================
// 🔥 UPDATE
// ============================================================
func (c *AudioAdminPodcastController) Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errorinterceptor.ErrorInterceptor(w, r, http.ErrNotSupported, http.StatusMethodNotAllowed, 0)
		return
	}

	if err := r.ParseMultipartForm(50 << 20); err != nil {
		errorinterceptor.ErrorInterceptor(w, r, errors.New("gagal membaca form data"), http.StatusBadRequest, 0)
		return
	}

	id := parseInt(r.FormValue("id"))
	if id == 0 {
		errorinterceptor.ErrorInterceptor(w, r, errors.New("id tidak valid"), http.StatusBadRequest, 0)
		return
	}

	audio, audioHeader, _ := r.FormFile("new_audio_file")
	thumb, thumbHeader, _ := r.FormFile("new_thumbnail")

	req := domainrequest.UpdateAudioAdminPodcastRequest{
		ID:          id,
		NewTitle:    r.FormValue("new_title"),
		NewDesc:     r.FormValue("new_description"),
		NewCategory: parseInt(r.FormValue("new_category_id")),
		NewAudio:    audioHeader,
		NewThumb:    thumbHeader,
	}

	if audio != nil {
		audio.Close()
	}
	if thumb != nil {
		thumb.Close()
	}

	item, err := c.usecase.Update(req)
	if err != nil {
		errorinterceptor.ErrorInterceptor(w, r, err, http.StatusBadRequest, id)
		return
	}

	successinterceptor.SuccessInterceptor(
		w, r,
		http.StatusOK,
		"podcast audio berhasil diperbarui",
		item,
		item.ID,
	)
}

//
// ============================================================
// 🔥 DELETE
// ============================================================
func (c *AudioAdminPodcastController) Delete(w http.ResponseWriter, r *http.Request) {
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

	if err := c.usecase.Delete(domainrequest.DeleteAudioAdminPodcastRequest{ID: body.ID}); err != nil {
		errorinterceptor.ErrorInterceptor(w, r, err, http.StatusBadRequest, body.ID)
		return
	}

	successinterceptor.SuccessInterceptor(
		w, r,
		http.StatusOK,
		"podcast audio berhasil dihapus",
		nil,
		body.ID,
	)
}

//
// ============================================================
// 🔧 UTIL
// ============================================================
func parseInt(str string) int {
	val, _ := strconv.Atoi(str)
	return val
}
