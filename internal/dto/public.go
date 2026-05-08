package dto

type PublicCreateListingRequest struct {
	UserID      int    `json:"user_id" form:"user_id"`
	ListingType string `json:"listing_type" form:"listing_type"`
	Price       int64  `json:"price" form:"price"`
}

type PublicCreateUserRequest struct {
	Name string `json:"name"`
}
