package playlistusecase

import (
	playlistentity "mqfm_backend/internal/domain/entities/playlist"
	playlistservice "mqfm_backend/internal/application/services/playlist"
	req "mqfm_backend/internal/domain/request/playlist"

)

// ------------------------------------------------------------
// 🔥 USECASE STRUCT (PLAYLIST)
// ------------------------------------------------------------
type PlaylistUsecase struct {
	service *playlistservice.PlaylistService
}

// ------------------------------------------------------------
// 🔥 CONSTRUCTOR
// ------------------------------------------------------------
func NewPlaylistUsecase(service *playlistservice.PlaylistService) *PlaylistUsecase {
	return &PlaylistUsecase{service: service}
}

// ============================================================
// 🔥 CREATE PLAYLIST
// ============================================================
func (u *PlaylistUsecase) Create(request req.CreatePlaylistRequest) (*playlistentity.Playlist, error) {

	item, err := u.service.Create(
		request.Title,
		request.Description,
		request.Thumbnail,
		request.UserID,
	)
	if err != nil {
		return nil, err
	}

	return item, nil
}

// ============================================================
// 🔥 GET ALL PLAYLIST BY USER
// ============================================================
func (u *PlaylistUsecase) GetAllByUser(userID int) ([]*playlistentity.Playlist, error) {
	return u.service.GetAllByUser(userID)
}

// ============================================================
// 🔥 GET PLAYLIST DETAIL
// ============================================================
func (u *PlaylistUsecase) GetByID(request req.GetPlaylistRequest) (*playlistentity.Playlist, error) {
	return u.service.GetByID(request.ID)
}

// ============================================================
// 🔥 UPDATE PLAYLIST
// ============================================================
func (u *PlaylistUsecase) Update(request req.UpdatePlaylistRequest) (*playlistentity.Playlist, error) {

	item, err := u.service.Update(
		request.ID,
		request.NewTitle,
		request.NewDesc,
		request.NewThumb,
	)
	if err != nil {
		return nil, err
	}

	return item, nil
}

// ============================================================
// 🔥 DELETE PLAYLIST
// ============================================================
func (u *PlaylistUsecase) Delete(request req.DeletePlaylistRequest) error {
	return u.service.Delete(request.ID)
}
