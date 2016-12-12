package main

import (
	"api"
	"config"
	"db"
	"driver"
	"fmt"
	"mqtt_client"
	"net/http"
)

func main() {
	go db.PassCodeStart()
	go db.CardStart()
	go db.DeviceInfoStart()
	go mqtt_client.Start()
	if config.InDevice {
		go driver.KeyStart()
	}
	err := http.ListenAndServe(":9090", api.Route())
	if err != nil {
		fmt.Println("listen and server err:", err)
	} else {
		fmt.Println("run in 9090")
	}

}
