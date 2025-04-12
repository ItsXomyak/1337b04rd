package rickmorty

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"1337b04rd/internal/app/common/logger"
)

type character struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Image string `json:"image"`
}

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string, httpClient *http.Client) *Client {
	return &Client{
		baseURL:    baseURL,
		httpClient: httpClient,
	}
}

func (c *Client) FetchCharacterByID(id int) (*character, error) {
	url := fmt.Sprintf("%s/character/%d", c.baseURL, id)
	resp, err := c.httpClient.Get(url)
	if err != nil {
		logger.Error("failed to send GET request to Rick and Morty API", "url", url, "error", err)
		return nil, fmt.Errorf("rickmorty GET error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Error("rickmorty API returned non-200 status", "url", url, "status", resp.StatusCode)
		return nil, fmt.Errorf("rickmorty API status: %d", resp.StatusCode)
	}

	var data character
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		logger.Error("failed to decode response from Rick and Morty API", "url", url, "error", err)
		return nil, fmt.Errorf("decode error: %w", err)
	}

	if data.Image == "" || data.Name == "" {
		logger.Error("rickmorty API returned character with missing fields", "id", id)
		return nil, errors.New("character missing name or image")
	}

	return &data, nil
}
