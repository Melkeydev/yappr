-- +goose Up
-- +goose StatementBegin
DO $$
BEGIN
    -- Add creator_id column if it doesn't exist
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
                   WHERE table_name = 'rooms' AND column_name = 'creator_id') THEN
        ALTER TABLE rooms ADD COLUMN creator_id UUID REFERENCES users(id) ON DELETE SET NULL;
    END IF;
    
    -- Create index if it doesn't exist
    IF NOT EXISTS (SELECT 1 FROM pg_indexes 
                   WHERE tablename = 'rooms' AND indexname = 'idx_rooms_creator_id') THEN
        CREATE INDEX idx_rooms_creator_id ON rooms(creator_id);
    END IF;
END $$;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_rooms_creator_id;
ALTER TABLE rooms DROP COLUMN creator_id;
-- +goose StatementEnd