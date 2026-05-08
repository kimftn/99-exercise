package dto

type CreateUserRequest struct {
	Name string `json:"name"`
}

type UserResponse struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

type GetUsersResponse struct {
	Result bool           `json:"result"`
	Users  []UserResponse `json:"users"`
}

type GetUserByIDResponse struct {
	Result bool          `json:"result"`
	Users  *UserResponse `json:"users"`
}

type CreateUserResponse struct {
	Result bool          `json:"result"`
	User   *UserResponse `json:"user"`
}
