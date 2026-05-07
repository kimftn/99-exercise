package service

import (
	"errors"

	"property-api/internal/domain"
	"property-api/internal/dto"
	"property-api/internal/repository"
)

var ErrUserNotFound = errors.New("user not found")

type UserService struct {
	repository repository.UserRepository
}

func NewUserService(repository repository.UserRepository) *UserService {
	return &UserService{repository: repository}
}

func (s *UserService) CreateUser(request dto.CreateUserRequest) dto.UserResponse {
	user := domain.User{
		Name:  request.Name,
		Email: request.Email,
		Role:  request.Role,
	}

	created := s.repository.Create(user)
	return toUserResponse(created)
}

func (s *UserService) GetUsers() []dto.UserResponse {
	users := s.repository.GetAll()
	response := make([]dto.UserResponse, 0, len(users))

	for _, user := range users {
		response = append(response, toUserResponse(user))
	}

	return response
}

func (s *UserService) GetUserByID(id int) (dto.UserResponse, error) {
	user, found := s.repository.GetByID(id)
	if !found {
		return dto.UserResponse{}, ErrUserNotFound
	}

	return toUserResponse(user), nil
}

func toUserResponse(user domain.User) dto.UserResponse {
	return dto.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}
}
