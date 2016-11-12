package db

import (
	"config"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"gopkg.in/mgo.v2"
	"time"
)

func db() *mgo.Database {

	url := "mongodb://127.0.0.1:27017"
	session, err := mgo.Dial(url)
	if err != nil {
		fmt.Println(err)
	}
	return session.DB("bee")
}

func getIdentifier() string {
	var deviceIdentifier = fmt.Sprintf("%s:%d", config.Identifier, time.Now().Unix())
	return deviceIdentifier
}

func getSignature() string {
	key := []byte(config.Secret)
	mac := hmac.New(sha1.New, key)
	mac.Write([]byte(getIdentifier()))
	signature := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	return signature
}
