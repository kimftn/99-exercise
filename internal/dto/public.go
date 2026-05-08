package dto

type PublicCreateListingRequest struct {
	Title    string  `json:"title"`
	City     string  `json:"city"`
	Price    float64 `json:"price"`
	Status   string  `json:"status"`
	Category string  `json:"category"`
}

type PublicCreateUserRequest struct {
	Name string `json:"name"`
}
