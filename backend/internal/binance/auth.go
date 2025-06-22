package binance

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"time"
)

// signRequest adds the required signature to the request parameters
func (c *Client) signRequest(params url.Values) error {
	if c.secretKey == "" {
		return nil // No signature needed for public endpoints
	}

	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	params.Set("timestamp", fmt.Sprintf("%d", timestamp))

	signature := generateSignature(params.Encode(), c.secretKey)
	params.Set("signature", signature)

	return nil
}

// generateSignature creates a SHA256 HMAC signature
func generateSignature(payload, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(payload))
	return hex.EncodeToString(h.Sum(nil))
}
