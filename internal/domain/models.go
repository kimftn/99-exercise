package domain

type Listing struct {
	ID          int    `json:"id"`
	UserID      int    `json:"user_id"`
	Price       int64  `json:"price"`
	ListingType string `json:"listing_type"`
	CreatedAt   int64  `json:"created_at"`
	UpdatedAt   int64  `json:"updated_at"`
}

type User struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}
