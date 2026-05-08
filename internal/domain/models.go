package domain

type Listing struct {
	ID       int     `json:"id"`
	Title    string  `json:"title"`
	City     string  `json:"city"`
	Price    float64 `json:"price"`
	Status   string  `json:"status"`
	Category string  `json:"category"`
}

type User struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}
