-- +goose Up
-- +goose StatementBegin
ALTER TABLE rooms ADD COLUMN is_pinned BOOLEAN DEFAULT FALSE;
ALTER TABLE rooms ADD COLUMN topic_title TEXT;
ALTER TABLE rooms ADD COLUMN topic_description TEXT;
ALTER TABLE rooms ADD COLUMN topic_url TEXT;
ALTER TABLE rooms ADD COLUMN topic_source VARCHAR(50);
ALTER TABLE rooms ADD COLUMN topic_updated_at TIMESTAMP WITH TIME ZONE;

-- Create index for pinned rooms
CREATE INDEX idx_rooms_is_pinned ON rooms(is_pinned) WHERE is_pinned = TRUE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_rooms_is_pinned;
ALTER TABLE rooms DROP COLUMN topic_updated_at;
ALTER TABLE rooms DROP COLUMN topic_source;
ALTER TABLE rooms DROP COLUMN topic_url;
ALTER TABLE rooms DROP COLUMN topic_description;
ALTER TABLE rooms DROP COLUMN topic_title;
ALTER TABLE rooms DROP COLUMN is_pinned;
-- +goose StatementEnd