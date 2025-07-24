package stats

import (
	"context"
	"log"

	"github.com/google/uuid"
	statsRepo "github.com/melkeydev/chat-go/internal/repo/stats"
)

type StatsService struct {
	statsRepo *statsRepo.StatsRepository
}

func NewStatsService(statsRepo *statsRepo.StatsRepository) *StatsService {
	return &StatsService{
		statsRepo: statsRepo,
	}
}

type CheckinResult struct {
	StreakCount int  `json:"streak_count"`
	IsNewCheckin bool `json:"is_new_checkin"`
}

type UserProfile struct {
	UserID               string `json:"user_id"`
	DailyStreak          int    `json:"daily_streak"`
	TotalCheckins        int    `json:"total_checkins"`
	TotalMessages        int    `json:"total_messages"`
	TotalUpvotesReceived int    `json:"total_upvotes_received"`
	CanReceiveUpvote     bool   `json:"can_receive_upvote"`
}

// ProcessDailyCheckin handles user check-in and returns streak info
func (s *StatsService) ProcessDailyCheckin(ctx context.Context, userID uuid.UUID) (*CheckinResult, error) {
	log.Printf("Processing daily check-in for user: %s", userID.String())
	
	streakCount, isNewCheckin, err := s.statsRepo.ProcessDailyCheckin(ctx, userID)
	if err != nil {
		log.Printf("Error processing check-in: %v", err)
		return nil, err
	}
	
	if isNewCheckin {
		log.Printf("New check-in recorded for user %s with streak: %d", userID.String(), streakCount)
	} else {
		log.Printf("User %s already checked in today, streak: %d", userID.String(), streakCount)
	}
	
	return &CheckinResult{
		StreakCount:  streakCount,
		IsNewCheckin: isNewCheckin,
	}, nil
}

// GetUserProfile returns user profile for display
func (s *StatsService) GetUserProfile(ctx context.Context, userID, viewerID uuid.UUID) (*UserProfile, error) {
	stats, err := s.statsRepo.GetUserProfile(ctx, userID)
	if err != nil {
		log.Printf("Error getting user profile: %v", err)
		return nil, err
	}
	
	// Check if viewer can upvote this user
	canUpvote := false
	if viewerID != userID { // Can't upvote yourself
		canUpvote, err = s.statsRepo.CanUserUpvote(ctx, viewerID, userID)
		if err != nil {
			log.Printf("Error checking upvote eligibility: %v", err)
			// Don't fail the whole request, just disable upvoting
			canUpvote = false
		}
	}
	
	return &UserProfile{
		UserID:               userID.String(),
		DailyStreak:          stats.DailyStreak,
		TotalCheckins:        stats.TotalCheckins,
		TotalMessages:        stats.TotalMessages,
		TotalUpvotesReceived: stats.TotalUpvotesReceived,
		CanReceiveUpvote:     canUpvote,
	}, nil
}

// GiveUpvote processes an upvote between users
func (s *StatsService) GiveUpvote(ctx context.Context, fromUserID, toUserID uuid.UUID) error {
	// Validate users are different
	if fromUserID == toUserID {
		return ErrCannotUpvoteSelf
	}
	
	// Check if upvote is allowed
	canUpvote, err := s.statsRepo.CanUserUpvote(ctx, fromUserID, toUserID)
	if err != nil {
		return err
	}
	
	if !canUpvote {
		return ErrUpvoteNotAllowed
	}
	
	log.Printf("Processing upvote from %s to %s", fromUserID.String(), toUserID.String())
	
	err = s.statsRepo.GiveUpvote(ctx, fromUserID, toUserID)
	if err != nil {
		log.Printf("Error giving upvote: %v", err)
		return err
	}
	
	log.Printf("Upvote successfully processed")
	return nil
}

// Custom errors
var (
	ErrCannotUpvoteSelf  = &StatsError{Code: "CANNOT_UPVOTE_SELF", Message: "Cannot upvote yourself"}
	ErrUpvoteNotAllowed  = &StatsError{Code: "UPVOTE_NOT_ALLOWED", Message: "Upvote not allowed - already upvoted this user or used daily upvote"}
)

type StatsError struct {
	Code    string
	Message string
}

func (e *StatsError) Error() string {
	return e.Message
}