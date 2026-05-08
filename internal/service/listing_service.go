package service

import (
	"errors"
	"fmt"
	"time"

	"property-api/internal/domain"
	"property-api/internal/dto"
	"property-api/internal/repository"
)

var ErrInvalidListingType = errors.New("invalid listing type")
var ErrListingUserNotFound = errors.New("listing user not found")

type ListingService struct {
	repository     repository.ListingRepository
	userRepository repository.UserRepository
}

const (
	defaultListingPageNum  = 1
	defaultListingPageSize = 10
)

func NewListingService(
	repository repository.ListingRepository,
	userRepository repository.UserRepository,
) *ListingService {
	return &ListingService{
		repository:     repository,
		userRepository: userRepository,
	}
}

func (s *ListingService) CreateListing(request dto.CreateListingRequest) (dto.ListingResponse, error) {
	if request.ListingType != "rent" && request.ListingType != "sale" {
		return dto.ListingResponse{}, ErrInvalidListingType
	}

	_, found, err := s.userRepository.GetByID(request.UserID)
	if err != nil {
		return dto.ListingResponse{}, fmt.Errorf("get user for listing: %w", err)
	}

	if !found {
		return dto.ListingResponse{}, ErrListingUserNotFound
	}

	now := time.Now().UnixMicro()
	listing := domain.Listing{
		UserID:      request.UserID,
		Price:       request.Price,
		ListingType: request.ListingType,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	created, err := s.repository.Create(listing)
	if err != nil {
		return dto.ListingResponse{}, fmt.Errorf("create listing: %w", err)
	}

	return toListingResponse(created), nil
}

func (s *ListingService) GetListings(pageNum int, pageSize int, userID *int) ([]dto.ListingResponse, error) {
	if pageNum <= 0 {
		pageNum = defaultListingPageNum
	}

	if pageSize <= 0 {
		pageSize = defaultListingPageSize
	}

	listings, err := s.repository.GetAll(pageNum, pageSize, userID)
	if err != nil {
		return nil, fmt.Errorf("get listings: %w", err)
	}

	response := make([]dto.ListingResponse, 0, len(listings))

	for _, listing := range listings {
		response = append(response, toListingResponse(listing))
	}

	return response, nil
}

func toListingResponse(listing domain.Listing) dto.ListingResponse {
	return dto.ListingResponse{
		ID:          listing.ID,
		UserID:      listing.UserID,
		Price:       listing.Price,
		ListingType: listing.ListingType,
		CreatedAt:   listing.CreatedAt,
		UpdatedAt:   listing.UpdatedAt,
	}
}
