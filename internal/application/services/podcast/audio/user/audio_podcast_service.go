package audiouservice

import (
    "errors"

    audiouserentity "mqfm_backend/internal/domain/entities/podcast/audio/user"
    audiouserrepository "mqfm_backend/internal/infrastructure/repository/podcast/audio/user"

)

//
// ===================================================
// 🔥 USER AUDIO PODCAST SERVICE (READ-ONLY)
// ===================================================
type AudioUserPodcastService struct {
    repo audiouserrepository.AudioUserPodcastRepository
}

//
// ---------------- CONSTRUCTOR ----------------
func NewAudioUserPodcastService(repo audiouserrepository.AudioUserPodcastRepository) *AudioUserPodcastService {
    return &AudioUserPodcastService{repo: repo}
}

//
// ---------------------------------------------------
// 🔥 GET ALL PODCAST (User baca dari table user)
// ---------------------------------------------------
func (s *AudioUserPodcastService) GetAll() ([]*audiouserentity.AudioUserPodcast, error) {
    items := s.repo.GetAll()

    if len(items) == 0 {
        return nil, errors.New("belum ada podcast audio")
    }

    return items, nil
}

//
// ---------------------------------------------------
// 🔥 GET BY ID (User baca dari table user)
// ---------------------------------------------------
func (s *AudioUserPodcastService) GetByID(id int) (*audiouserentity.AudioUserPodcast, error) {
    item, found := s.repo.FindByID(id)

    if !found {
        return nil, errors.New("podcast audio tidak ditemukan")
    }

    return item, nil
}
