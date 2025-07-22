package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/melkeydev/chat-go/internal/api/model"
	"github.com/melkeydev/chat-go/internal/service/user"
	"github.com/melkeydev/chat-go/util"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req model.RequestCreateUser
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("CreateUser - JSON decode error: %v", err)
		util.WriteError(w, http.StatusBadRequest, "invalid JSON payload")
		return
	}

	log.Printf("CreateUser - Request received: username=%s, email=%s", req.Username, req.Email)

	user, err := h.userService.CreateUser(r.Context(), req)
	if err != nil {
		log.Printf("CreateUser - Service error: %v", err)
		util.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("CreateUser - Success: user created with ID=%s, username=%s", user.ID, user.Username)

	// Set JWT cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    user.AccessToken,
		Path:     "/",
		MaxAge:   60 * 60 * 24,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

	util.WriteJSON(w, http.StatusCreated, user)
}
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req model.RequestLoginUser
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.WriteError(w, http.StatusBadRequest, "invalid JSON payload")
		return
	}

	user, err := h.userService.Login(r.Context(), req)
	if err != nil {
		util.WriteError(w, http.StatusUnauthorized, err.Error())
		return
	}

	// TODO: prod vs dev
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    user.AccessToken,
		Path:     "/",
		MaxAge:   60 * 60 * 24,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

	util.WriteJSON(w, http.StatusOK, user)
}

func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

	util.WriteJSON(w, http.StatusOK, map[string]string{"message": "logout successful"})
}

func (h *UserHandler) UpdateUsername(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by JWT middleware)
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		util.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req struct {
		Username string `json:"username"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.WriteError(w, http.StatusBadRequest, "invalid JSON payload")
		return
	}

	// Validate username
	if len(req.Username) < 3 || len(req.Username) > 20 {
		util.WriteError(w, http.StatusBadRequest, "username must be between 3 and 20 characters")
		return
	}

	user, err := h.userService.UpdateUsername(r.Context(), userID, req.Username)
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	util.WriteJSON(w, http.StatusOK, user)
}
