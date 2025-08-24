package pbclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Client struct {
	BaseURL string
	Token   string
}

func NewClient(baseURL string) *Client {
	return &Client{BaseURL: baseURL}
}

// Login authenticates with PocketBase and sets the admin token
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

// CreateRecord inserts a new record into a collection
func (c *Client) CreateRecord(collection string, data map[string]interface{}) (*map[string]interface{}, error) {
	url := fmt.Sprintf("%s/api/collections/%s/records", c.BaseURL, collection)

	payload, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	req.Header.Set("Authorization", c.Token)
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// FindRecordByField checks if a record already exists by field match
func (c *Client) FindRecordByField(collection, field string, value interface{}) (*map[string]interface{}, error) {
	// For title field, normalize to lowercase and trim spaces for case-insensitive matching
	strVal, isStr := value.(string)
	if field == "title" && isStr {
		strVal = strings.ToLower(strings.TrimSpace(strVal))
		// Use LIKE for case-insensitive match (PocketBase supports basic LIKE)
		url := fmt.Sprintf("%s/api/collections/%s/records?filter=(LOWER(TRIM(%s))='%s')", c.BaseURL, collection, field, strVal)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Authorization", c.Token)
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()
		var data struct {
			Items []map[string]interface{} `json:"items"`
		}
		if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
			return nil, err
		}
		if len(data.Items) > 0 {
			return &data.Items[0], nil
		}
		return nil, nil
	}
	// Default: exact match
	url := fmt.Sprintf("%s/api/collections/%s/records?filter=(%s='%v')", c.BaseURL, collection, field, value)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", c.Token)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var data struct {
		Items []map[string]interface{} `json:"items"`
	}
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return nil, err
	}
	if len(data.Items) > 0 {
		return &data.Items[0], nil
	}
	return nil, nil
}
