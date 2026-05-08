package service

import (
	"errors"

	"property-api/internal/client"
	"property-api/internal/dto"
)

var ErrPublicListingRejected = errors.New("public listing rejected")

type PublicService struct {
	listingClient client.ListingClient
	userService   *UserService
}

func NewPublicService(listingClient client.ListingClient, userService *UserService) *PublicService {
	return &PublicService{
		listingClient: listingClient,
		userService:   userService,
	}
}

func (s *PublicService) GetListings(pageNum int, pageSize int, userID *int) ([]dto.ListingResponse, error) {
	return s.listingClient.GetListings(pageNum, pageSize, userID)
}

func (s *PublicService) CreateListing(request dto.PublicCreateListingRequest) (dto.ListingResponse, error) {
	listing, err := s.listingClient.CreateListing(request)
	if err != nil {
		if err == client.ErrListingRequestRejected {
			return dto.ListingResponse{}, ErrPublicListingRejected
		}
		return dto.ListingResponse{}, err
	}

	return listing, nil
}

func (s *PublicService) CreateUser(request dto.PublicCreateUserRequest) (dto.UserResponse, error) {
	return s.userService.CreateUser(dto.CreateUserRequest{
		Name: request.Name,
	})
}
