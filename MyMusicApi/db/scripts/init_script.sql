/*
 * MUSIC DATABASE SCHEMA
 * Features:
 * - Auto-incrementing IDs
 * - Automatic timestamp management (without triggers)
 * - Playlist song auto-positioning
 * - Full-text search support
 * - Proper indexing
 * - Data integrity constraints
 */

-- =============================================
-- SECTION 1: EXTENSIONS
-- =============================================

-- Enable UUID generation capability
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Enable cryptographic functions
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- =============================================
-- SECTION 2: CORE TABLES
-- =============================================

/*
 * SONG TABLE
 * Stores metadata about individual songs
 */
CREATE TABLE Song (
    Id SERIAL PRIMARY KEY,
    Name VARCHAR(255) NOT NULL,
    SourceURL VARCHAR(512) NOT NULL,
    Path VARCHAR(512),
    Duration INTEGER,
    CreatedAt TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UpdatedAt TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT sourceurl_unique UNIQUE (SourceURL),  -- Prevent duplicate source URLs
    CONSTRAINT path_unique UNIQUE (Path) -- Prevent multiple songs having the same path
);

-- Column comments for documentation
COMMENT ON COLUMN Song.SourceURL IS 'URL where the song can be obtained (not streaming URL)';
COMMENT ON COLUMN Song.Duration IS 'Duration in seconds';

/*
 * LOGGING TABLE
 * Stores system events and messages
 */
CREATE TABLE Logging (
    Id SERIAL PRIMARY KEY,
    Timestamp TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    Message TEXT NOT NULL,
    Type SERIAL NOT NULL, -- 0 info, 1 warning, 2 error  -- Enforce valid log types
    Context JSONB  -- Store additional structured data
);

COMMENT ON COLUMN Logging.Context IS 'Additional structured logging data in JSON format';

/*
 * PLAYLIST TABLE
 * Stores playlist information
 */
CREATE TABLE Playlist (
    Id SERIAL PRIMARY KEY,
    Name VARCHAR(255) NOT NULL,
    Description TEXT,
    CreationDate TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UpdatedAt TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    IsPublic BOOLEAN NOT NULL DEFAULT FALSE,  -- Default to private playlists
    CONSTRAINT name_unique UNIQUE (Name) -- Prevent playlist with the same name
);

-- =============================================
-- SECTION 3: PLAYLIST SONG JUNCTION TABLE
-- =============================================

/*
 * PLAYLISTSONG TABLE
 * Manages many-to-many relationship between playlists and songs
 * with automatic position assignment
 */
CREATE TABLE PlaylistSong (
    SongId INTEGER NOT NULL REFERENCES Song(Id) ON DELETE CASCADE,
    PlaylistId INTEGER NOT NULL REFERENCES Playlist(Id) ON DELETE CASCADE,
    Position INTEGER NOT NULL,  -- Will be auto-set by the rule below
    AddedAt TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (SongId, PlaylistId)
);

-- =============================================
-- SECTION 4: AUTO-POSITIONING SYSTEM
-- =============================================

/*
 * GET_NEXT_POSITION FUNCTION
 * Helper function to calculate the next available position in a playlist
 */
CREATE OR REPLACE FUNCTION get_next_position(playlist_id INTEGER) 
RETURNS INTEGER AS $$
DECLARE
    max_pos INTEGER;
BEGIN
    -- Find current max position and increment by 1
    SELECT COALESCE(MAX(Position), -1) + 1 INTO max_pos 
    FROM PlaylistSong 
    WHERE PlaylistId = playlist_id;
    
    RETURN max_pos;
END;
$$ LANGUAGE plpgsql;

/*
 * AUTO_POSITION_PLAYLIST_SONG RULE
 * Automatically sets position for new entries when not specified
 */
