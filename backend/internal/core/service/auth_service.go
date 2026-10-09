package service

import (
	"errors"
	"strings"
	"time"

	"strconv"

	"workshop-registration/backend/internal/adapter/storage/postgres/models"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrInvalidUserData    = errors.New("invalid user information")
	ErrEmailExists        = errors.New("email already exists")
)

type AuthService struct {
	DB        *gorm.DB
	JWTSecret []byte
}

func NewAuthService(db *gorm.DB, secret string) *AuthService {
	return &AuthService{
		DB:        db,
		JWTSecret: []byte(secret),
	}
}

func (s *AuthService) Login(email, password string) (string, error) {
	var user models.UserModel

	email = strings.ToLower(strings.TrimSpace(email))

	err := s.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", ErrInvalidCredentials
		}
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(password),
	); err != nil {
		return "", ErrInvalidCredentials
	}

	claims := jwt.RegisteredClaims{
		Subject:   "",
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(8 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	// Use database user ID as JWT subject
	claims.Subject = fmtUserID(user.ID)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(s.JWTSecret)
}

func fmtUserID(id uint) string {
	return strconv.FormatUint(uint64(id), 10)
}

func (s *AuthService) CreateUser(
	name, email, password, role string,
) (*models.UserModel, error) {

	name = strings.TrimSpace(name)
	email = strings.ToLower(strings.TrimSpace(email))

	if name == "" || email == "" || len(password) < 8 ||
		(role != "manager" && role != "staff") {
		return nil, ErrInvalidUserData
	}

	var count int64
	if err := s.DB.Model(&models.UserModel{}).
		Where("email = ?", email).
		Count(&count).Error; err != nil {
		return nil, err
	}

	if count > 0 {
		return nil, ErrEmailExists
	}

	hash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return nil, err
	}

	user := models.UserModel{
		Name:         name,
		Email:        email,
		PasswordHash: string(hash),
		Role:         role,
	}

	if err := s.DB.Create(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, ErrEmailExists
		}
		return nil, err
	}

	return &user, nil
}
