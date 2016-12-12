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
	card, card_present := db.ValidCard(wiegandNo)
	if card_present {
		if config.InDevice {
			driver.OpenDoor()
		}
		err := db.CreateCardAccessLog(card)
		db.UploadAccessLog()
		if err != nil {
			fmt.Println("create card access log error: ", err)
		}
		fmt.Fprintf(w, "open success")
	} else {
		fmt.Println("not found the card")
		fmt.Fprintf(w, "card not found")
	}
}
