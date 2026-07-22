package iyzico

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	baseURL    string
	apiKey     string
	secretKey  string
	httpClient *http.Client
}

func NewClient(baseURL string, apiKey string, secretKey string) *Client {
	return &Client{
		baseURL:   baseURL,
		apiKey:    apiKey,
		secretKey: secretKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) Initialize3DS(ctx context.Context,req Init3DSRequest) (*Init3DSResponse, error) {

	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	var response Init3DSResponse

	err = c.doRequest(ctx,http.MethodPost,"payment/3dsecure/initialize",body,&response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *Client) doRequest(ctx context.Context, method string, path string, body []byte, result interface{}) error {

	randomKey := fmt.Sprintf("%d", time.Now().UnixNano())

	authorization := generateAuthorization(c.apiKey, c.secretKey, randomKey, path, body)

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, bytes.NewReader(body))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", authorization)
	req.Header.Set("x-iyzi-rnd",randomKey)
	req.Header.Set("Content-Type","application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(result); 
	if err != nil {
		return err
	}

	return nil
}
