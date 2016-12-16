package db

import (
	"config"
	"encoding/json"
	"errors"
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

func ValidCard(wiegand_no string) (Card, bool) {
	db, err := bolt.Open("bolt.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	defer db.Close()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	// 从只读事务中读取值.
	var card Card
	present_flag := false
	if err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(DB_CARDS))
		card_bt := b.Get([]byte(wiegand_no))
		if len(card_bt) > 0 {
			err := json.Unmarshal(card_bt, &card)
			if err != nil {
				fmt.Println("json 解析错误：", err)
			}
			now := time.Now().Unix()
			if card.ValidTill > now {
				present_flag = true
			} else {
				fmt.Println("卡片有效期已过,有效时间戳:", card.ValidTill)
			}
			fmt.Printf("The card username is: %s\n", card.UserName)
			fmt.Printf("The card mobile is: %s\n", card.UserMobile)
		}
		return nil
	}); err != nil {
		fmt.Println("db查询错误：", err)
	}
	if err := db.Close(); err != nil {
		fmt.Println("db关闭出错：", err)
	}
	return card, present_flag
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
		fmt.Println("-数据库卡数量：", count)
		return nil
	}); err != nil {
		fmt.Println(err)
	}
	if err := db.Close(); err != nil {
		fmt.Println(err)
	}
	return count > 0
}

type CardData struct {
	Cards []Card
}

func GetRemoteCards() ([]Card, error) {
	url := fmt.Sprintf("%s/cards?device_identifier=%s&signature=%s", config.Host, config.GetIdentifier(), config.GetSignature())
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	if resp.StatusCode == 200 {
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
		return data.Cards, err
	} else {
		fmt.Println("card 获取远程数据出错 code:", resp.StatusCode)
		var cards []Card
		errStr := fmt.Sprintf("error status code: %s", resp.StatusCode)
		err := errors.New(errStr)
		return cards, err
	}
}

func AddCards(cards []Card) error {
	db, err := bolt.Open("bolt.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	defer db.Close()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(DB_CARDS))
		for _, card := range cards {
			enCard, err := json.Marshal(card)
			fmt.Println("添加卡 wiegangdNo：", card.WiegandNo)
			err = b.Put([]byte(card.WiegandNo), enCard)
			if err != nil {
				fmt.Print("保存出错：", err)
			}
		}
		return nil
	})
	if err := db.Close(); err != nil {
		fmt.Println(err)
	}
	return err
}

func DeleteCards(cards []Card) error {
	db, err := bolt.Open("bolt.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	defer db.Close()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(DB_CARDS))
		for _, card := range cards {
			fmt.Println("删除card数据 wiegandNo:", card.WiegandNo)
			err = b.Delete([]byte(card.WiegandNo))
			if err != nil {
				fmt.Print("删除出错：", err)
			}
		}
		return err
	})
	if err := db.Close(); err != nil {
		fmt.Println(err)
	}
	return err
}

func ResetCards() error {
	db, err := bolt.Open("bolt.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	defer db.Close()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	cards, err := GetRemoteCards()
	if err != nil {
		fmt.Println(err)
	} else {
		db.Update(func(tx *bolt.Tx) error {
			err := tx.DeleteBucket([]byte(DB_CARDS))
			_, err = tx.CreateBucket([]byte(DB_CARDS))
			return err
		})
		if err := db.Close(); err != nil {
			fmt.Println(err)
		}
		err = AddCards(cards)
		if err == nil {
			err = DoneCardsSync()
		}
	}
	return err
}

func DoneCardsSync() error {
	client := &http.Client{}
	url := fmt.Sprintf("%s/cards/done_cards_sync?device_identifier=%s&signature=%s", config.Host, config.GetIdentifier(), config.GetSignature())
	req, err := http.NewRequest("PUT", url, nil)
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("更新服务器门禁卡同步状态 code:", res.StatusCode)
	return err
}

func CardStart() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
			time.Sleep(5 * time.Second)
			CardStart()
		}
	}()
	if !IsCardsPresent() {
		cards, err := GetRemoteCards()
		if err != nil {
			fmt.Println(err)
		} else {
			AddCards(cards)
			go DoneCardsSync()
		}
	}
}
