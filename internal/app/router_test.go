package app

import (
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"

	"property-api/internal/client"
	"property-api/internal/domain"
	"property-api/internal/dto"
	"property-api/internal/repository"
)

type stubUserRepository struct {
	users []domain.User
}

type stubListingRepository struct {
	listings []domain.Listing
}

type stubPublicListingClient struct {
	listings []domain.Listing
}

func newStubUserRepository() repository.UserRepository {
	return &stubUserRepository{
		users: []domain.User{
			{ID: 1, Name: "Ava Hartono", CreatedAt: 1715000000, UpdatedAt: 1715000000},
			{ID: 2, Name: "Rizky Pratama", CreatedAt: 1715003600, UpdatedAt: 1715003600},
		},
	}
}

func newStubListingRepository() repository.ListingRepository {
	return &stubListingRepository{
		listings: []domain.Listing{
			{ID: 1, UserID: 1, Price: 1250000000, ListingType: "sale", CreatedAt: 1715000000, UpdatedAt: 1715000000},
			{ID: 2, UserID: 2, Price: 4500000, ListingType: "rent", CreatedAt: 1715003600, UpdatedAt: 1715003600},
		},
	}
}

func newStubPublicListingClient() client.ListingClient {
	return &stubPublicListingClient{
		listings: []domain.Listing{
			{ID: 1, UserID: 1, Price: 1250000000, ListingType: "sale", CreatedAt: 1715000000, UpdatedAt: 1715000000},
			{ID: 2, UserID: 2, Price: 4500000, ListingType: "rent", CreatedAt: 1715003600, UpdatedAt: 1715003600},
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

func (r *stubListingRepository) Create(listing domain.Listing) (domain.Listing, error) {
	listing.ID = len(r.listings) + 1
	r.listings = append(r.listings, listing)
	return listing, nil
}

func (r *stubListingRepository) GetAll(pageNum int, pageSize int, userID *int) ([]domain.Listing, error) {
	filtered := make([]domain.Listing, 0)
	for _, listing := range r.listings {
		if userID != nil && listing.UserID != *userID {
			continue
		}
		filtered = append(filtered, listing)
	}

	offset := (pageNum - 1) * pageSize
	if offset >= len(filtered) {
		return []domain.Listing{}, nil
	}

	end := offset + pageSize
	if end > len(filtered) {
		end = len(filtered)
	}

	return filtered[offset:end], nil
}

func (c *stubPublicListingClient) GetListings(pageNum int, pageSize int, userID *int) ([]dto.ListingResponse, error) {
	filtered := make([]dto.ListingResponse, 0)
	for _, listing := range c.listings {
		if userID != nil && listing.UserID != *userID {
			continue
		}
		filtered = append(filtered, dto.ListingResponse{
			ID:          listing.ID,
			UserID:      listing.UserID,
			Price:       listing.Price,
			ListingType: listing.ListingType,
			CreatedAt:   listing.CreatedAt,
			UpdatedAt:   listing.UpdatedAt,
		})
	}

	offset := (pageNum - 1) * pageSize
	if offset >= len(filtered) {
		return []dto.ListingResponse{}, nil
	}

	end := offset + pageSize
	if end > len(filtered) {
		end = len(filtered)
	}

	return filtered[offset:end], nil
}

func (c *stubPublicListingClient) CreateListing(request dto.PublicCreateListingRequest) (dto.ListingResponse, error) {
	if request.ListingType != "rent" && request.ListingType != "sale" {
		return dto.ListingResponse{}, client.ErrListingRequestRejected
	}

	userFound := false
	for _, listing := range c.listings {
		if listing.UserID == request.UserID {
			userFound = true
			break
		}
	}
	if !userFound {
		return dto.ListingResponse{}, client.ErrListingRequestRejected
	}

	response := dto.ListingResponse{
		ID:          len(c.listings) + 1,
		UserID:      request.UserID,
		Price:       request.Price,
		ListingType: request.ListingType,
		CreatedAt:   1715007200,
		UpdatedAt:   1715007200,
	}
	c.listings = append(c.listings, domain.Listing{
		ID:          response.ID,
		UserID:      response.UserID,
		Price:       response.Price,
		ListingType: response.ListingType,
		CreatedAt:   response.CreatedAt,
		UpdatedAt:   response.UpdatedAt,
	})

	return response, nil
}

func TestHealthEndpoint(t *testing.T) {
	server := newServerWithDependencies(newStubUserRepository(), newStubListingRepository(), newStubPublicListingClient())

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
	server := newServerWithDependencies(newStubUserRepository(), newStubListingRepository(), newStubPublicListingClient())

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
	server := newServerWithDependencies(newStubUserRepository(), newStubListingRepository(), newStubPublicListingClient())

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
	server := newServerWithDependencies(newStubUserRepository(), newStubListingRepository(), newStubPublicListingClient())

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
	server := newServerWithDependencies(newStubUserRepository(), newStubListingRepository(), newStubPublicListingClient())

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
	server := newServerWithDependencies(newStubUserRepository(), newStubListingRepository(), newStubPublicListingClient())

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
	server := newServerWithDependencies(newStubUserRepository(), newStubListingRepository(), newStubPublicListingClient())

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
	server := newServerWithDependencies(newStubUserRepository(), newStubListingRepository(), newStubPublicListingClient())

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

func TestGetListingsEndpointPaginationAndFilter(t *testing.T) {
	server := newServerWithDependencies(newStubUserRepository(), newStubListingRepository(), newStubPublicListingClient())

	req := httptest.NewRequest("GET", "/listings?page_num=1&page_size=1&user_id=2", nil)
	resp, err := server.Test(req)
	if err != nil {
		t.Fatalf("filtered listings request failed: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	var payload struct {
		Result   bool `json:"result"`
		Listings []struct {
			ID          int    `json:"id"`
			UserID      int    `json:"user_id"`
			Price       int64  `json:"price"`
			ListingType string `json:"listing_type"`
		} `json:"listings"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		t.Fatalf("decode listings response: %v", err)
	}

	if !payload.Result {
		t.Fatalf("expected result true")
	}

	if len(payload.Listings) != 1 {
		t.Fatalf("expected 1 listing, got %d", len(payload.Listings))
	}

	if payload.Listings[0].UserID != 2 {
		t.Fatalf("expected listing user_id 2, got %d", payload.Listings[0].UserID)
	}
}

func TestCreateListingEndpointFormResponse(t *testing.T) {
	server := newServerWithDependencies(newStubUserRepository(), newStubListingRepository(), newStubPublicListingClient())

	body := "user_id=1&listing_type=rent&price=7000000"
	req := httptest.NewRequest("POST", "/listings", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := server.Test(req)
	if err != nil {
		t.Fatalf("create listing request failed: %v", err)
	}

	if resp.StatusCode != 201 {
		t.Fatalf("expected status 201, got %d", resp.StatusCode)
	}

	var payload struct {
		Result  bool `json:"result"`
		Listing struct {
			ID          int    `json:"id"`
			UserID      int    `json:"user_id"`
			Price       int64  `json:"price"`
			ListingType string `json:"listing_type"`
		} `json:"listing"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		t.Fatalf("decode create listing response: %v", err)
	}

	if !payload.Result {
		t.Fatalf("expected result true")
	}

	if payload.Listing.UserID != 1 {
		t.Fatalf("expected user_id 1, got %d", payload.Listing.UserID)
	}

	if payload.Listing.ListingType != "rent" {
		t.Fatalf("expected listing_type rent, got %s", payload.Listing.ListingType)
	}
}

func TestCreateListingEndpointInvalidListingTypeReturns422(t *testing.T) {
	server := newServerWithDependencies(newStubUserRepository(), newStubListingRepository(), newStubPublicListingClient())

	body := "user_id=1&listing_type=lease&price=7000000"
	req := httptest.NewRequest("POST", "/listings", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := server.Test(req)
	if err != nil {
		t.Fatalf("invalid listing type request failed: %v", err)
	}

	if resp.StatusCode != 422 {
		t.Fatalf("expected status 422, got %d", resp.StatusCode)
	}

	var payload struct {
		Return   bool `json:"return"`
		Listings any  `json:"listings"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		t.Fatalf("decode invalid listing type response: %v", err)
	}

	if payload.Return {
		t.Fatalf("expected return false")
	}

	if payload.Listings != nil {
		t.Fatalf("expected listings to be nil, got %#v", payload.Listings)
	}
}

func TestCreateListingEndpointUserNotFoundReturns422(t *testing.T) {
	server := newServerWithDependencies(newStubUserRepository(), newStubListingRepository(), newStubPublicListingClient())

	body := "user_id=999&listing_type=rent&price=7000000"
	req := httptest.NewRequest("POST", "/listings", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := server.Test(req)
	if err != nil {
		t.Fatalf("listing user not found request failed: %v", err)
	}

	if resp.StatusCode != 422 {
		t.Fatalf("expected status 422, got %d", resp.StatusCode)
	}

	var payload struct {
		Return   bool `json:"return"`
		Listings any  `json:"listings"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		t.Fatalf("decode listing user not found response: %v", err)
	}

	if payload.Return {
		t.Fatalf("expected return false")
	}

	if payload.Listings != nil {
		t.Fatalf("expected listings to be nil, got %#v", payload.Listings)
	}
}
