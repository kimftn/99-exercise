package repository

import "property-api/internal/domain"

type ListingRepository interface {
	Create(listing domain.Listing) domain.Listing
	GetAll() []domain.Listing
}

type InMemoryListingRepository struct {
	listings []domain.Listing
}

func NewInMemoryListingRepository() *InMemoryListingRepository {
	return &InMemoryListingRepository{
		listings: []domain.Listing{
			{
				ID:       1,
				Title:    "Sunny Apartment",
				City:     "Jakarta",
				Price:    1250000000,
				Status:   "available",
				Category: "buy",
			},
			{
				ID:       2,
				Title:    "Cozy Studio",
				City:     "Bandung",
				Price:    4500000,
				Status:   "available",
				Category: "rent",
			},
		},
	}
}

func (r *InMemoryListingRepository) Create(listing domain.Listing) domain.Listing {
	listing.ID = len(r.listings) + 1
	r.listings = append(r.listings, listing)
	return listing
}

func (r *InMemoryListingRepository) GetAll() []domain.Listing {
	return r.listings
}
