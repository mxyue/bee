package config

import (
	"bufio"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	Version  = "0.1"
	MqttPort = "1885"
	InDevice = false
)

var MqttHost, Host, Identifier, Secret string
var txtConfigs = make(map[string]string)

func init() {
	err := readLines(configPath())
	if err != nil {
		fmt.Println("读取configs.txt错误：", err)
	}
	Identifier = txtConfigs["identifier"]
	Secret = txtConfigs["secret"]
	MqttHost = txtConfigs["mqtt_host"]
	Host = txtConfigs["host"]
	if MqttHost == "" {
		MqttHost = "voip.didikon.com"
	}
	if Host == "" {
		Host = "http://api.didikon.com:80/device"
	}
}

func configPath() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println(err)
	}
	return fmt.Sprintf("%s/configs.txt", dir)
}

func readLines(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		setConf(scanner.Text())
	}
	return scanner.Err()
}

func setConf(text string) {
	if strings.HasPrefix(text, "#") {
		return
	}
	temp_arr := strings.Split(text, "=")
	if len(temp_arr) == 2 {
		txtConfigs[temp_arr[0]] = temp_arr[1]
	}
}

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
