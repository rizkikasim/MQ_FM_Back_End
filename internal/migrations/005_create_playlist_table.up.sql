CREATE TABLE IF NOT EXISTS playlists (
    id INT AUTO_INCREMENT PRIMARY KEY,

    -- Pemilik playlist (user)
    user_id INT NOT NULL,

    -- Nama playlist
    name VARCHAR(255) NOT NULL,

    -- Deskripsi playlist (optional)
    description TEXT,

    -- Thumbnail playlist (optional)
    thumbnail VARCHAR(255),

    -- Timestamp
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
