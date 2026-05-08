package dto

type CreateListingRequest struct {
	UserID      int    `json:"user_id" form:"user_id"`
	ListingType string `json:"listing_type" form:"listing_type"`
	Price       int64  `json:"price" form:"price"`
}

type ListingResponse struct {
	ID          int    `json:"id"`
	UserID      int    `json:"user_id"`
	Price       int64  `json:"price"`
	ListingType string `json:"listing_type"`
	CreatedAt   int64  `json:"created_at"`
	UpdatedAt   int64  `json:"updated_at"`
}

type GetListingsResponse struct {
	Result   bool              `json:"result"`
	Listings []ListingResponse `json:"listings"`
}

type CreateListingResponse struct {
	Result  bool             `json:"result"`
	Listing *ListingResponse `json:"listing"`
}

type ListingErrorResponse struct {
	Return   bool              `json:"return"`
	Listings []ListingResponse `json:"listings"`
}
