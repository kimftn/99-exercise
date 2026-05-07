package dto

type CreateListingRequest struct {
	Title    string  `json:"title"`
	City     string  `json:"city"`
	Price    float64 `json:"price"`
	Status   string  `json:"status"`
	Category string  `json:"category"`
}

type ListingResponse struct {
	ID       int     `json:"id"`
	Title    string  `json:"title"`
	City     string  `json:"city"`
	Price    float64 `json:"price"`
	Status   string  `json:"status"`
	Category string  `json:"category"`
}
