package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"

	"github.com/liza/labwork_45/internal/config"
	"github.com/liza/labwork_45/internal/model"

	"golang.org/x/crypto/bcrypt"
)

// These constansts represent token's duration
const (
	accessTokenTTL  = 24 * time.Hour
	refreshTokenTTL = 72 * time.Hour
)

// tokenClaims struct contains information about the claims associated with the given token
type tokenClaims struct {
	Role string `json:"role"`
	jwt.StandardClaims
}

type AuthApiService struct {
	rps AuthApiRepository
}

func NewAuthApiService(rps AuthApiRepository) *AuthApiService {
	return &AuthApiService{rps: rps}
}

type AuthApiRepository interface {
	GetAll(ctx context.Context) ([]*model.User, error)
	InsertUser(ctx context.Context, user *model.SaveUser) (uuid.UUID, error)
	GetUserByLogin(ctx context.Context, login string) (*model.HashedLogin, error)
	SaveRefreshToken(ctx context.Context, ID uuid.UUID, refreshToken []byte) error
	GetUserByID(ctx context.Context, ID uuid.UUID) (*model.User, error)
	DeleteUserByID(ctx context.Context, ID uuid.UUID) error
}

func (srv *AuthApiService) GetAll(ctx context.Context) ([]*model.User, error) {
	return srv.rps.GetAll(ctx)
}

func (srv *AuthApiService) LoginUser(ctx context.Context, auth *model.Login) (string, string, error) {
	cfg, err := config.NewConfig()
	if err != nil {
		return "", "", fmt.Errorf("NewConfig: %w", err)
	}
	selectedUser, err := srv.rps.GetUserByLogin(ctx, auth.Login)
	if err != nil {
		return "", "", fmt.Errorf("GetUserByLogin: %w", err)
	}
	err = bcrypt.CompareHashAndPassword(selectedUser.Password, []byte(auth.Password))
	if err != nil {
		return "", "", fmt.Errorf("CompareHashAndPassword: %w", err)
	}
	// GenerateAccessToken
	accessToken, refreshToken, err := GenerateAccessAndRefreshTokens(cfg.SigningKey, selectedUser.Role, selectedUser.ID)
	if err != nil {
		return "", "", fmt.Errorf("GenerateAccessAndRefreshTokens: %w", err)
	}
	// HashRefreshToken
	hashedRefreshToken, err := HashRefreshToken(refreshToken)
	if err != nil {
		return "", "", fmt.Errorf("HashRefreshToken: %w", err)
	}
	// SaveRefreshToken
	err = srv.rps.SaveRefreshToken(ctx, selectedUser.ID, hashedRefreshToken)
	if err != nil {
		return "", "", fmt.Errorf("SaveRefreshToken: %w", err)
	}
	// CompareTokenIDs
	compID, err := CompareTokenIDs(accessToken, refreshToken, cfg.SigningKey)
	if err != nil {
		return "", "", fmt.Errorf("CompareTokenIDs: %w", err)
	}
	if !compID {
		return "", "", fmt.Errorf("invalid token(campare error): %w", err)
	}
	return accessToken, refreshToken, nil
}

func (srv *AuthApiService) SignUpUser(ctx context.Context, user *model.SignUp) (uuid.UUID, error) {

	hashedPassword := hashPassword([]byte(user.Password))
	saveUser := &model.SaveUser{
		Login:    user.Login,
		Password: hashedPassword,
		Username: user.Username,
		Role:     user.Role,
	}
	id, err := srv.rps.InsertUser(ctx, saveUser)
	if err != nil {
		return uuid.Nil, fmt.Errorf("InsertUser: %w", err)
	}
	return id, nil
}

func (srv *AuthApiService) GetPersonalInfo(ctx context.Context, ID uuid.UUID) (*model.User, error) {
	userInfo, err := srv.rps.GetUserByID(ctx, ID)
	if err != nil {
		return nil, fmt.Errorf("GetUserByID: %w", err)
	}
	return userInfo, nil
}

func (srv *AuthApiService) DeleteUserByID(ctx context.Context, ID uuid.UUID) error {
	err := srv.rps.DeleteUserByID(ctx, ID)
	if err != nil {
		return fmt.Errorf("DeleteUserByID: %w", err)
	}
	return nil
}

// HashPassword func returns hashed password using bcrypt algorithm
func hashPassword(password []byte) []byte {
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return nil
	}
	return hashedPassword
}

// HashRefreshToken func returns hashed refresh token using bcrypt algorithm
func HashRefreshToken(refreshToken string) ([]byte, error) {
	hash := sha256.New()
	hash.Write([]byte(refreshToken))
	hashBytes := hash.Sum(nil)
	hashString := hex.EncodeToString(hashBytes)

	return []byte(hashString), nil
}

// GenerateAccessAndRefreshTokens func returns access & refresh tokens
func GenerateAccessAndRefreshTokens(key, role string, id uuid.UUID) (access, refresh string, err error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		Role: role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(accessTokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
			Id:        id.String(),
		},
	})

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(refreshTokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
			Id:        id.String(),
		},
	})
	access, err = accessToken.SignedString([]byte(key))
	if err != nil {
		return "", "", fmt.Errorf("SignedString(access): %w", err)
	}
	refresh, err = refreshToken.SignedString([]byte(key))
	if err != nil {
		return "", "", fmt.Errorf("SignedString(refresh): %w", err)
	}
	return access, refresh, err
}

// CompareTokenIDs func compares token ids
func CompareTokenIDs(accessToken, refreshToken, key string) (bool, error) {
	accessID, err := ExtractIDFromToken(accessToken, key)
	if err != nil {
		return false, fmt.Errorf("ExtractIDFromToken: %w", err)
	}

	refreshID, err := ExtractIDFromToken(refreshToken, key)
	if err != nil {
		return false, fmt.Errorf("ExtractIDFromToken: %w", err)
	}
	return accessID == refreshID, nil
}

// ExtractIDFromToken extracts the identifier (ID) from the payload (claims) of the token.
func ExtractIDFromToken(tokenString, key string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})
	if err != nil {
		return "", fmt.Errorf("Parse(): %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if id, ok := claims["jti"].(string); ok {
			return id, nil
		}
	}

	return "", fmt.Errorf("error extracting ID from token: %v", token)
}
