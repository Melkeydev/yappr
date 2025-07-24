package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	statsService "github.com/melkeydev/chat-go/internal/service/stats"
	"github.com/melkeydev/chat-go/util"
)

type StatsHandler struct {
	statsService *statsService.StatsService
}

func NewStatsHandler(statsService *statsService.StatsService) *StatsHandler {
	return &StatsHandler{
		statsService: statsService,
	}
}

// CheckIn handles daily check-in requests
func (h *StatsHandler) CheckIn(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	// Get user ID from context (requires JWT middleware)
	userIDStr, ok := ctx.Value("userID").(string)
	if !ok {
		util.WriteError(w, http.StatusUnauthorized, "user not authenticated")
		return
	}
	
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, "invalid user ID")
		return
	}
	
	log.Printf("Check-in request from user: %s", userIDStr)
	
	result, err := h.statsService.ProcessDailyCheckin(ctx, userID)
	if err != nil {
		log.Printf("Error processing check-in: %v", err)
		util.WriteError(w, http.StatusInternalServerError, "failed to process check-in")
		return
	}
	
	util.WriteJSON(w, http.StatusOK, result)
}

// GetUserProfile returns user profile information
func (h *StatsHandler) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	// Get target user ID from URL
	targetUserIDStr := chi.URLParam(r, "userId")
	targetUserID, err := uuid.Parse(targetUserIDStr)
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, "invalid user ID")
		return
	}
	
	// Get viewer ID from context (can be empty for anonymous users)
	viewerUserID := uuid.Nil
	if viewerIDStr, ok := ctx.Value("userID").(string); ok {
		if viewerID, err := uuid.Parse(viewerIDStr); err == nil {
			viewerUserID = viewerID
		}
	}
	
	profile, err := h.statsService.GetUserProfile(ctx, targetUserID, viewerUserID)
	if err != nil {
		log.Printf("Error getting user profile: %v", err)
		util.WriteError(w, http.StatusInternalServerError, "failed to get user profile")
		return
	}
	
	util.WriteJSON(w, http.StatusOK, profile)
}

// GiveUpvote handles upvote requests
func (h *StatsHandler) GiveUpvote(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	// Get user ID from context (requires JWT middleware)
	fromUserIDStr, ok := ctx.Value("userID").(string)
	if !ok {
		util.WriteError(w, http.StatusUnauthorized, "user not authenticated")
		return
	}
	
	fromUserID, err := uuid.Parse(fromUserIDStr)
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, "invalid user ID")
		return
	}
	
	// Parse request body
	var req struct {
		ToUserID string `json:"to_user_id"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.WriteError(w, http.StatusBadRequest, "invalid JSON payload")
		return
	}
	
	toUserID, err := uuid.Parse(req.ToUserID)
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, "invalid target user ID")
		return
	}
	
	log.Printf("Upvote request from %s to %s", fromUserIDStr, req.ToUserID)
	
	err = h.statsService.GiveUpvote(ctx, fromUserID, toUserID)
	if err != nil {
		// Handle specific stats errors
		if statsErr, ok := err.(*statsService.StatsError); ok {
			switch statsErr.Code {
			case "CANNOT_UPVOTE_SELF":
				util.WriteError(w, http.StatusBadRequest, statsErr.Message)
			case "UPVOTE_NOT_ALLOWED":
				util.WriteError(w, http.StatusConflict, statsErr.Message)
			default:
				util.WriteError(w, http.StatusInternalServerError, "failed to process upvote")
			}
			return
		}
		
		log.Printf("Error processing upvote: %v", err)
		util.WriteError(w, http.StatusInternalServerError, "failed to process upvote")
		return
	}
	
	util.WriteJSON(w, http.StatusOK, map[string]string{"status": "success"})
}