CREATE TABLE Migration (
    Id SERIAL PRIMARY KEY,
    FileName VARCHAR(75) NOT NULL,
    Contents TEXT,
    AppliedOn TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT filename_unique UNIQUE (FileName),
    CONSTRAINT contents_unique UNIQUE (Contents)
);