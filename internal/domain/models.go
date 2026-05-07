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
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}
