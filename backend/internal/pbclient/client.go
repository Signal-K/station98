package pbclient

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type Client struct {
	BaseURL string
	Token   string
}

func NewClient(baseURL string) *Client {
	return &Client{BaseURL: baseURL}
}

func (c *Client) Login(email, password string) error {
	body := map[string]string{"identity": email, "password": password}
	data, _ := json.Marshal(body)

	resp, err := http.Post(c.BaseURL+"/api/admins/auth-with-password", "application/json", bytes.NewReader(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var result struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	c.Token = result.Token
	return nil
}
