package service

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
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

func (s *UserService) CreateUser(ctx context.Context, req model.RequestCreateUser) (*model.ResponseLoginUser, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	log.Printf("UserService.CreateUser - Starting user creation for: %s", req.Email)

	// Validate input
	if req.Username == "" || req.Email == "" || req.Password == "" {
		log.Printf("UserService.CreateUser - Validation failed: missing required fields")
		return nil, fmt.Errorf("username, email, and password are required")
	}

	if len(req.Password) < 6 {
		log.Printf("UserService.CreateUser - Validation failed: password too short")
		return nil, fmt.Errorf("password must be at least 6 characters")
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		log.Printf("UserService.CreateUser - Password hashing failed: %v", err)
		return nil, fmt.Errorf("failed to process password")
	}

	u := &repo.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: &hashedPassword,
	}

	user, err := s.userRepo.CreateUser(ctx, u)
	if err != nil {
		log.Printf("UserService.CreateUser - Database error: %v", err)
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique") {
			return nil, fmt.Errorf("username or email already exists")
		}
		return nil, fmt.Errorf("failed to create user: %v", err)
	}

	log.Printf("UserService.CreateUser - User created successfully in database: %s", user.ID.String())

	// Generate JWT token for the new user
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

	return &model.ResponseLoginUser{
		AccessToken: ss,
		Username:    user.Username,
		ID:          user.ID.String(),
	}, nil
}

func (s *UserService) Login(ctx context.Context, req model.RequestLoginUser) (*model.ResponseLoginUser, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	log.Printf("UserService.Login - Starting login attempt for email: %s", req.Email)

	user, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		log.Printf("UserService.Login - Database error: %v", err)
		return nil, fmt.Errorf("failed to authenticate user")
	}
	
	if user == nil {
		log.Printf("UserService.Login - User not found for email: %s", req.Email)
		return nil, fmt.Errorf("invalid email or password")
	}
	
	if user.PasswordHash == nil {
		log.Printf("UserService.Login - User has no password hash: %s", req.Email)
		return nil, fmt.Errorf("invalid user account")
	}

	err = util.CheckPassword(req.Password, *user.PasswordHash)
	if err != nil {
		log.Printf("UserService.Login - Password check failed for user: %s", user.ID.String())
		return nil, fmt.Errorf("invalid email or password")
	}

	log.Printf("UserService.Login - Password verified successfully for user: %s", user.ID.String())

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
		log.Printf("UserService.Login - JWT signing failed for user: %s, error: %v", user.ID.String(), err)
		return nil, fmt.Errorf("failed to generate authentication token")
	}

	log.Printf("UserService.Login - Login successful for user: %s (%s)", user.ID.String(), user.Username)
	return &model.ResponseLoginUser{AccessToken: ss, Username: user.Username, ID: user.ID.String()}, nil
}

func (s *UserService) GetUserByID(ctx context.Context, id uuid.UUID) (*repo.User, error) {
	return s.userRepo.GetUserByID(ctx, id)
}

func (s *UserService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return s.userRepo.DeleteUser(ctx, id)
}

func (s *UserService) UpdateUsername(ctx context.Context, userID string, newUsername string) (*model.ResponseLoginUser, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	// Parse user ID
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	// Update username in database
	user, err := s.userRepo.UpdateUsername(ctx, uid, newUsername)
	if err != nil {
		return nil, err
	}

	return &model.ResponseLoginUser{
		ID:       user.ID.String(),
		Username: user.Username,
	}, nil
}
