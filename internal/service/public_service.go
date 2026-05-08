package service

import "property-api/internal/dto"

type PublicService struct {
	listingService *ListingService
	userService    *UserService
}

func NewPublicService(listingService *ListingService, userService *UserService) *PublicService {
	return &PublicService{
		listingService: listingService,
		userService:    userService,
	}
}

func (s *PublicService) GetListings() []dto.ListingResponse {
	return s.listingService.GetListings()
}

func (s *PublicService) CreateListing(request dto.PublicCreateListingRequest) dto.ListingResponse {
	return s.listingService.CreateListing(dto.CreateListingRequest{
		Title:    request.Title,
		City:     request.City,
		Price:    request.Price,
		Status:   request.Status,
		Category: request.Category,
	})
}

func (s *PublicService) CreateUser(request dto.PublicCreateUserRequest) (dto.UserResponse, error) {
	return s.userService.CreateUser(dto.CreateUserRequest{
		Name: request.Name,
	})
}
