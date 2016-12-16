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

type Passcode struct {
	ValidTill int64  `json:"valid_till"`
	ValidFrom int64  `json:"valid_from"`
	Id        int32  `json:"user_id"`
	Content   string `json:"passcode"`
}

func IsValidPasscode(number string) bool {
	db, err := bolt.Open("bolt.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	defer db.Close()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	var passwd string
	// 从只读事务中读取值.
	if err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(DB_PASSCODES))
		passwd = fmt.Sprintf("%s", b.Get([]byte("content")))
		fmt.Printf("The value of 'passcode' is: %s\n", passwd)
		return nil
	}); err != nil {
		fmt.Println(err)
	}
	// 关闭数据库释放锁
	if err := db.Close(); err != nil {
		fmt.Println(err)
	}
	return passwd == number
}

func IsPassCodePresent() bool {
	db, err := bolt.Open("bolt.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	defer db.Close()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	var passwd string
	// 从只读事务中读取值.
	if err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(DB_PASSCODES))
		passwd = fmt.Sprintf("%s", b.Get([]byte("content")))
		return nil
	}); err != nil {
		fmt.Println(err)
	}
	// 关闭数据库释放锁
	if err := db.Close(); err != nil {
		fmt.Println(err)
	}
	fmt.Printf("-数据库门密码: %s\n", passwd)
	return passwd != ""
}

type PasscodeData struct {
	Passcodes []Passcode
}

func GetRemotePasscodes() []Passcode {
	url := fmt.Sprintf("%s/passcodes?device_identifier=%s&signature=%s", config.Host, config.GetIdentifier(), config.GetSignature())
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
		var data PasscodeData
		err = json.Unmarshal(body, &data)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("passcode 数量：", len(data.Passcodes))
		return data.Passcodes
	} else {
		fmt.Println("passcode 获取远程数据出错 code:", resp.StatusCode)
		var passcodes []Passcode
		return passcodes
	}

}

func AddPasscodes(passcodes []Passcode) error {
	db, err := bolt.Open("bolt.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	defer db.Close()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println("修改密码数据")
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(DB_PASSCODES))
		for _, passcode := range passcodes {
			fmt.Println(passcode.Content)
			err := b.Put([]byte("content"), []byte(passcode.Content))
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
	return err
}

func DonePasscodesSync() error {
	client := &http.Client{}
	url := fmt.Sprintf("%s/passcodes/done_passcodes_sync?device_identifier=%s&signature=%s", config.Host, config.GetIdentifier(), config.GetSignature())
	req, err := http.NewRequest("PUT", url, nil)
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("更新服务器密码同步状态 code:", res.StatusCode)
	return err
}

func PassCodeStart() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
			time.Sleep(4 * time.Second)
			PassCodeStart()
		}
	}()
	if !IsPassCodePresent() {
		AddPasscodes(GetRemotePasscodes())
		go DonePasscodesSync()
	}
}
