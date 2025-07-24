-- +goose Up
-- +goose StatementBegin
ALTER TABLE rooms ADD COLUMN creator_id UUID REFERENCES users(id) ON DELETE SET NULL;
CREATE INDEX idx_rooms_creator_id ON rooms(creator_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_rooms_creator_id;
ALTER TABLE rooms DROP COLUMN creator_id;
-- +goose StatementEnd