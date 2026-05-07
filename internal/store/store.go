package store

import "property-api/internal/domain"

type Store struct {
	Listings []domain.Listing
	Users    []domain.User
}

func New() *Store {
	return &Store{
		Listings: []domain.Listing{
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
		Users: []domain.User{
			{
				ID:    1,
				Name:  "Ava Hartono",
				Email: "ava@example.com",
				Role:  "agent",
			},
			{
				ID:    2,
				Name:  "Rizky Pratama",
				Email: "rizky@example.com",
				Role:  "buyer",
			},
		},
	}
}
