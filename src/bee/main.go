package main

import (
	"api"
	"db"
	"driver"
	"fmt"
	"mqtt_client"
	"net/http"
)

func main() {
	if db.CardCount() == 0 {
		db.AddPasscodes(db.GetRemotePasscodes())
	}

	mqtt_client.Start()

	if db.PasscodeCount() == 0 {
		db.AddCards(db.GetRemoteCards())
	}

	driver.KeyStart()
	fmt.Println("run in 9090")
	err := http.ListenAndServe(":9090", api.Route())
	if err != nil {
		fmt.Println("listen and server err:", err)
	}

}
