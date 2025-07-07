package handler

import (
	"encoding/json"
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
		util.WriteError(w, http.StatusBadRequest, "invalid JSON payload")
		return
	}

	res, err := h.userService.CreateUser(r.Context(), req)
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	util.WriteJSON(w, http.StatusCreated, res)
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
		Domain:   "localhost",
		MaxAge:   60 * 60 * 24,
		HttpOnly: true,
		Secure:   false,
	})

	util.WriteJSON(w, http.StatusOK, user)
}

func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    "",
		Path:     "/",
		Domain:   "localhost",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   false,
	})

	util.WriteJSON(w, http.StatusOK, map[string]string{"message": "logout successful"})
}
