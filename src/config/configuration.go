package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	Version  = "0.1"
	Host     = "http://api.didikon.com:80/device"
	MqttHost = "voip.didikon.com"
	MqttPort = "1885"
	InDevice = false
)

var txtConfigs = make(map[string]string)
var Identifier string
var Secret string

func init() {
	err := readLines(configPath())
	if err != nil {
		fmt.Println("读取configs.txt错误：", err)
	}
	Identifier = txtConfigs["identifier"]
	Secret = txtConfigs["secret"]
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
	temp_arr := strings.Split(text, "=")
	if len(temp_arr) == 2 {
		txtConfigs[temp_arr[0]] = temp_arr[1]
	}
}
