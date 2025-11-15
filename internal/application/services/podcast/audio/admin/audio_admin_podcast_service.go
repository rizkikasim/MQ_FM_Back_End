package audiadminservice

import (
	"errors"
	"time"

	audioentities "mqfm_backend/internal/domain/entities/podcast/audio/admin"
	audiorepository "mqfm_backend/internal/infrastructure/repository/podcast/audio/admin"

)

// ---------------- SERVICE STRUCT ----------------
type AudioAdminPodcastService struct {
	repo audiorepository.AudioAdminPodcastRepository
}

// ---------------- CONSTRUCTOR ----------------
func NewAudioAdminPodcastService(repo audiorepository.AudioAdminPodcastRepository) *AudioAdminPodcastService {
	return &AudioAdminPodcastService{repo: repo}
}

// ---------------- CREATE ----------------
func (s *AudioAdminPodcastService) Create(
	title string,
	description string,
	audioURL string,
	duration int,
	categoryID int,
	thumbnail string,
) (*audioentities.AudioAdminPodcast, error) {

	// Cek duplikasi judul
	if existing, found := s.repo.FindByTitle(title); found && existing != nil {
		return nil, errors.New("judul podcast sudah digunakan")
	}

	podcast := audioentities.NewAudioAdminPodcast(
		title,
		description,
		audioURL,
		duration,
		categoryID,
		thumbnail,
	)

	podcast.CreatedAt = time.Now()
	podcast.UpdatedAt = time.Now()

	if err := s.repo.Save(podcast); err != nil {
		return nil, err
	}

	return podcast, nil
}

// ---------------- GET ALL ----------------
func (s *AudioAdminPodcastService) GetAll() ([]*audioentities.AudioAdminPodcast, error) {
	items := s.repo.GetAll()
	if len(items) == 0 {
		return nil, errors.New("belum ada podcast audio")
	}
	return items, nil
}

// ---------------- GET BY ID ----------------
func (s *AudioAdminPodcastService) GetByID(id int) (*audioentities.AudioAdminPodcast, error) {
	item, found := s.repo.FindByID(id)
	if !found {
		return nil, errors.New("podcast audio tidak ditemukan")
	}
	return item, nil
}

// ---------------- UPDATE ----------------
func (s *AudioAdminPodcastService) Update(
	id int,
	newTitle string,
	newDesc string,
	newAudioURL string,
	newDuration int,
	newCategoryID int,
	newThumbnail string,
) (*audioentities.AudioAdminPodcast, error) {

	item, found := s.repo.FindByID(id)
	if !found {
		return nil, errors.New("podcast audio tidak ditemukan")
	}

	if newTitle != "" {
		item.Title = newTitle
	}
	if newDesc != "" {
		item.Description = newDesc
	}
	if newAudioURL != "" {
		item.AudioURL = newAudioURL
	}
	if newDuration > 0 {
		item.Duration = newDuration
	}
	if newCategoryID > 0 {
		item.CategoryID = newCategoryID
	}
	if newThumbnail != "" {
		item.Thumbnail = newThumbnail
	}

	item.UpdatedAt = time.Now()

	if err := s.repo.Update(item); err != nil {
		return nil, err
	}

	return item, nil
}

// ---------------- DELETE ----------------
func (s *AudioAdminPodcastService) Delete(id int) error {
	item, found := s.repo.FindByID(id)
	if !found {
		return errors.New("podcast audio tidak ditemukan")
	}
	return s.repo.Delete(item.ID)
}
