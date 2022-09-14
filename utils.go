package authorizer

import (
	"encoding/base64"
	"math/rand"
)

// CreateRandomString returns a random string 43 characters
func CreateRandomString() string {
	letterBytes := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, 43)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

// EncodeB64 returns string which is base64 encoded
func EncodeB64(message string) string {
	base64Text := make([]byte, base64.StdEncoding.EncodedLen(len(message)))
	base64.StdEncoding.Encode(base64Text, []byte(message))
	return string(base64Text)
}

// NewStringRef returns a reference to a string with given value
func NewStringRef(v string) *string {
	return &v
}

// StringValue returns the value of the given string ref
func StringValue(r *string, defaultValue ...string) string {
	if r != nil {
		return *r
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return ""
}
