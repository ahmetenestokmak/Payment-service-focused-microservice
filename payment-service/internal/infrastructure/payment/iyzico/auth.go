package iyzico

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
)



func generateAuthorization(apiKey string, secretKey string, randomKey string, uriPath string, body []byte) string {

	payload := fmt.Sprintf("%s%s%s", randomKey, uriPath, string(body))

	mac := hmac.New(sha256.New, []byte(secretKey))

	mac.Write([]byte(payload))

	signature := fmt.Sprintf("%x", mac.Sum(nil))

	authorizationString := fmt.Sprintf("apiKey:%s&randomKey:%s&signature:%s", apiKey, randomKey, signature)

	encoded := base64.StdEncoding.EncodeToString([]byte(authorizationString))

	return "IYZWSv2 " + encoded
}

func (c *Client) Auth3DS(ctx context.Context, req Auth3DSRequest) (*Auth3DSResponse, error) {

	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	var response Auth3DSResponse

	err = c.doRequest(
		ctx,
		http.MethodPost,
		"/payment/3dsecure/auth",
		body,
		&response,
	)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
