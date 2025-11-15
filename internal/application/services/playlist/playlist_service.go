package playlistservice

import (
	"errors"
	"time"

	playlistentity "mqfm_backend/internal/domain/entities/playlist"
	playlistrepo "mqfm_backend/internal/infrastructure/repository/playlist"

)

// =====================================================
// 🔥 PLAYLIST SERVICE
// =====================================================
type PlaylistService struct {
	repo playlistrepo.PlaylistRepository
}

// -----------------------------------------------------
// 🔥 CONSTRUCTOR
// -----------------------------------------------------
func NewPlaylistService(repo playlistrepo.PlaylistRepository) *PlaylistService {
	return &PlaylistService{repo: repo}
}

// -----------------------------------------------------
// 🔥 CREATE PLAYLIST
// -----------------------------------------------------
func (s *PlaylistService) Create(
	title string,
	description string,
	thumbnail string,
	userID int,
) (*playlistentity.Playlist, error) {

	if title == "" {
		return nil, errors.New("judul playlist wajib diisi")
	}

	// Cek duplikasi playlist untuk user yg sama
	if existing, found := s.repo.FindByTitle(title, userID); found && existing != nil {
		return nil, errors.New("judul playlist sudah digunakan")
	}

	pl := playlistentity.NewPlaylist(
		title,
		description,
		thumbnail,
		userID,
	)

	pl.CreatedAt = time.Now()
	pl.UpdatedAt = time.Now()

	if err := s.repo.Save(pl); err != nil {
		return nil, err
	}

	return pl, nil
}

// -----------------------------------------------------
// 🔥 GET ALL (playlist user)
// -----------------------------------------------------
func (s *PlaylistService) GetAll(userID int) ([]*playlistentity.Playlist, error) {
	items := s.repo.GetAll(userID)

	if len(items) == 0 {
		return nil, errors.New("belum ada playlist")
	}

	return items, nil
}

// -----------------------------------------------------
// 🔥 GET BY ID
// -----------------------------------------------------
func (s *PlaylistService) GetByID(id int, userID int) (*playlistentity.Playlist, error) {
	item, found := s.repo.FindByID(id, userID)
	if !found {
		return nil, errors.New("playlist tidak ditemukan")
	}
	return item, nil
}

// -----------------------------------------------------
// 🔥 UPDATE PLAYLIST
// -----------------------------------------------------
func (s *PlaylistService) Update(
	id int,
	userID int,
	newTitle string,
	newDescription string,
	newThumbnail string,
) (*playlistentity.Playlist, error) {

	item, found := s.repo.FindByID(id, userID)
	if !found {
		return nil, errors.New("playlist tidak ditemukan")
	}

	if newTitle != "" {
		item.Title = newTitle
	}
	if newDescription != "" {
		item.Description = newDescription
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

// -----------------------------------------------------
// 🔥 DELETE
// -----------------------------------------------------
func (s *PlaylistService) Delete(id int, userID int) error {
	item, found := s.repo.FindByID(id, userID)
	if !found {
		return errors.New("playlist tidak ditemukan")
	}

	return s.repo.Delete(item.ID, userID)
}
