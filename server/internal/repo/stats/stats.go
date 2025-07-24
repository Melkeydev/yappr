package stats

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/google/uuid"
)

type UserStats struct {
	UserID                uuid.UUID  `db:"user_id"`
	DailyStreak           int        `db:"daily_streak"`
	TotalCheckins         int        `db:"total_checkins"`
	TotalMessages         int        `db:"total_messages"`
	TotalUpvotesReceived  int        `db:"total_upvotes_received"`
	LastCheckinDate       *time.Time `db:"last_checkin_date"`
	LastUpvoteGivenDate   *time.Time `db:"last_upvote_given_date"`
	CreatedAt             time.Time  `db:"created_at"`
	UpdatedAt             time.Time  `db:"updated_at"`
}

type DailyCheckin struct {
	ID           uuid.UUID `db:"id"`
	UserID       uuid.UUID `db:"user_id"`
	CheckinDate  time.Time `db:"checkin_date"`
	StreakCount  int       `db:"streak_count"`
	CreatedAt    time.Time `db:"created_at"`
}

type Upvote struct {
	ID         uuid.UUID `db:"id"`
	FromUserID uuid.UUID `db:"from_user_id"`
	ToUserID   uuid.UUID `db:"to_user_id"`
	CreatedAt  time.Time `db:"created_at"`
}

type Achievement struct {
	ID            uuid.UUID `db:"id"`
	Name          string    `db:"name"`
	Description   string    `db:"description"`
	Icon          string    `db:"icon"`
	ThresholdType string    `db:"threshold_type"`
	ThresholdValue int      `db:"threshold_value"`
	EarnedAt      *time.Time `db:"earned_at,omitempty"`
}

type StatsRepository struct {
	db *sql.DB
}

func NewStatsRepository(db *sql.DB) *StatsRepository {
	return &StatsRepository{
		db: db,
	}
}

// GetOrCreateUserStats gets existing stats or creates new entry
func (r *StatsRepository) GetOrCreateUserStats(ctx context.Context, userID uuid.UUID) (*UserStats, error) {
	stats := &UserStats{}
	
	// Try to get existing stats
	query := `
		SELECT user_id, daily_streak, total_checkins, total_messages, 
		       total_upvotes_received, last_checkin_date, last_upvote_given_date,
		       created_at, updated_at
		FROM user_stats 
		WHERE user_id = $1`
	
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&stats.UserID, &stats.DailyStreak, &stats.TotalCheckins,
		&stats.TotalMessages, &stats.TotalUpvotesReceived,
		&stats.LastCheckinDate, &stats.LastUpvoteGivenDate,
		&stats.CreatedAt, &stats.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		// Create new stats entry
		insertQuery := `
			INSERT INTO user_stats (user_id) 
			VALUES ($1) 
			RETURNING user_id, daily_streak, total_checkins, total_messages,
			         total_upvotes_received, last_checkin_date, last_upvote_given_date,
			         created_at, updated_at`
		
		err = r.db.QueryRowContext(ctx, insertQuery, userID).Scan(
			&stats.UserID, &stats.DailyStreak, &stats.TotalCheckins,
			&stats.TotalMessages, &stats.TotalUpvotesReceived,
			&stats.LastCheckinDate, &stats.LastUpvoteGivenDate,
			&stats.CreatedAt, &stats.UpdatedAt,
		)
	}
	
	return stats, err
}

// ProcessDailyCheckin handles check-in logic and returns new streak count
func (r *StatsRepository) ProcessDailyCheckin(ctx context.Context, userID uuid.UUID) (int, bool, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, false, err
	}
	defer tx.Rollback()

	today := time.Now().UTC().Truncate(24 * time.Hour)
	
	// Get current stats
	stats, err := r.GetOrCreateUserStats(ctx, userID)
	if err != nil {
		return 0, false, err
	}
	
	// Check if already checked in today
	if stats.LastCheckinDate != nil {
		lastCheckin := stats.LastCheckinDate.UTC().Truncate(24 * time.Hour)
		if lastCheckin.Equal(today) {
			// Already checked in today
			return stats.DailyStreak, false, nil
		}
	}
	
	// Calculate new streak
	newStreak := 1
	if stats.LastCheckinDate != nil {
		yesterday := today.Add(-24 * time.Hour)
		lastCheckin := stats.LastCheckinDate.UTC().Truncate(24 * time.Hour)
		if lastCheckin.Equal(yesterday) {
			// Consecutive day
			newStreak = stats.DailyStreak + 1
		}
		// If more than 1 day gap, streak resets to 1
	}
	
	// Update user stats
	updateQuery := `
		UPDATE user_stats 
		SET daily_streak = $1, total_checkins = total_checkins + 1,
		    last_checkin_date = $2, updated_at = NOW()
		WHERE user_id = $3`
	
	_, err = tx.ExecContext(ctx, updateQuery, newStreak, today, userID)
	if err != nil {
		return 0, false, err
	}
	
	// Insert checkin record
	insertQuery := `
		INSERT INTO daily_checkins (user_id, checkin_date, streak_count)
		VALUES ($1, $2, $3)`
	
	_, err = tx.ExecContext(ctx, insertQuery, userID, today, newStreak)
	if err != nil {
		return 0, false, err
	}
	
	err = tx.Commit()
	return newStreak, true, err
}

