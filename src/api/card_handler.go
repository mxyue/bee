package api

import (
	"config"
	"db"
	"driver"
	"fmt"
	"net/http"
)

func OpenByCard(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	wiegandNo := r.Form["no"][0]
	fmt.Println("维根：", wiegandNo)
	if db.IsValidCard(wiegandNo) {
		if config.InDevice {
			driver.OpenDoor()
		}
		fmt.Fprintf(w, "open success")
	} else {
		fmt.Println("没有该卡数据")
		fmt.Fprintf(w, "card not found")
	}
}
