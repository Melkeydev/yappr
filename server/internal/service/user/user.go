package service

import (
	"context"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	model "github.com/melkeydev/chat-go/internal/api/model"
	repo "github.com/melkeydev/chat-go/internal/repo/user"
	"github.com/melkeydev/chat-go/util"
)

type JWTClaims struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type UserService struct {
	userRepo *repo.UserRepository
	timeout  time.Duration
}

func NewUserService(userRepo *repo.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
		timeout:  time.Duration(2) * time.Second,
	}
}

func (s *UserService) CreateUser(ctx context.Context, req model.RequestCreateUser) (*repo.User, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	u := &repo.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: &hashedPassword,
	}

	user, err := s.userRepo.CreateUser(ctx, u)
	if err != nil {
		return nil, err
	}

	return user, nil

}

func (s *UserService) Login(ctx context.Context, req model.RequestLoginUser) (*model.ResponseLoginUser, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	user, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	err = util.CheckPassword(req.Password, *user.PasswordHash)
	if err != nil {
		return nil, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, JWTClaims{
		ID:       user.ID.String(),
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    user.ID.String(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	})

	secretKey := os.Getenv("secretKey")

	ss, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return nil, err
	}

	return &model.ResponseLoginUser{AccessToken: ss, Username: user.Username, ID: user.ID.String()}, nil
}

func (s *UserService) GetUserByID(ctx context.Context, id uuid.UUID) (*repo.User, error) {
	return s.userRepo.GetUserByID(ctx, id)
}

func (s *UserService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return s.userRepo.DeleteUser(ctx, id)
}
