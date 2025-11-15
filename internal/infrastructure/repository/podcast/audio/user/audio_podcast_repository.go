package audiouserrepository

import (
	"sync"

	audioentities "mqfm_backend/internal/domain/entities/podcast/audio/user"

)

//
// ============================================================
// 🔥 INTERFACE — khusus USER (READ ONLY)
// ============================================================
type AudioUserPodcastRepository interface {
	FindByID(id int) (*audioentities.AudioUserPodcast, bool)
	GetAll() []*audioentities.AudioUserPodcast
}

//
// ============================================================
// 🔥 IN-MEMORY IMPLEMENTATION (AMAN BUAT USER)
// ============================================================
type audioUserPodcastRepo struct {
	data map[int]*audioentities.AudioUserPodcast
	mu   sync.RWMutex
}

//
// ============================================================
// 🔥 CONSTRUCTOR
// ============================================================
func NewAudioUserPodcastRepository() AudioUserPodcastRepository {
	return &audioUserPodcastRepo{
		data: make(map[int]*audioentities.AudioUserPodcast),
	}
}

//
// ============================================================
// 🔥 FIND BY ID
// ============================================================
func (r *audioUserPodcastRepo) FindByID(id int) (*audioentities.AudioUserPodcast, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	item, ok := r.data[id]
	return item, ok
}

//
// ============================================================
// 🔥 GET ALL (USER READ ONLY)
// ============================================================
func (r *audioUserPodcastRepo) GetAll() []*audioentities.AudioUserPodcast {
	r.mu.RLock()
	defer r.mu.RUnlock()

	list := make([]*audioentities.AudioUserPodcast, 0, len(r.data))
	for _, v := range r.data {
		list = append(list, v)
	}
	return list
}
