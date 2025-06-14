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
    Name VARCHAR(100) NOT NULL, --- fixed lenght title
    SourceId VARCHAR(11) NOT NULL, --- youtube video id fixed lenght of 11
    Path VARCHAR(111) NOT NULL, --- for now this is fine
    ThumbnailPath VARCHAR(255) NOT NULL,
    Duration INTEGER NOT NULL,
    CreatedAt TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UpdatedAt TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT sourceurl_unique UNIQUE (SourceId),
    CONSTRAINT path_unique UNIQUE (Path),
    CONSTRAINT songsthumbnailpath_unique UNIQUE (ThumbnailPath)
);

COMMENT ON COLUMN Song.SourceId IS 'URL where the song can be obtained (not streaming URL)';
COMMENT ON COLUMN Song.Duration IS 'Duration in seconds';

-- Used for ytdlp output
CREATE TABLE TaskLog(
    Id SERIAL PRIMARY KEY,
    StartTime TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    EndTime TIMESTAMP WITH TIME ZONE NULL,
    Status INTEGER NOT NULL CHECK (Status BETWEEN 0 AND 4),
    OutputLog JSONB
);

COMMENT ON COLUMN TaskLog.Status IS '0=Pending, 1=Downloading, 2=Updating, 3=Done, 4=Error';

CREATE TABLE Playlist (
    Id SERIAL PRIMARY KEY,
    Name VARCHAR(255) NOT NULL,
    Description TEXT NOT NULL,
    ThumbnailPath VARCHAR(255) NOT NULL,
    CreationDate TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UpdatedAt TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    IsPublic BOOLEAN NOT NULL DEFAULT FALSE,
    CONSTRAINT name_unique UNIQUE (Name)
);

CREATE TABLE PlaylistSong (
    SongId INTEGER NOT NULL REFERENCES Song(Id) ON DELETE CASCADE,
    PlaylistId INTEGER NOT NULL REFERENCES Playlist(Id) ON DELETE CASCADE,
    Position INTEGER NOT NULL,
    AddedAt TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (PlaylistId, SongId)
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

CREATE INDEX idx_tasklog_starttime ON TaskLog(StartTime);
CREATE INDEX idx_tasklog_status ON TaskLog(Status);
CREATE INDEX idx_song_name ON Song(Name);
CREATE INDEX idx_song_created ON Song(CreatedAt);
CREATE INDEX idx_song_search ON Song USING gin(to_tsvector('english', Name));
CREATE INDEX idx_playlist_name ON Playlist(Name);
CREATE INDEX idx_playlist_public ON Playlist(IsPublic) WHERE IsPublic = TRUE;
CREATE INDEX idx_playlistsong_position ON PlaylistSong(PlaylistId, Position);
CREATE INDEX idx_playlistsong_playlist ON PlaylistSong(PlaylistId);
CREATE INDEX idx_playlistsong_song ON PlaylistSong(SongId);

-- =============================================
-- SECTION 6: Default DATA
-- =============================================

-- Main playlist for all songs
-- Should not show up in playlist list, instead under tab all songs 
INSERT INTO Playlist (Name, Description, ThumbnailPath, IsPublic) VALUES
('Library', 'Default playlist for all', 'https://picsum.photos/100', FALSE);

-- =============================================
-- SECTION 7: SAMPLE DATA
-- =============================================

-- INSERT INTO Song (Name, SourceURL, Path, Duration) VALUES
-- ('Summer Vibes', 'https://example.com/tracks/summer123', '/path/one', 180),
-- ('Night Drive', 'https://example.com/tracks/night456', '/path/two', 240),
-- ('Morning Coffee', 'https://example.com/tracks/morning789', '/path/three', 210);

-- INSERT INTO Playlist (Name, IsPublic) VALUES
-- ('My Favorite Mix', FALSE),
-- ('Public Workout Jams', TRUE),
-- ('Chill Evening Tracks', FALSE);

-- These inserts will now work without recursion errors
-- INSERT INTO PlaylistSong (SongId, PlaylistId) VALUES
-- (1, 1),  -- Position 0
-- (2, 1),  -- Position 1
-- (3, 1);  -- Position 2

-- INSERT INTO PlaylistSong (SongId, PlaylistId) VALUES
-- (1, 2),  -- Position 0
-- (3, 2);  -- Position 1