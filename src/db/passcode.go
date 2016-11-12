package db

import (
	"config"
	"encoding/json"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"net/http"
)

type Passcode struct {
	ValidTill int64  `json:"valid_till" bson:"valid_till"`
	ValidFrom int64  `json:"valid_from" bson:"valid_from"`
	Id        int32  `json:"user_id" bson:"user_id"`
	Content   string `json:"passcode" bson:"content"`
}

func IsValidPasscode(number string) bool {
	count, err := db().C("passcodes").Find(bson.M{"content": number}).Count()
	if err != nil {
		fmt.Println(err)
	}
	if count > 0 {
		return true
	} else {
		return false
	}

}

func PasscodeCount() int {
	count, err := db().C("passcodes").Count()
	if err != nil {
		fmt.Println(err)
	}
	return count
}

type PasscodeData struct {
	Passcodes []Passcode
}

func GetRemotePasscodes() []Passcode {

	url := fmt.Sprintf("%s/passcodes?device_identifier=%s&signature=%s", config.Host, getIdentifier(), getSignature())
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	var data PasscodeData
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("passcode 数量：", len(data.Passcodes))
	return data.Passcodes
}

func AddPasscode(passcode Passcode) {
	err := db().C("passcodes").Insert(passcode)
	if err != nil {
		fmt.Print("保存出错：", err)
	}
}

func ClearPasscodes() {
	db().C("passcodes").RemoveAll(nil)
}

func AddPasscodes(passcodes []Passcode) {
	fmt.Println("添加密码数据")
	for _, passcode := range passcodes {
		err := db().C("passcodes").Insert(passcode)
		if err != nil {
			fmt.Print("保存出错：", err)
		}
	}
}
