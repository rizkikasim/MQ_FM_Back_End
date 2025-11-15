package playlistrepository

import (
	"fmt"
	"sync"

	playlistentity "mqfm_backend/internal/domain/entities/playlist"

)

//
// ============================================================
// 🔥 INTERFACE — Playlist Repository
// ============================================================
type PlaylistRepository interface {
	Save(item *playlistentity.Playlist) error
	FindByID(id int) (*playlistentity.Playlist, bool)
	FindByUser(userID int) []*playlistentity.Playlist
	GetAll() []*playlistentity.Playlist
	Update(item *playlistentity.Playlist) error
	Delete(id int) error
	Count() int
	Clear()
}

//
// ============================================================
// 🔥 MEMORY IMPLEMENTATION (AMAN BUAT DEV/TEST KECIL)
// ============================================================
type playlistRepo struct {
	data map[int]*playlistentity.Playlist
	mu   sync.RWMutex
	auto int
}

//
// ============================================================
// 🔥 CONSTRUCTOR
// ============================================================
func NewPlaylistRepository() PlaylistRepository {
	return &playlistRepo{
		data: make(map[int]*playlistentity.Playlist),
		auto: 1,
	}
}

//
// ============================================================
// 🔥 SAVE (CREATE)
// ============================================================
func (r *playlistRepo) Save(item *playlistentity.Playlist) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	item.ID = r.auto
	r.data[r.auto] = item
	r.auto++

	return nil
}

//
// ============================================================
// 🔥 FIND BY ID
// ============================================================
func (r *playlistRepo) FindByID(id int) (*playlistentity.Playlist, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	item, ok := r.data[id]
	return item, ok
}

//
// ============================================================
// 🔥 FIND BY USER_ID
// ============================================================
func (r *playlistRepo) FindByUser(userID int) []*playlistentity.Playlist {
	r.mu.RLock()
	defer r.mu.RUnlock()

	list := []*playlistentity.Playlist{}
	for _, v := range r.data {
		if v.UserID == userID {
			list = append(list, v)
		}
	}
	return list
}

//
// ============================================================
// 🔥 GET ALL
// ============================================================
func (r *playlistRepo) GetAll() []*playlistentity.Playlist {
	r.mu.RLock()
	defer r.mu.RUnlock()

	list := make([]*playlistentity.Playlist, 0, len(r.data))
	for _, v := range r.data {
		list = append(list, v)
	}
	return list
}

//
// ============================================================
// 🔥 UPDATE
// ============================================================
func (r *playlistRepo) Update(item *playlistentity.Playlist) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.data[item.ID]; !ok {
		return fmt.Errorf("playlist tidak ditemukan")
	}

	r.data[item.ID] = item
	return nil
}

//
// ============================================================
// 🔥 DELETE
// ============================================================
func (r *playlistRepo) Delete(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.data[id]; !ok {
		return fmt.Errorf("playlist tidak ditemukan")
	}

	delete(r.data, id)
	return nil
}

//
// ============================================================
// 🔥 COUNT
// ============================================================
func (r *playlistRepo) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return len(r.data)
}

//
// ============================================================
// 🔥 CLEAR (RESET – buat testing)
// ============================================================
func (r *playlistRepo) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.data = make(map[int]*playlistentity.Playlist)
	r.auto = 1
}
