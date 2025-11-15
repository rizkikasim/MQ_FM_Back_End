package audioadminusecase

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	audioentities "mqfm_backend/internal/domain/entities/podcast/audio/admin"
	audioservice "mqfm_backend/internal/application/services/podcast/audio/admin"
	req "mqfm_backend/internal/domain/request/podcast/audio/admin"
	utils "mqfm_backend/internal/domain/utils"

)

// ---------------- USECASE STRUCT ----------------
type AudioAdminPodcastUsecase struct {
	service *audioservice.AudioAdminPodcastService
}

// ---------------- CONSTRUCTOR ----------------
func NewAudioAdminPodcastUsecase(service *audioservice.AudioAdminPodcastService) *AudioAdminPodcastUsecase {
	return &AudioAdminPodcastUsecase{service: service}
}

// =================================================
// 🔥 CREATE AUDIO PODCAST (AUTO DURASI)
// =================================================
func (u *AudioAdminPodcastUsecase) Create(request req.CreateAudioAdminPodcastRequest) (*audioentities.AudioAdminPodcast, error) {

	uploadDir := "storage/audio_podcast"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return nil, fmt.Errorf("gagal membuat folder upload: %v", err)
	}

	// ---------------------------------------------
	// 🔥 SIMPAN AUDIO FILE
	// ---------------------------------------------
	audioFilename := ""
	audioFullPath := ""

	if request.AudioFile != nil {

		audioFilename = fmt.Sprintf("audio_%s", request.AudioFile.Filename)
		audioFullPath = filepath.Join(uploadDir, audioFilename)

		src, err := request.AudioFile.Open()
		if err != nil {
			return nil, fmt.Errorf("gagal membuka file audio: %v", err)
		}
		defer src.Close()

		dst, err := os.Create(audioFullPath)
		if err != nil {
			return nil, fmt.Errorf("gagal menyimpan file audio: %v", err)
		}
		defer dst.Close()

		if _, err := io.Copy(dst, src); err != nil {
			return nil, fmt.Errorf("gagal menulis file audio: %v", err)
		}
	}

	// ---------------------------------------------
	// 🔥 AUTO DECODE DURASI MENGGUNAKAN FFPROBE
	// ---------------------------------------------
	duration := 0
	if audioFullPath != "" {
		dur, err := utils.GetAudioDuration(audioFullPath)
		if err != nil {
			return nil, fmt.Errorf("gagal mengambil durasi audio: %v", err)
		}
		duration = dur
	}

	// ---------------------------------------------
	// 🔥 SIMPAN THUMBNAIL
	// ---------------------------------------------
	thumbFilename := ""
	if request.Thumbnail != nil {

		thumbFilename = fmt.Sprintf("thumb_%s", request.Thumbnail.Filename)
		filePath := filepath.Join(uploadDir, thumbFilename)

		src, err := request.Thumbnail.Open()
		if err != nil {
			return nil, fmt.Errorf("gagal membuka file thumbnail: %v", err)
		}
		defer src.Close()

		dst, err := os.Create(filePath)
		if err != nil {
			return nil, fmt.Errorf("gagal menyimpan thumbnail: %v", err)
		}
		defer dst.Close()

		if _, err := io.Copy(dst, src); err != nil {
			return nil, fmt.Errorf("gagal menulis thumbnail: %v", err)
		}
	}

	// ---------------------------------------------
	// 🔥 CALL SERVICE CREATE
	// ---------------------------------------------
	item, err := u.service.Create(
		request.Title,
		request.Description,
		audioFilename,
		duration, // <-- AUTO duration, bukan dari frontend
		request.CategoryID,
		thumbFilename,
	)

	if err != nil {
		return nil, err
	}

	return item, nil
}

// =================================================
// 🔥 GET ALL
// =================================================
func (u *AudioAdminPodcastUsecase) GetAll() ([]*audioentities.AudioAdminPodcast, error) {
	return u.service.GetAll()
}

// =================================================
// 🔥 GET BY ID
// =================================================
func (u *AudioAdminPodcastUsecase) GetByID(request req.GetAudioAdminPodcastRequest) (*audioentities.AudioAdminPodcast, error) {
	return u.service.GetByID(request.ID)
}

// =================================================
// 🔥 UPDATE AUDIO PODCAST
// =================================================
func (u *AudioAdminPodcastUsecase) Update(request req.UpdateAudioAdminPodcastRequest) (*audioentities.AudioAdminPodcast, error) {

	uploadDir := "storage/audio_podcast"
	os.MkdirAll(uploadDir, os.ModePerm)

	var newAudioFilename string
	var newThumbFilename string
	var newAudioDuration int

	// ---------------------------------------------
	// 🔥 UPDATE AUDIO FILE (jika ada)
	// ---------------------------------------------
	if request.NewAudio != nil {

		newAudioFilename = fmt.Sprintf("audio_%d_%s", request.ID, request.NewAudio.Filename)
		filePath := filepath.Join(uploadDir, newAudioFilename)

		src, err := request.NewAudio.Open()
		if err != nil {
			return nil, fmt.Errorf("gagal membuka file audio update: %v", err)
		}
		defer src.Close()

		dst, err := os.Create(filePath)
		if err != nil {
			return nil, fmt.Errorf("gagal menyimpan file audio baru: %v", err)
		}
		defer dst.Close()

		if _, err := io.Copy(dst, src); err != nil {
			return nil, fmt.Errorf("gagal menulis file audio baru: %v", err)
		}

		// 🔥 AUTO DURASI BARU
		dur, err := utils.GetAudioDuration(filePath)
		if err == nil {
			newAudioDuration = dur
		}
	}

	// ---------------------------------------------
	// 🔥 UPDATE THUMBNAIL (jika ada)
	// ---------------------------------------------
	if request.NewThumb != nil {

		newThumbFilename = fmt.Sprintf("thumb_%d_%s", request.ID, request.NewThumb.Filename)
		filePath := filepath.Join(uploadDir, newThumbFilename)

		src, err := request.NewThumb.Open()
		if err != nil {
			return nil, fmt.Errorf("gagal membuka thumbnail update: %v", err)
		}
		defer src.Close()

		dst, err := os.Create(filePath)
		if err != nil {
			return nil, fmt.Errorf("gagal menyimpan thumbnail baru: %v", err)
		}
		defer dst.Close()

		if _, err := io.Copy(dst, src); err != nil {
			return nil, fmt.Errorf("gagal menulis thumbnail baru: %v", err)
		}
	}

	// ---------------------------------------------
	// 🔥 CALL SERVICE UPDATE
	// ---------------------------------------------
	item, err := u.service.Update(
		request.ID,
		request.NewTitle,
		request.NewDesc,
		newAudioFilename,
		newAudioDuration,
		request.NewCategory,
		newThumbFilename,
	)

	if err != nil {
		return nil, err
	}

	return item, nil
}

// =================================================
// 🔥 DELETE
// =================================================
func (u *AudioAdminPodcastUsecase) Delete(request req.DeleteAudioAdminPodcastRequest) error {
	return u.service.Delete(request.ID)
}
