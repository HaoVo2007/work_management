package user

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"time"
	"work-management/internal/domain/user/dto/request"
	"work-management/internal/domain/user/dto/response"
	"work-management/internal/domain/user/model"
	"work-management/internal/pkg/aws"

	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	RegisterUser(ctx context.Context, req request.CreateUserRequest) (*response.CreateUserResponse, error)
	LoginUser(ctx context.Context, req request.LoginUserRequest) (string, error)
	LogoutUser(ctx context.Context, userID string) error
	UploadAvatar(ctx context.Context, userID string, avatar *multipart.FileHeader) (string, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{
		repository: repository,
	}
}

func (s *service) RegisterUser(ctx context.Context, req request.CreateUserRequest) (*response.CreateUserResponse, error) {

	if req.Name == "" {
		return nil, fmt.Errorf("name is required")
	}

	if req.Email == "" {
		return nil, fmt.Errorf("email is required")
	}

	if req.Password == "" {
		return nil, fmt.Errorf("password is required")
	}

	user, err := s.repository.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	if user != nil {
		return nil, fmt.Errorf("user already exists")
	}

	hashedPassword := s.hashPassword(req.Password)
	newUserID := primitive.NewObjectID()
	token, refreshToken := s.generateToken(newUserID.Hex(), req.Name, "user")
	inviteLink := fmt.Sprintf("INVITE_LINK/%s", newUserID.Hex())
	userData := &model.User{
		ID:           newUserID,
		Name:         req.Name,
		Email:        req.Email,
		Password:     hashedPassword,
		Role:         "user",
		InviteLink:   inviteLink,
		Token:        token,
		RefreshToken: refreshToken,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err = s.repository.Create(ctx, userData)
	if err != nil {
		return nil, err
	}

	return &response.CreateUserResponse{
		ID:           userData.ID.Hex(),
		Name:         userData.Name,
		Email:        userData.Email,
		Role:         userData.Role,
		InviteLink:   userData.InviteLink,
		Token:        userData.Token,
		RefreshToken: userData.RefreshToken,
		CreatedAt:    userData.CreatedAt,
		UpdatedAt:    userData.UpdatedAt,
	}, nil
}

func (s *service) LoginUser(ctx context.Context, req request.LoginUserRequest) (string, error) {

	if req.Email == "" && req.Password == "" {
		return "", fmt.Errorf("email and password are required")
	}

	user, err := s.repository.FindByEmail(ctx, req.Email)
	if err != nil {
		return "", err
	}

	if user == nil {
		return "", fmt.Errorf("invalid email or password")
	}

	isValid, _ := s.verifyPassword(user.Password, req.Password)
	if !isValid {
		return "", fmt.Errorf("invalid email or password")
	}

	updateFileds := bson.M{
		"token":         user.Token,
		"refresh_token": user.RefreshToken,
		"updated_at":    time.Now(),
	}

	err = s.repository.UpdateByID(ctx, user.ID, updateFileds)
	if err != nil {
		return "", fmt.Errorf("failed to update user tokens: %w", err)
	}

	return user.Token, nil

}

func (s *service) LogoutUser(ctx context.Context, userID string) error {

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	updateFields := bson.M{
		"token":         "",
		"refresh_token": "",
		"updated_at":    time.Now().Format(time.RFC3339),
	}

	err = s.repository.UpdateByID(ctx, objectID, updateFields)
	if err != nil {
		return fmt.Errorf("failed to update user tokens: %w", err)
	}

	return nil

}

func (s *service) UploadAvatar(ctx context.Context, userID string, avatar *multipart.FileHeader) (string, error) {
	
	if userID == "" {
		return "", fmt.Errorf("user_id is required")
	}

	if avatar == nil {
		return "", fmt.Errorf("avatar is required")
	}

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return "", err
	}

	user, err := s.repository.FindByID(ctx, objectID)
	if err != nil {
		return "", err
	}

	if user == nil {
		return "", fmt.Errorf("user not found")
	}

	if user.Avatar != nil {
		err = aws.DeleteFile(ctx, *user.Avatar)
		if err != nil {
			return "", err
		}
	}

	key, err := aws.UploadPrivateFile(ctx, avatar, "avatars")
	if err != nil {
		return "", err
	}
	
	uploadFields := bson.M{
		"avatar":        key,
		"updated_at":    time.Now(),
	}

	err = s.repository.UpdateByID(ctx, objectID, uploadFields)
	if err != nil {
		return "", fmt.Errorf("failed to update user avatar: %w", err)
	}

	avatarUrl, err := aws.GetPresignedURL(ctx, key, time.Hour)
	if err != nil {
		return "", err
	}

	return *avatarUrl, nil

}

func (s *service) hashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

func (s *service) generateToken(userID, userName, role string) (string, string) {

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Panic("JWT_SECRET not set")
	}

	claims := jwt.MapClaims{
		"user_id":   userID,
		"user_name": userName,
		"role":      role,
		"exp":       jwt.NewNumericDate(time.Now().Add(time.Hour * 8)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		log.Panic(err)
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
	})
	refreshTokenString, err := refreshToken.SignedString([]byte(secret))
	if err != nil {
		log.Panic(err)
	}

	return tokenString, refreshTokenString

}

func (s *service) verifyPassword(userPassword string, providedPassword string) (bool, string) {

	err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(providedPassword))
	check := true
	msg := ""

	if err != nil {
		msg = "Login or password is incorrect"
		check = false
	}

	return check, msg

}
