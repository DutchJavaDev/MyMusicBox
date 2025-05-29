/*
 * MUSIC DATABASE SCHEMA
 * Features:
 * - Auto-incrementing IDs
 * - Automatic timestamp management
 * - Playlist song auto-positioning (using triggers)
 * - Full-text search support
 * - Proper indexing
 * - Data integrity constraints
 */

-- =============================================
-- SECTION 1: EXTENSIONS
-- =============================================

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- =============================================
-- SECTION 2: CORE TABLES
-- =============================================

CREATE TABLE Song (
    Id SERIAL PRIMARY KEY,
    Name VARCHAR(255) NOT NULL,
    SourceURL VARCHAR(512) NOT NULL,
    Path VARCHAR(512),
    Duration INTEGER,
    CreatedAt TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UpdatedAt TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT sourceurl_unique UNIQUE (SourceURL),
    CONSTRAINT path_unique UNIQUE (Path)
);

COMMENT ON COLUMN Song.SourceURL IS 'URL where the song can be obtained (not streaming URL)';
COMMENT ON COLUMN Song.Duration IS 'Duration in seconds';

CREATE TABLE Logging (
    Id SERIAL PRIMARY KEY,
    Timestamp TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    Message TEXT NOT NULL,
    Type INTEGER NOT NULL CHECK (Type BETWEEN 0 AND 2),
    Context JSONB
);

COMMENT ON COLUMN Logging.Type IS '0=Info, 1=Warning, 2=Error';

CREATE TABLE Playlist (
    Id SERIAL PRIMARY KEY,
    Name VARCHAR(255) NOT NULL,
    Description TEXT,
    CreationDate TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UpdatedAt TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    IsPublic BOOLEAN NOT NULL DEFAULT FALSE,
    CONSTRAINT name_unique UNIQUE (Name)
);

CREATE TABLE PlaylistSong (
    SongId INTEGER NOT NULL REFERENCES Song(Id) ON DELETE CASCADE,
    PlaylistId INTEGER NOT NULL REFERENCES Playlist(Id) ON DELETE CASCADE,
    Position INTEGER,
    AddedAt TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (SongId, PlaylistId)
);

-- =============================================
-- SECTION 3: AUTO-POSITIONING SYSTEM (TRIGGER-BASED)
-- =============================================

CREATE OR REPLACE FUNCTION set_playlist_position()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.Position IS NULL THEN
        NEW.Position := (
            SELECT COALESCE(MAX(Position), -1) + 1
            FROM PlaylistSong
            WHERE PlaylistId = NEW.PlaylistId
        );
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_playlist_position
BEFORE INSERT ON PlaylistSong
FOR EACH ROW EXECUTE FUNCTION set_playlist_position();

-- =============================================
-- SECTION 4: UPDATE HANDLING
-- =============================================

CREATE OR REPLACE FUNCTION update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.UpdatedAt = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Song update trigger
CREATE TRIGGER trg_song_update
BEFORE UPDATE ON Song
FOR EACH ROW EXECUTE FUNCTION update_timestamp();

-- Playlist update trigger
CREATE TRIGGER trg_playlist_update
BEFORE UPDATE ON Playlist
FOR EACH ROW EXECUTE FUNCTION update_timestamp();

-- =============================================
-- SECTION 5: INDEXES FOR PERFORMANCE
-- =============================================

CREATE INDEX idx_logging_timestamp ON Logging(Timestamp);
CREATE INDEX idx_logging_type ON Logging(Type);
CREATE INDEX idx_song_name ON Song(Name);
CREATE INDEX idx_song_created ON Song(CreatedAt);
CREATE INDEX idx_song_search ON Song USING gin(to_tsvector('english', Name));
CREATE INDEX idx_playlist_name ON Playlist(Name);
CREATE INDEX idx_playlist_public ON Playlist(IsPublic) WHERE IsPublic = TRUE;
CREATE INDEX idx_playlistsong_position ON PlaylistSong(PlaylistId, Position);
CREATE INDEX idx_playlistsong_playlist ON PlaylistSong(PlaylistId);
CREATE INDEX idx_playlistsong_song ON PlaylistSong(SongId);

-- =============================================
-- SECTION 6: SAMPLE DATA
-- =============================================

-- INSERT INTO Song (Name, SourceURL, Path, Duration) VALUES
-- ('Summer Vibes', 'https://example.com/tracks/summer123', '/path/one', 180),
-- ('Night Drive', 'https://example.com/tracks/night456', '/path/two', 240),
-- ('Morning Coffee', 'https://example.com/tracks/morning789', '/path/three', 210);

-- INSERT INTO Playlist (Name, IsPublic) VALUES
-- ('My Favorite Mix', FALSE),
-- ('Public Workout Jams', TRUE),
-- ('Chill Evening Tracks', FALSE);

-- -- These inserts will now work without recursion errors
-- INSERT INTO PlaylistSong (SongId, PlaylistId) VALUES
-- (1, 1),  -- Position 0
-- (2, 1),  -- Position 1
-- (3, 1);  -- Position 2

-- INSERT INTO PlaylistSong (SongId, PlaylistId) VALUES
-- (1, 2),  -- Position 0
-- (3, 2);  -- Position 1