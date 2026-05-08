package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"property-api/internal/dto"
)

var ErrListingRequestRejected = errors.New("listing request rejected")

type ListingClient interface {
	GetListings(pageNum int, pageSize int, userID *int) ([]dto.ListingResponse, error)
	CreateListing(request dto.PublicCreateListingRequest) (dto.ListingResponse, error)
}

type HTTPListingClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewHTTPListingClient(baseURL string) *HTTPListingClient {
	return &HTTPListingClient{
		baseURL: strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *HTTPListingClient) GetListings(pageNum int, pageSize int, userID *int) ([]dto.ListingResponse, error) {
	query := url.Values{}
	query.Set("page_num", strconv.Itoa(pageNum))
	query.Set("page_size", strconv.Itoa(pageSize))
	if userID != nil {
		query.Set("user_id", strconv.Itoa(*userID))
	}

	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("%s/listings?%s", c.baseURL, query.Encode()),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("build get listings request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send get listings request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get listings returned status %d", resp.StatusCode)
	}

	var payload dto.GetListingsResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("decode get listings response: %w", err)
	}

	return payload.Listings, nil
}

func (c *HTTPListingClient) CreateListing(request dto.PublicCreateListingRequest) (dto.ListingResponse, error) {
	form := url.Values{}
	form.Set("user_id", strconv.Itoa(request.UserID))
	form.Set("listing_type", request.ListingType)
	form.Set("price", strconv.FormatInt(request.Price, 10))

	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%s/listings", c.baseURL),
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		return dto.ListingResponse{}, fmt.Errorf("build create listing request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return dto.ListingResponse{}, fmt.Errorf("send create listing request: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusCreated:
		var payload dto.CreateListingResponse
		if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
			return dto.ListingResponse{}, fmt.Errorf("decode create listing response: %w", err)
		}
		if payload.Listing == nil {
			return dto.ListingResponse{}, fmt.Errorf("create listing response missing listing")
		}
		return *payload.Listing, nil
	case http.StatusUnprocessableEntity:
		return dto.ListingResponse{}, ErrListingRequestRejected
	default:
		return dto.ListingResponse{}, fmt.Errorf("create listing returned status %d", resp.StatusCode)
	}
}
