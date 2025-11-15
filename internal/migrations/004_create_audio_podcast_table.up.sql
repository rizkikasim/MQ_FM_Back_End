CREATE TABLE IF NOT EXISTS audio_admin_podcast (
    id INT AUTO_INCREMENT PRIMARY KEY,

    -- Judul podcast
    title VARCHAR(255) NOT NULL UNIQUE,

    -- Deskripsi podcast
    description TEXT,

    -- File audio yang sudah disimpan (lokal atau URL)
    audio_url VARCHAR(255),

    -- Durasi dalam detik
    duration INT DEFAULT 0,

    -- Relasi ke kategori (kategori admin)
    category_id INT NOT NULL,

    -- Thumbnail untuk cover podcast
    thumbnail VARCHAR(255),

    -- Timestamps
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
