package db

import (
	"config"
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"io/ioutil"
	"net/http"
	"time"
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
	db, err := bolt.Open("bolt.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	defer db.Close()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	// 从只读事务中读取值.
	var card Card
	var present_flag bool
	if err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(DB_CARDS))
		card_bt := b.Get([]byte(wiegand_no))
		if len(card_bt) > 0 {
			present_flag = true
			err := json.Unmarshal(card_bt, &card)
			fmt.Printf("The card username is: %s\n", card.UserName)
			fmt.Printf("The card mobile is: %s\n", card.UserMobile)
			if err != nil {
				fmt.Println("json 解析错误：", err)
			}
		} else {
			present_flag = false
		}
		return nil
	}); err != nil {
		fmt.Println("db查询错误：", err)
	}
	// 关闭数据库释放锁
	if err := db.Close(); err != nil {
		fmt.Println("db关闭出错：", err)
	}
	return present_flag
}

func IsCardsPresent() bool {
	db, err := bolt.Open("bolt.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	defer db.Close()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	var count int
	if err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(DB_CARDS))
		count = b.Stats().KeyN
		fmt.Println("数据库卡数量：", count)
		return nil
	}); err != nil {
		fmt.Println(err)
	}
	// 关闭数据库释放锁
	if err := db.Close(); err != nil {
		fmt.Println(err)
	}
	return count > 0
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
	fmt.Println("服务端card 数量：", len(data.Cards))
	return data.Cards
}

func AddCards(cards []Card) {
	db, err := bolt.Open("bolt.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	defer db.Close()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println("添加card数据")
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(DB_CARDS))
		for _, card := range cards {
			enCard, err := json.Marshal(card)
			fmt.Println("server 卡维根：", card.WiegandNo)
			err = b.Put([]byte(card.WiegandNo), enCard)
			if err != nil {
				fmt.Print("保存出错：", err)
			}
		}
		return nil
	})
	// 关闭数据库释放锁
	if err := db.Close(); err != nil {
		fmt.Println(err)
	}
}

func CardStart() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
			time.Sleep(2 * time.Second)
			CardStart()
		}
	}()
	if !IsCardsPresent() {
		AddCards(GetRemoteCards())
	}
}
