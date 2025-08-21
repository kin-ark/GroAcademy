package services

import (
	"math"
	"strconv"

	"github.com/kin-ark/GroAcademy/internal/models"
	"github.com/kin-ark/GroAcademy/internal/repositories"
)

type UserService interface {
	GetUsers(query models.SearchQuery) ([]models.User, models.PaginationResponse, error)
	BuildUsersResponse(users []models.User) []models.UsersResponse
}

type userService struct {
	userRepo repositories.UserRepository
}

func NewUserService(r repositories.UserRepository) UserService {
	return &userService{userRepo: r}
}

func (s *userService) GetUsers(query models.SearchQuery) ([]models.User, models.PaginationResponse, error) {
	query.Normalize()

	users, totalItems, err := s.userRepo.GetAllUsers(query)
	if err != nil {
		return nil, models.PaginationResponse{}, err
	}

	totalPages := int(math.Ceil(float64(totalItems) / float64(query.Limit)))
	if query.Page > totalPages && totalPages > 0 {
		query.Page = totalPages
	}

	pagination := models.PaginationResponse{
		CurrentPage: query.Page,
		TotalPages:  totalPages,
		TotalItems:  int(totalItems),
	}

	return users, pagination, nil
}

func (s *userService) BuildUsersResponse(users []models.User) []models.UsersResponse {
	var res []models.UsersResponse
	for _, user := range users {
		stringId := strconv.FormatUint(uint64(user.ID), 10)
		res = append(res, models.UsersResponse{
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			Username:  user.Username,
			Balance:   user.Balance,
			ID:        stringId,
		})
	}
	return res
}