// CanUserUpvote checks if user can give an upvote (hasn't given one today and hasn't upvoted target)
func (r *StatsRepository) CanUserUpvote(ctx context.Context, fromUserID, toUserID uuid.UUID) (bool, error) {
	// Check if already upvoted this user
	existsQuery := `
		SELECT EXISTS(
			SELECT 1 FROM upvotes 
			WHERE from_user_id = $1 AND to_user_id = $2
		)`
	
	var alreadyUpvoted bool
	err := r.db.QueryRowContext(ctx, existsQuery, fromUserID, toUserID).Scan(&alreadyUpvoted)
	if err != nil {
		return false, err
	}
	
	if alreadyUpvoted {
		return false, nil
	}
	
	// Check if gave an upvote today
	today := time.Now().UTC().Truncate(24 * time.Hour)
	todayQuery := `
		SELECT last_upvote_given_date FROM user_stats 
		WHERE user_id = $1`
	
	var lastUpvoteDate *time.Time
	err = r.db.QueryRowContext(ctx, todayQuery, fromUserID).Scan(&lastUpvoteDate)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}
	
	if lastUpvoteDate != nil {
		lastUpvote := lastUpvoteDate.UTC().Truncate(24 * time.Hour)
		if lastUpvote.Equal(today) {
			return false, nil // Already gave upvote today
		}
	}
	
	return true, nil
}

// GiveUpvote processes an upvote between users
func (r *StatsRepository) GiveUpvote(ctx context.Context, fromUserID, toUserID uuid.UUID) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	
	today := time.Now().UTC().Truncate(24 * time.Hour)
	
	// Insert upvote record
	insertQuery := `
		INSERT INTO upvotes (from_user_id, to_user_id)
		VALUES ($1, $2)`
	
	_, err = tx.ExecContext(ctx, insertQuery, fromUserID, toUserID)
	if err != nil {
		return err
	}
	
	// Update giver's last upvote date
	updateGiverQuery := `
		UPDATE user_stats 
		SET last_upvote_given_date = $1, updated_at = NOW()
		WHERE user_id = $2`
	
	_, err = tx.ExecContext(ctx, updateGiverQuery, today, fromUserID)
	if err != nil {
		return err
	}
	
	// Update receiver's upvote count
	updateReceiverQuery := `
		UPDATE user_stats 
		SET total_upvotes_received = total_upvotes_received + 1, updated_at = NOW()
		WHERE user_id = $1`
	
	_, err = tx.ExecContext(ctx, updateReceiverQuery, toUserID)
	if err != nil {
		return err
	}
	
	return tx.Commit()
}

// GetUserProfile returns user stats for profile display
func (r *StatsRepository) GetUserProfile(ctx context.Context, userID uuid.UUID) (*UserStats, error) {
	return r.GetOrCreateUserStats(ctx, userID)
}

// IncrementMessageCount increments the user's total message count
func (r *StatsRepository) IncrementMessageCount(ctx context.Context, userID uuid.UUID) error {
	query := `
		UPDATE user_stats 
		SET total_messages = total_messages + 1, updated_at = NOW()
		WHERE user_id = $1`
	
	result, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	// If no rows affected, create the user stats entry first
	if rowsAffected == 0 {
		_, err := r.GetOrCreateUserStats(ctx, userID)
		if err != nil {
			return err
		}
		
		// Try the update again
		_, err = r.db.ExecContext(ctx, query, userID)
		return err
	}
	
	return nil
}

