package service

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
)

func VerifyMidtransSignature(
	orderID string,
	statusCode string,
	grossAmount string,
	serverKey string,
	signature string,
) bool {
	raw := orderID + statusCode + grossAmount + serverKey

	hash := sha512.Sum512([]byte(raw))
	expectedSignature := hex.EncodeToString(hash[:])

	return hmac.Equal(
		[]byte(expectedSignature),
		[]byte(signature),
	)
}
