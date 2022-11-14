package clients

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client .
type Client struct {
	token      string
	baseURL    string
	HTTPClient *http.Client
}

// NewClient .
func NewClient(BaseURL string, token string) *Client {
	return &Client{
		HTTPClient: &http.Client{
			Timeout: 5 * time.Minute,
		},
		baseURL: BaseURL,
		token:   token,
	}
}

type errorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type successResponse struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}

func (c *Client) sendRequest(req *http.Request, v interface{}) error {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json; charset=utf-8")
	req.Header.Set("Authorization", c.token)

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(res.Body)

	// Try to unmarshall into errorResponse
	if res.StatusCode != http.StatusOK {
		var errRes errorResponse
		if err = json.NewDecoder(res.Body).Decode(&errRes); err == nil {
			return errors.New(errRes.Message)
		}

		return fmt.Errorf("unknown error, status code: %d", res.StatusCode)
	}

	// Unmarshall and populate v
	fullResponse := successResponse{
		Code: res.StatusCode,
		Data: v,
	}
	if err = json.NewDecoder(res.Body).Decode(&fullResponse.Data); err != nil {
		return err
	}
	return nil
}
