package service

import (
	"property-api/internal/domain"
	"property-api/internal/dto"
	"property-api/internal/repository"
)

type ListingService struct {
	repository repository.ListingRepository
}

func NewListingService(repository repository.ListingRepository) *ListingService {
	return &ListingService{repository: repository}
}

func (s *ListingService) CreateListing(request dto.CreateListingRequest) dto.ListingResponse {
	listing := domain.Listing{
		Title:    request.Title,
		City:     request.City,
		Price:    request.Price,
		Status:   request.Status,
		Category: request.Category,
	}

	if listing.Status == "" {
		listing.Status = "available"
	}

	created := s.repository.Create(listing)
	return toListingResponse(created)
}

func (s *ListingService) GetListings() []dto.ListingResponse {
	listings := s.repository.GetAll()
	response := make([]dto.ListingResponse, 0, len(listings))

	for _, listing := range listings {
		response = append(response, toListingResponse(listing))
	}

	return response
}

func toListingResponse(listing domain.Listing) dto.ListingResponse {
	return dto.ListingResponse{
		ID:       listing.ID,
		Title:    listing.Title,
		City:     listing.City,
		Price:    listing.Price,
		Status:   listing.Status,
		Category: listing.Category,
	}
}
