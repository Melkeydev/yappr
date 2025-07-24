-- +goose Up
-- +goose StatementBegin

-- User statistics table
CREATE TABLE IF NOT EXISTS user_stats (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    daily_streak INTEGER NOT NULL DEFAULT 0,
    total_checkins INTEGER NOT NULL DEFAULT 0,
    total_messages INTEGER NOT NULL DEFAULT 0,
    total_upvotes_received INTEGER NOT NULL DEFAULT 0,
    last_checkin_date DATE,
    last_upvote_given_date DATE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Daily check-ins history
CREATE TABLE IF NOT EXISTS daily_checkins (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    checkin_date DATE NOT NULL,
    streak_count INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Upvotes between users
CREATE TABLE IF NOT EXISTS upvotes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    from_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    to_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    UNIQUE(from_user_id, to_user_id),
    CHECK(from_user_id != to_user_id)
);

-- Achievement types
CREATE TABLE IF NOT EXISTS achievement_types (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT NOT NULL,
    icon VARCHAR(50) NOT NULL,
    threshold_type VARCHAR(50) NOT NULL, -- 'streak', 'messages', 'upvotes'
    threshold_value INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- User achievements
CREATE TABLE IF NOT EXISTS user_achievements (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    achievement_type_id UUID NOT NULL REFERENCES achievement_types(id) ON DELETE CASCADE,
    earned_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, achievement_type_id)
);

-- Indexes for performance
CREATE INDEX idx_user_stats_user_id ON user_stats(user_id);
CREATE INDEX idx_daily_checkins_user_id ON daily_checkins(user_id);
CREATE INDEX idx_daily_checkins_date ON daily_checkins(checkin_date);
CREATE INDEX idx_upvotes_from_user ON upvotes(from_user_id);
CREATE INDEX idx_upvotes_to_user ON upvotes(to_user_id);
CREATE INDEX idx_user_achievements_user_id ON user_achievements(user_id);

-- Insert default achievement types
INSERT INTO achievement_types (name, description, icon, threshold_type, threshold_value) VALUES
('First Steps', 'Complete your first daily check-in', '🌟', 'streak', 1),
('Weekly Warrior', 'Maintain a 7-day streak', '🔥', 'streak', 7),
('Monthly Master', 'Maintain a 30-day streak', '👑', 'streak', 30),
('Chatter', 'Send your first 10 messages', '💬', 'messages', 10),
('Conversationalist', 'Send 100 messages', '🗣️', 'messages', 100),
('Popular', 'Receive your first 5 upvotes', '⭐', 'upvotes', 5),
('Beloved', 'Receive 25 upvotes', '💖', 'upvotes', 25);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_achievements;
DROP TABLE IF EXISTS achievement_types;
DROP TABLE IF EXISTS upvotes;
DROP TABLE IF EXISTS daily_checkins;
DROP TABLE IF EXISTS user_stats;
-- +goose StatementEnd