package rickmorty

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"

	"1337b04rd/internal/domain/avatar"
)

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

type characterListResponse struct {
	Info struct {
		Count int `json:"count"`
	} `json:"info"`
}

type character struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Image string `json:"image"`
}

func (c *Client) GetRandomAvatar() (*avatar.Avatar, error) {
	count, err := c.getCharacterCount()
	if err != nil {
		return nil, err
	}

	randomID := rand.Intn(count) + 1

	character, err := c.fetchCharacter(randomID)
	if err != nil {
		return nil, err
	}

	return &avatar.Avatar{
		URL:         character.Image,
		DisplayName: character.Name,
	}, nil
}

func (c *Client) getCharacterCount() (int, error) {
	resp, err := c.httpClient.Get(fmt.Sprintf("%s/character", c.baseURL))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("rickmorty API: status %d", resp.StatusCode)
	}

	var data characterListResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, err
	}

	return data.Info.Count, nil
}

func (c *Client) fetchCharacter(id int) (*character, error) {
	resp, err := c.httpClient.Get(fmt.Sprintf("%s/character/%d", c.baseURL, id))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("rickmorty API: status %d", resp.StatusCode)
	}

	var data character
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	if data.Image == "" || data.Name == "" {
		return nil, errors.New("invalid character response")
	}

	return &data, nil
}
