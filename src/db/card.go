package db

import (
	"config"

	"encoding/json"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"net/http"
)

type Card struct {
	ValidTill  int64  `json:"valid_till" bson:"valid_till"`
	ValidFrom  int64  `json:"valid_from" bson:"valid_from"`
	UserId     int32  `json:"user_id" bson:"user_id"`
	UserMobile string `json:"user_mobile" bson:"user_mobile"`
	UserName   string `json:"username" bson:"username"`
	CardNo     string `json:"card_no" bson:"card_no"`
	WiegandNo  string `json:"wiegand_no" bson:"wiegand_no"`
}

func IsValidCard(wiegand_no string) bool {
	count, err := db().C("cards").Find(bson.M{"wiegand_no": wiegand_no}).Count()
	if err != nil {
		fmt.Println(err)
	}
	if count > 0 {
		return true
	} else {
		return false
	}

}
func FindCard(wiegand_no string) {
	result := Card{}
	db().C("cards").Find(bson.M{"wiegand_no": wiegand_no}).One(&result)
	fmt.Println("Username:", result.UserName, result.UserMobile)
}

func CardCount() int {
	count, err := db().C("cards").Count()
	if err != nil {
		fmt.Println(err)
	}
	return count
}

type CardData struct {
	Cards []Card
}

func GetRemoteCards() []Card {

	url := fmt.Sprintf("%s/cards?device_identifier=%s&signature=%s", config.Host, getIdentifier(), getSignature())
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	var data CardData
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("card 数量：", len(data.Cards))
	return data.Cards
}

func AddCard(card Card) {
	err := db().C("cards").Insert(card)
	if err != nil {
		fmt.Print("保存出错：", err)
	}
}

func ClearCards() {
	db().C("cards").RemoveAll(nil)
}

func AddCards(cards []Card) {
	fmt.Println("添加卡数据")
	for _, card := range cards {
		err := db().C("cards").Insert(card)
		if err != nil {
			fmt.Print("保存出错：", err)
		}
	}
}
