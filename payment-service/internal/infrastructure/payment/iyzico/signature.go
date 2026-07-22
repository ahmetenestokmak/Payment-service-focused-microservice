package iyzico

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

func GenerateAuthorization(apiKey string,secretKey string,randomKey string,requestString string,) string {

	signatureData := randomKey + requestString

	hash := sha256.Sum256([]byte(secretKey + signatureData),)

	signature := fmt.Sprintf("%x", hash)

	authString := fmt.Sprintf(
		"%s:%s:%s",
		apiKey,
		randomKey,
		signature,
	)

	encoded := base64.StdEncoding.EncodeToString(
		[]byte(authString),
	)

	return "IYZWSv2 " + encoded
}
