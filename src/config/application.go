package config

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"time"
)

func GetIdentifier() string {
	var deviceIdentifier = fmt.Sprintf("%s:%d", Identifier, time.Now().Unix())
	return deviceIdentifier
}

func GetSignature() string {
	key := []byte(Secret)
	mac := hmac.New(sha1.New, key)
	mac.Write([]byte(GetIdentifier()))
	signature := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	return signature
}
