package audiouserrepository

import (
	"sync"

	audioentities "mqfm_backend/internal/domain/entities/podcast/audio/user"

)


// ------------------------------------------------------------
// 🔥 IN-MEMORY REPOSITORY UNTUK AUDIO PODCAST USER (READ ONLY)
// ------------------------------------------------------------
type InMemoryAudioUserPodcastRepository struct {
	items map[int]*audioentities.AudioUserPodcast
	mu    sync.RWMutex
}

// ------------------------------------------------------------
// 🔥 CONSTRUCTOR
// ------------------------------------------------------------
func NewInMemoryAudioUserPodcastRepository() *InMemoryAudioUserPodcastRepository {
	return &InMemoryAudioUserPodcastRepository{
		items: make(map[int]*audioentities.AudioUserPodcast),
	}
}

// ------------------------------------------------------------
// 🔥 FIND BY ID
// ------------------------------------------------------------
func (r *InMemoryAudioUserPodcastRepository) FindByID(id int) (*audioentities.AudioUserPodcast, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	val, exists := r.items[id]
	return val, exists
}

// ------------------------------------------------------------
// 🔥 GET ALL
// ------------------------------------------------------------
func (r *InMemoryAudioUserPodcastRepository) GetAll() []*audioentities.AudioUserPodcast {
	r.mu.RLock()
	defer r.mu.RUnlock()

	list := make([]*audioentities.AudioUserPodcast, 0, len(r.items))
	for _, v := range r.items {
		list = append(list, v)
	}
	return list
}
