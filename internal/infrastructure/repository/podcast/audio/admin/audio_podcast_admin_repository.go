package audioadminrepository

import (
	"sync"
	"fmt" // ← INI YANG KURANG

	audioentities "mqfm_backend/internal/domain/entities/podcast/audio/admin"

)

// ============================================================
// 🔥 INTERFACE
// ============================================================
type AudioAdminPodcastRepository interface {
	Save(item *audioentities.AudioAdminPodcast) error
	FindByID(id int) (*audioentities.AudioAdminPodcast, bool)
	FindByTitle(title string) (*audioentities.AudioAdminPodcast, bool)
	GetAll() []*audioentities.AudioAdminPodcast
	Update(item *audioentities.AudioAdminPodcast) error
	Delete(id int) error
	Count() int
	Clear()
}

// ============================================================
// 🔥 MEMORY IMPLEMENTATION (SAFE UNTUK PRODUCTION KECIL / TEST)
// ============================================================
type audioAdminPodcastRepo struct {
	data map[int]*audioentities.AudioAdminPodcast
	mu   sync.RWMutex
	auto int
}

// ============================================================
// 🔥 CONSTRUCTOR
// ============================================================
func NewAudioAdminPodcastRepository() AudioAdminPodcastRepository {
	return &audioAdminPodcastRepo{
		data: make(map[int]*audioentities.AudioAdminPodcast),
		auto: 1,
	}
}

// ============================================================
// 🔥 SAVE
// ============================================================
func (r *audioAdminPodcastRepo) Save(item *audioentities.AudioAdminPodcast) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	item.ID = r.auto
	r.data[r.auto] = item
	r.auto++

	return nil
}

// ============================================================
// 🔥 FIND BY ID
// ============================================================
func (r *audioAdminPodcastRepo) FindByID(id int) (*audioentities.AudioAdminPodcast, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	item, ok := r.data[id]
	return item, ok
}

// ============================================================
// 🔥 FIND BY TITLE
// ============================================================
func (r *audioAdminPodcastRepo) FindByTitle(title string) (*audioentities.AudioAdminPodcast, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, v := range r.data {
		if v.Title == title {
			return v, true
		}
	}
	return nil, false
}

// ============================================================
// 🔥 GET ALL
// ============================================================
func (r *audioAdminPodcastRepo) GetAll() []*audioentities.AudioAdminPodcast {
	r.mu.RLock()
	defer r.mu.RUnlock()

	list := make([]*audioentities.AudioAdminPodcast, 0, len(r.data))
	for _, v := range r.data {
		list = append(list, v)
	}
	return list
}

// ============================================================
// 🔥 UPDATE
// ============================================================
func (r *audioAdminPodcastRepo) Update(item *audioentities.AudioAdminPodcast) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.data[item.ID]; !ok {
		return fmt.Errorf("podcast audio tidak ditemukan")
	}

	r.data[item.ID] = item
	return nil
}

// ============================================================
// 🔥 DELETE
// ============================================================
func (r *audioAdminPodcastRepo) Delete(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.data[id]; !ok {
		return fmt.Errorf("podcast audio tidak ditemukan")
	}

	delete(r.data, id)
	return nil
}

// ============================================================
// 🔥 COUNT
// ============================================================
func (r *audioAdminPodcastRepo) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.data)
}

// ============================================================
// 🔥 CLEAR (untuk testing)
// ============================================================
func (r *audioAdminPodcastRepo) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.data = make(map[int]*audioentities.AudioAdminPodcast)
	r.auto = 1
}