CREATE OR REPLACE RULE auto_position_playlist_song AS
ON INSERT TO PlaylistSong
WHERE NEW.Position IS NULL
DO INSTEAD (
    INSERT INTO PlaylistSong (SongId, PlaylistId, Position, AddedAt)
    VALUES (
        NEW.SongId, 
        NEW.PlaylistId, 
        get_next_position(NEW.PlaylistId),  -- Auto-calculate position
        COALESCE(NEW.AddedAt, CURRENT_TIMESTAMP)  -- Use provided timestamp or current time
    )
);

-- =============================================
-- SECTION 5: UPDATE HANDLING (TRIGGER-FREE)
-- =============================================

/*
 * UPDATE HANDLING VIEWS
 * Provide automatic UpdatedAt maintenance without triggers
 */

-- Song update view
CREATE OR REPLACE VIEW SongWithUpdate AS
SELECT * FROM Song;

CREATE OR REPLACE RULE update_song AS
ON UPDATE TO SongWithUpdate DO INSTEAD
UPDATE Song SET
    Name = NEW.Name,
    SourceURL = NEW.SourceURL,
    Path = NEW.Path,
    Duration = NEW.Duration,
    UpdatedAt = CURRENT_TIMESTAMP  -- Auto-update timestamp
WHERE Id = OLD.Id;

-- Playlist update view
CREATE OR REPLACE VIEW PlaylistWithUpdate AS
SELECT * FROM Playlist;

CREATE OR REPLACE RULE update_playlist AS
ON UPDATE TO PlaylistWithUpdate DO INSTEAD
UPDATE Playlist SET
    Name = NEW.Name,
    Description = NEW.Description,
    IsPublic = NEW.IsPublic,
    UpdatedAt = CURRENT_TIMESTAMP  -- Auto-update timestamp
WHERE Id = OLD.Id;

-- =============================================
-- SECTION 6: INDEXES FOR PERFORMANCE
-- =============================================

-- Logging indexes
CREATE INDEX idx_logging_timestamp ON Logging(Timestamp);
CREATE INDEX idx_logging_type ON Logging(Type);

-- Song indexes
CREATE INDEX idx_song_name ON Song(Name);
CREATE INDEX idx_song_created ON Song(CreatedAt);
CREATE INDEX idx_song_search ON Song USING gin(to_tsvector('english', Name));  -- Full-text search

-- Playlist indexes
CREATE INDEX idx_playlist_name ON Playlist(Name);
CREATE INDEX idx_playlist_public ON Playlist(IsPublic) WHERE IsPublic = TRUE;  -- Partial index

-- PlaylistSong indexes
CREATE INDEX idx_playlistsong_position ON PlaylistSong(PlaylistId, Position);
CREATE INDEX idx_playlistsong_playlist ON PlaylistSong(PlaylistId);
CREATE INDEX idx_playlistsong_song ON PlaylistSong(SongId);

-- =============================================
-- SECTION 7: SAMPLE DATA (OPTIONAL)
-- =============================================

-- Sample songs
-- INSERT INTO Song (Name, SourceURL, Duration) VALUES
-- ('Summer Vibes', 'https://example.com/tracks/summer123', 180),
-- ('Night Drive', 'https://example.com/tracks/night456', 240),
-- ('Morning Coffee', 'https://example.com/tracks/morning789', 210);

-- -- Sample playlists
-- INSERT INTO Playlist (Name, IsPublic) VALUES
-- ('My Favorite Mix', FALSE),
-- ('Public Workout Jams', TRUE),
-- ('Chill Evening Tracks', FALSE);

-- -- Sample playlist entries (positions will be auto-assigned)
-- INSERT INTO PlaylistSong (SongId, PlaylistId) VALUES
-- (1, 1),  -- Position 0
-- (2, 1),  -- Position 1
-- (3, 1);  -- Position 2

-- INSERT INTO PlaylistSong (SongId, PlaylistId) VALUES
-- (1, 2),  -- Position 0
-- (3, 2);  -- Position 1