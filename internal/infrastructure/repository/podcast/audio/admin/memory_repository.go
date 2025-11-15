package audioadminrepository

import (
	"sync"

	audioentities "mqfm_backend/internal/domain/entities/podcast/audio/admin"

)

// ------------------------------------------------------------
// 🔥 IN-MEMORY REPOSITORY UNTUK AUDIO PODCAST ADMIN
// ------------------------------------------------------------
type InMemoryAudioAdminPodcastRepository struct {
	items map[int]*audioentities.AudioAdminPodcast
	mu    sync.RWMutex
	nextID int
}

// ------------------------------------------------------------
// 🔥 CONSTRUCTOR
// ------------------------------------------------------------
func NewInMemoryAudioAdminPodcastRepository() *InMemoryAudioAdminPodcastRepository {
	return &InMemoryAudioAdminPodcastRepository{
		items:  make(map[int]*audioentities.AudioAdminPodcast),
		nextID: 1,
	}
}

// ------------------------------------------------------------
// 🔥 SAVE (CREATE)
// ------------------------------------------------------------
func (r *InMemoryAudioAdminPodcastRepository) Save(item *audioentities.AudioAdminPodcast) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if item.ID == 0 {
		item.ID = r.nextID
		r.nextID++
	}

	r.items[item.ID] = item
	return nil
}

// ------------------------------------------------------------
// 🔥 FIND BY ID
// ------------------------------------------------------------
func (r *InMemoryAudioAdminPodcastRepository) FindByID(id int) (*audioentities.AudioAdminPodcast, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	val, exists := r.items[id]
	return val, exists
}

// ------------------------------------------------------------
// 🔥 FIND BY TITLE
// ------------------------------------------------------------
func (r *InMemoryAudioAdminPodcastRepository) FindByTitle(title string) (*audioentities.AudioAdminPodcast, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, v := range r.items {
		if v.Title == title {
			return v, true
		}
	}
	return nil, false
}

// ------------------------------------------------------------
// 🔥 UPDATE
// ------------------------------------------------------------
func (r *InMemoryAudioAdminPodcastRepository) Update(item *audioentities.AudioAdminPodcast) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.items[item.ID]; !exists {
		return nil
	}

	r.items[item.ID] = item
	return nil
}

// ------------------------------------------------------------
// 🔥 DELETE
// ------------------------------------------------------------
func (r *InMemoryAudioAdminPodcastRepository) Delete(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.items, id)
	return nil
}

// ------------------------------------------------------------
// 🔥 GET ALL
// ------------------------------------------------------------
func (r *InMemoryAudioAdminPodcastRepository) GetAll() []*audioentities.AudioAdminPodcast {
	r.mu.RLock()
	defer r.mu.RUnlock()

	list := make([]*audioentities.AudioAdminPodcast, 0, len(r.items))
	for _, v := range r.items {
		list = append(list, v)
	}
	return list
}

// ------------------------------------------------------------
// 🔥 COUNT
// ------------------------------------------------------------
func (r *InMemoryAudioAdminPodcastRepository) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return len(r.items)
}

// ------------------------------------------------------------
// 🔥 CLEAR (Testing / Reset Data)
// ------------------------------------------------------------
func (r *InMemoryAudioAdminPodcastRepository) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.items = make(map[int]*audioentities.AudioAdminPodcast)
	r.nextID = 1
}