// CheckAndAwardAchievements checks if user has earned new achievements and awards them
func (r *StatsRepository) CheckAndAwardAchievements(ctx context.Context, userID uuid.UUID) ([]Achievement, error) {
	// Get user's current stats
	stats, err := r.GetOrCreateUserStats(ctx, userID)
	if err != nil {
		return nil, err
	}
	
	// Get all achievement types
	achievementTypes, err := r.getAllAchievementTypes(ctx)
	if err != nil {
		return nil, err
	}
	
	// Get user's already earned achievements
	earnedAchievements, err := r.getUserAchievements(ctx, userID)
	if err != nil {
		return nil, err
	}
	
	// Check for new achievements
	newAchievements := []Achievement{}
	
	for _, achType := range achievementTypes {
		// Skip if already earned
		if r.hasAchievement(earnedAchievements, achType.ID) {
			continue
		}
		
		// Check if threshold is met
		var currentValue int
		switch achType.ThresholdType {
		case "streak":
			currentValue = stats.DailyStreak
		case "messages":
			currentValue = stats.TotalMessages
		case "upvotes":
			currentValue = stats.TotalUpvotesReceived
		default:
			continue
		}
		
		if currentValue >= achType.ThresholdValue {
			// Award the achievement
			if err := r.awardAchievement(ctx, userID, achType.ID); err != nil {
				log.Printf("Failed to award achievement %s to user %s: %v", achType.Name, userID.String(), err)
				continue
			}
			
			achievement := Achievement{
				ID:            achType.ID,
				Name:          achType.Name,
				Description:   achType.Description,
				Icon:          achType.Icon,
				ThresholdType: achType.ThresholdType,
				ThresholdValue: achType.ThresholdValue,
			}
			newAchievements = append(newAchievements, achievement)
		}
	}
	
	return newAchievements, nil
}

// getAllAchievementTypes gets all available achievement types
func (r *StatsRepository) getAllAchievementTypes(ctx context.Context) ([]Achievement, error) {
	query := `
		SELECT id, name, description, icon, threshold_type, threshold_value
		FROM achievement_types
		ORDER BY threshold_value ASC`
	
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var achievements []Achievement
	for rows.Next() {
		var ach Achievement
		err := rows.Scan(&ach.ID, &ach.Name, &ach.Description, &ach.Icon, 
			&ach.ThresholdType, &ach.ThresholdValue)
		if err != nil {
			return nil, err
		}
		achievements = append(achievements, ach)
	}
	
	return achievements, rows.Err()
}

// getUserAchievements gets user's earned achievements
func (r *StatsRepository) getUserAchievements(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	query := `
		SELECT achievement_type_id 
		FROM user_achievements 
		WHERE user_id = $1`
	
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var achievementIDs []uuid.UUID
	for rows.Next() {
		var achID uuid.UUID
		if err := rows.Scan(&achID); err != nil {
			return nil, err
		}
		achievementIDs = append(achievementIDs, achID)
	}
	
	return achievementIDs, rows.Err()
}

// hasAchievement checks if user already has an achievement
func (r *StatsRepository) hasAchievement(earnedAchievements []uuid.UUID, achievementID uuid.UUID) bool {
	for _, achID := range earnedAchievements {
		if achID == achievementID {
			return true
		}
	}
	return false
}

// awardAchievement awards an achievement to a user
func (r *StatsRepository) awardAchievement(ctx context.Context, userID, achievementID uuid.UUID) error {
	query := `
		INSERT INTO user_achievements (user_id, achievement_type_id)
		VALUES ($1, $2)
		ON CONFLICT (user_id, achievement_type_id) DO NOTHING`
	
	_, err := r.db.ExecContext(ctx, query, userID, achievementID)
	return err
}

// GetUserAchievementsWithDetails gets user's achievements with full details
func (r *StatsRepository) GetUserAchievementsWithDetails(ctx context.Context, userID uuid.UUID) ([]Achievement, error) {
	query := `
		SELECT at.id, at.name, at.description, at.icon, at.threshold_type, at.threshold_value, ua.earned_at
		FROM user_achievements ua
		JOIN achievement_types at ON ua.achievement_type_id = at.id
		WHERE ua.user_id = $1
		ORDER BY ua.earned_at DESC`
	
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var achievements []Achievement
	for rows.Next() {
		var ach Achievement
		err := rows.Scan(&ach.ID, &ach.Name, &ach.Description, &ach.Icon,
			&ach.ThresholdType, &ach.ThresholdValue, &ach.EarnedAt)
		if err != nil {
			return nil, err
		}
		achievements = append(achievements, ach)
	}
	
	return achievements, rows.Err()
}