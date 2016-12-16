package api

import (
	"config"
	"db"
	"driver"
	"fmt"
	"net/http"
)

func OpenByPassword(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	password := r.Form["code"][0]
	fmt.Println("密码：", password)
	valid_passcode := db.IsValidPasscode(password)
	fmt.Println("密码是否正确：", valid_passcode)
	if valid_passcode {
		if config.InDevice {
			driver.OpenDoor()
		}
		err := db.CreatePasswordAccessLog(password)
		db.UploadAccessLog()
		if err != nil {
			fmt.Println("create password access log error: ", err)
		}
		fmt.Fprintf(w, "open success")
	} else {
		fmt.Fprintf(w, "passcode error")
	}
}
