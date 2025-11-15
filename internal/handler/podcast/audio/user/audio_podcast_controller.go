package audiousercontroller

import (
	"errors"
	"net/http"
	"strconv"

	audiouserusecase "mqfm_backend/internal/domain/usecases/podcast/audio/user"
	errorinterceptor "mqfm_backend/internal/presentation/interceptor/errors"
	successinterceptor "mqfm_backend/internal/presentation/interceptor/success"

)

type AudioUserPodcastController struct {
	usecase *audiouserusecase.AudioUserPodcastUsecase
}

func NewAudioUserPodcastController(usecase *audiouserusecase.AudioUserPodcastUsecase) *AudioUserPodcastController {
	return &AudioUserPodcastController{usecase: usecase}
}

//
// ============================================================
// 🔥 GET ALL PODCAST AUDIO (USER)
// ============================================================
func (c *AudioUserPodcastController) GetAll(w http.ResponseWriter, r *http.Request) {
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
// 🔥 GET PODCAST AUDIO BY ID (USER)
// ============================================================
func (c *AudioUserPodcastController) GetByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorinterceptor.ErrorInterceptor(w, r, http.ErrNotSupported, http.StatusMethodNotAllowed, 0)
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		errorinterceptor.ErrorInterceptor(w, r, errors.New("id wajib diisi"), http.StatusBadRequest, 0)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		errorinterceptor.ErrorInterceptor(w, r, errors.New("id tidak valid"), http.StatusBadRequest, 0)
		return
	}

	item, err := c.usecase.GetByID(id)
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
