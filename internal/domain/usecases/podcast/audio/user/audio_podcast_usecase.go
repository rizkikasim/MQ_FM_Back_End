package audiouserusecase

import (
    audiouserentity "mqfm_backend/internal/domain/entities/podcast/audio/user"
    audiouservice "mqfm_backend/internal/application/services/podcast/audio/user"

)

// ======================================================
// 🔥 USECASE STRUCT (USER)
// ======================================================
type AudioUserPodcastUsecase struct {
    service *audiouservice.AudioUserPodcastService
}

// ======================================================
// 🔥 CONSTRUCTOR
// ======================================================
func NewAudioUserPodcastUsecase(service *audiouservice.AudioUserPodcastService) *AudioUserPodcastUsecase {
    return &AudioUserPodcastUsecase{service: service}
}

// ======================================================
// 🔥 GET ALL PODCAST (USER)
// ======================================================
func (u *AudioUserPodcastUsecase) GetAll() ([]*audiouserentity.AudioUserPodcast, error) {
    return u.service.GetAll()
}

// ======================================================
// 🔥 GET PODCAST BY ID (USER)
// ======================================================
func (u *AudioUserPodcastUsecase) GetByID(id int) (*audiouserentity.AudioUserPodcast, error) {
    return u.service.GetByID(id)
}
