package services

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kin-ark/GroAcademy/internal/models"
	"github.com/kin-ark/GroAcademy/internal/repositories"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService interface {
	RegisterUser(input models.RegisterInput) (*models.User, error)
	LoginUser(input models.LoginInput) (string, string, error)
}

type authService struct {
	userRepo repositories.UserRepository
}

func NewAuthService(r repositories.UserRepository) AuthService {
	return &authService{userRepo: r}
}

var ErrInvalidCredentials = errors.New("invalid identifier or password")

func (s *authService) RegisterUser(body models.RegisterInput) (*models.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	if err != nil {
		return nil, err
	}

	user := models.User{FirstName: body.FirstName, LastName: body.LastName, Username: body.Username, Email: body.Email, Password: string(hash), Role: "user", Balance: 0}
	if err := s.userRepo.Create(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *authService) LoginUser(input models.LoginInput) (string, string, error) {
	user, err := s.userRepo.FindByIdentifier(input.Identifier)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", "", ErrInvalidCredentials
		}
		return "", "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		return "", "", ErrInvalidCredentials
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  user.Username,
		"role": user.Role,
		"exp":  time.Now().Add(time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		return "", "", err
	}

	return tokenString, user.Username, nil
}
