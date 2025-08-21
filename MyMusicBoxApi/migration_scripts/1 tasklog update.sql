-- Used for ytdlp output
CREATE TABLE IF NOT EXISTS ParentTaskLog(
    Id SERIAL PRIMARY KEY,
    Url VARCHAR(72) NOT NULL,
    AddTime TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS ChildTaskLog(
    Id SERIAL PRIMARY KEY,
    ParentId INTEGER NOT NULL REFERENCES ParentTaskLog(Id) ON DELETE CASCADE,
    StartTime TIMESTAMP WITH TIME ZONE NULL,
    EndTime TIMESTAMP WITH TIME ZONE NULL,
    Status INTEGER NOT NULL CHECK (Status BETWEEN 0 AND 4),
    OutputLog JSONB
);

COMMENT ON COLUMN ChildTaskLog.Status IS '0=Pending, 1=Downloading, 2=Updating, 3=Done, 4=Error';

CREATE INDEX IF NOT EXISTS idx_parenttasklog_url ON ParentTaskLog(Url);
CREATE INDEX IF NOT EXISTS idx_parenttasklog_addtime ON ParentTaskLog(AddTime);

CREATE INDEX IF NOT EXISTS idx_childtasklog_starttime ON ChildTaskLog(StartTime);
CREATE INDEX IF NOT EXISTS idx_childtasklog_endtime ON ChildTaskLog(EndTime);
CREATE INDEX IF NOT EXISTS idx_childtasklog_status ON ChildTaskLog(Status);

--- Remove prev indexs
DROP INDEX IF EXISTS idx_tasklog_starttime;
DROP INDEX IF EXISTS idx_tasklog_status;

DROP TABLE IF EXISTS TaskLog;