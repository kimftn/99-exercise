package repository

import (
	"property-api/internal/domain"
	"property-api/internal/store"
)

type ListingRepository interface {
	Create(listing domain.Listing) domain.Listing
	GetAll() []domain.Listing
}

type InMemoryListingRepository struct {
	store *store.Store
}

func NewInMemoryListingRepository(store *store.Store) *InMemoryListingRepository {
	return &InMemoryListingRepository{store: store}
}

func (r *InMemoryListingRepository) Create(listing domain.Listing) domain.Listing {
	listing.ID = len(r.store.Listings) + 1
	r.store.Listings = append(r.store.Listings, listing)
	return listing
}

func (r *InMemoryListingRepository) GetAll() []domain.Listing {
	return r.store.Listings
}
