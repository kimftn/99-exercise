package app

import (
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"

	"property-api/internal/domain"
	"property-api/internal/repository"
)

type stubUserRepository struct {
	users []domain.User
}

func newStubUserRepository() repository.UserRepository {
	return &stubUserRepository{
		users: []domain.User{
			{ID: 1, Name: "Ava Hartono", CreatedAt: 1715000000, UpdatedAt: 1715000000},
			{ID: 2, Name: "Rizky Pratama", CreatedAt: 1715003600, UpdatedAt: 1715003600},
		},
	}
}

func (r *stubUserRepository) Create(user domain.User) (domain.User, error) {
	user.ID = len(r.users) + 1
	r.users = append(r.users, user)
	return user, nil
}

func (r *stubUserRepository) GetAll(pageNum int, pageSize int) ([]domain.User, error) {
	offset := (pageNum - 1) * pageSize
	if offset >= len(r.users) {
		return []domain.User{}, nil
	}

	end := offset + pageSize
	if end > len(r.users) {
		end = len(r.users)
	}

	return r.users[offset:end], nil
}

func (r *stubUserRepository) GetByID(id int) (domain.User, bool, error) {
	for _, user := range r.users {
		if user.ID == id {
			return user, true, nil
		}
	}

	return domain.User{}, false, nil
}

func TestHealthEndpoint(t *testing.T) {
	server := newServerWithDependencies(newStubUserRepository())

	req := httptest.NewRequest("GET", "/health", nil)
	resp, err := server.Test(req)
	if err != nil {
		t.Fatalf("health request failed: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestGetListingsEndpoint(t *testing.T) {
	server := newServerWithDependencies(newStubUserRepository())

	req := httptest.NewRequest("GET", "/listings", nil)
	resp, err := server.Test(req)
	if err != nil {
		t.Fatalf("listings request failed: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestCreatePublicUserEndpoint(t *testing.T) {
	server := newServerWithDependencies(newStubUserRepository())

	body := `{"name":"Dina Saputra"}`
	req := httptest.NewRequest("POST", "/public-api/users", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := server.Test(req)
	if err != nil {
		t.Fatalf("public user request failed: %v", err)
	}

	if resp.StatusCode != 201 {
		t.Fatalf("expected status 201, got %d", resp.StatusCode)
	}

	var payload struct {
		Result bool `json:"result"`
		User   struct {
			ID        int    `json:"id"`
			Name      string `json:"name"`
			CreatedAt int64  `json:"created_at"`
			UpdatedAt int64  `json:"updated_at"`
		} `json:"user"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if !payload.Result {
		t.Fatalf("expected result true")
	}

	if payload.User.Name != "Dina Saputra" {
		t.Fatalf("expected created user name Dina Saputra, got %s", payload.User.Name)
	}
}

func TestCreateUserEndpointResponseShape(t *testing.T) {
	server := newServerWithDependencies(newStubUserRepository())

	body := `{"name":"Nadia Putri"}`
	req := httptest.NewRequest("POST", "/users", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := server.Test(req)
	if err != nil {
		t.Fatalf("create user request failed: %v", err)
	}

	if resp.StatusCode != 201 {
		t.Fatalf("expected status 201, got %d", resp.StatusCode)
	}

	var payload struct {
		Result bool `json:"result"`
		User   struct {
			ID        int    `json:"id"`
			Name      string `json:"name"`
			CreatedAt int64  `json:"created_at"`
			UpdatedAt int64  `json:"updated_at"`
		} `json:"user"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		t.Fatalf("decode create user response: %v", err)
	}

	if !payload.Result {
		t.Fatalf("expected result true")
	}

	if payload.User.Name != "Nadia Putri" {
		t.Fatalf("expected created user name Nadia Putri, got %s", payload.User.Name)
	}
}

func TestGetUsersEndpointResponseShape(t *testing.T) {
	server := newServerWithDependencies(newStubUserRepository())

	req := httptest.NewRequest("GET", "/users", nil)
	resp, err := server.Test(req)
	if err != nil {
		t.Fatalf("users request failed: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	var payload struct {
		Result bool `json:"result"`
		Users  []struct {
			ID        int    `json:"id"`
			Name      string `json:"name"`
			CreatedAt int64  `json:"created_at"`
			UpdatedAt int64  `json:"updated_at"`
		} `json:"users"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		t.Fatalf("decode users response: %v", err)
	}

	if !payload.Result {
		t.Fatalf("expected result true")
	}

	if len(payload.Users) == 0 {
		t.Fatalf("expected seeded users")
	}
}

func TestGetUsersEndpointPagination(t *testing.T) {
	server := newServerWithDependencies(newStubUserRepository())

	req := httptest.NewRequest("GET", "/users?page_num=2&page_size=1", nil)
	resp, err := server.Test(req)
	if err != nil {
		t.Fatalf("paginated users request failed: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	var payload struct {
		Result bool `json:"result"`
		Users  []struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"users"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		t.Fatalf("decode paginated users response: %v", err)
	}

	if len(payload.Users) != 1 {
		t.Fatalf("expected 1 user, got %d", len(payload.Users))
	}

	if payload.Users[0].ID != 2 {
		t.Fatalf("expected user id 2, got %d", payload.Users[0].ID)
	}
}

func TestGetUserByIDEndpointResponseShape(t *testing.T) {
	server := newServerWithDependencies(newStubUserRepository())

	req := httptest.NewRequest("GET", "/users/1", nil)
	resp, err := server.Test(req)
	if err != nil {
		t.Fatalf("user by id request failed: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	var payload struct {
		Result bool `json:"result"`
		Users  struct {
			ID        int    `json:"id"`
			Name      string `json:"name"`
			CreatedAt int64  `json:"created_at"`
			UpdatedAt int64  `json:"updated_at"`
		} `json:"users"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		t.Fatalf("decode user by id response: %v", err)
	}

	if !payload.Result {
		t.Fatalf("expected result true")
	}

	if payload.Users.ID != 1 {
		t.Fatalf("expected user id 1, got %d", payload.Users.ID)
	}
}

func TestGetUserByIDEndpointNotFoundReturns200(t *testing.T) {
	server := newServerWithDependencies(newStubUserRepository())

	req := httptest.NewRequest("GET", "/users/999", nil)
	resp, err := server.Test(req)
	if err != nil {
		t.Fatalf("user by id not found request failed: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	var payload struct {
		Result bool `json:"result"`
		Users  any  `json:"users"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		t.Fatalf("decode user by id not found response: %v", err)
	}

	if payload.Result {
		t.Fatalf("expected result false")
	}

	if payload.Users != nil {
		t.Fatalf("expected users to be nil, got %#v", payload.Users)
	}
}
