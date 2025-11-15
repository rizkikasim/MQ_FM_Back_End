package playlistrepository

import (
	"sync"

	playlistentity "mqfm_backend/internal/domain/entities/playlist"

)

//
// ============================================================
// 🔥 IN-MEMORY REPOSITORY UNTUK PLAYLIST (DEV/TEST)
// ============================================================
type InMemoryPlaylistRepository struct {
	items  map[int]*playlistentity.Playlist
	mu     sync.RWMutex
	nextID int
}

//
// ============================================================
// 🔥 CONSTRUCTOR
// ============================================================
func NewInMemoryPlaylistRepository() *InMemoryPlaylistRepository {
	return &InMemoryPlaylistRepository{
		items:  make(map[int]*playlistentity.Playlist),
		nextID: 1,
	}
}

//
// ============================================================
// 🔥 SAVE (CREATE)
// ============================================================
func (r *InMemoryPlaylistRepository) Save(item *playlistentity.Playlist) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// auto ID
	if item.ID == 0 {
		item.ID = r.nextID
		r.nextID++
	}

	r.items[item.ID] = item
	return nil
}

//
// ============================================================
// 🔥 FIND BY ID
// ============================================================
func (r *InMemoryPlaylistRepository) FindByID(id int) (*playlistentity.Playlist, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	val, exists := r.items[id]
	return val, exists
}

//
// ============================================================
// 🔥 FIND BY USER_ID
// ============================================================
func (r *InMemoryPlaylistRepository) FindByUserID(userID int) []*playlistentity.Playlist {
	r.mu.RLock()
	defer r.mu.RUnlock()

	list := []*playlistentity.Playlist{}
	for _, v := range r.items {
		if v.UserID == userID {
			list = append(list, v)
		}
	}
	return list
}

//
// ============================================================
// 🔥 UPDATE
// ============================================================
func (r *InMemoryPlaylistRepository) Update(item *playlistentity.Playlist) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.items[item.ID]; !exists {
		return nil // atau return error, terserah
	}

	r.items[item.ID] = item
	return nil
}

//
// ============================================================
// 🔥 DELETE
// ============================================================
func (r *InMemoryPlaylistRepository) Delete(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.items, id)
	return nil
}

//
// ============================================================
// 🔥 GET ALL
// ============================================================
func (r *InMemoryPlaylistRepository) GetAll() []*playlistentity.Playlist {
	r.mu.RLock()
	defer r.mu.RUnlock()

	list := make([]*playlistentity.Playlist, 0, len(r.items))
	for _, v := range r.items {
		list = append(list, v)
	}
	return list
}

//
// ============================================================
// 🔥 COUNT
// ============================================================
func (r *InMemoryPlaylistRepository) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return len(r.items)
}

//
// ============================================================
// 🔥 CLEAR (Testing / Reset Data)
// ============================================================
func (r *InMemoryPlaylistRepository) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.items = make(map[int]*playlistentity.Playlist)
	r.nextID = 1
}
