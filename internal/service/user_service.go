package service

import (
	"errors"
	"fmt"
	"time"

	"property-api/internal/domain"
	"property-api/internal/dto"
	"property-api/internal/repository"
)

var ErrUserNotFound = errors.New("user not found")

type UserService struct {
	repository repository.UserRepository
}

const (
	defaultPageNum  = 1
	defaultPageSize = 10
)

func NewUserService(repository repository.UserRepository) *UserService {
	return &UserService{repository: repository}
}

func (s *UserService) CreateUser(request dto.CreateUserRequest) (dto.UserResponse, error) {
	now := time.Now().UnixMicro()
	user := domain.User{
		Name:      request.Name,
		CreatedAt: now,
		UpdatedAt: now,
	}

	created, err := s.repository.Create(user)
	if err != nil {
		return dto.UserResponse{}, fmt.Errorf("create user: %w", err)
	}

	return toUserResponse(created), nil
}

func (s *UserService) GetUsers(pageNum int, pageSize int) ([]dto.UserResponse, error) {
	if pageNum <= 0 {
		pageNum = defaultPageNum
	}

	if pageSize <= 0 {
		pageSize = defaultPageSize
	}

	users, err := s.repository.GetAll(pageNum, pageSize)
	if err != nil {
		return nil, fmt.Errorf("get users: %w", err)
	}

	response := make([]dto.UserResponse, 0, len(users))

	for _, user := range users {
		response = append(response, toUserResponse(user))
	}

	return response, nil
}

func (s *UserService) GetUserByID(id int) (dto.UserResponse, error) {
	user, found, err := s.repository.GetByID(id)
	if err != nil {
		return dto.UserResponse{}, fmt.Errorf("get user by id: %w", err)
	}

	if !found {
		return dto.UserResponse{}, ErrUserNotFound
	}

	return toUserResponse(user), nil
}

func toUserResponse(user domain.User) dto.UserResponse {
	return dto.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
