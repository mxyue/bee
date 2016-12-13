package db

import (
	// "bytes"
	"config"
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	// "io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	// "strings"
	"time"
)

type AccessLog struct {
	AccessTimestamp int64  `json:"access_timestamp"`
	UserId          int32  `json:"user_id"`
	UserMobile      string `json:"user_mobile"`
	UserName        string `json:"username"`
	CardNo          string `json:"card_no"`
	AccessType      string `json:"access_type"`
	Door            string `json:"door"`
	DoorId          int32  `json:"door_id"`
	Unit            string `json:"unit"`
	UnitId          int32  `json:"unit_id"`
	Building        string `json:"building"`
	BuildingId      int32  `json:"building_id"`
	Compound        string `json:"compound"`
	CompoundId      int32  `json:"compound_id"`
}

func CreateCardAccessLog(card Card) error {
	deviceInfo, err := GetDeviceInfo()
	if err != nil {
		fmt.Println("err:", err)
	}
	accessLog := AccessLog{
		AccessTimestamp: time.Now().Unix(),
		UserId:          card.UserId,
		UserName:        card.UserName,
		UserMobile:      card.UserMobile,
		CardNo:          card.CardNo,
		AccessType:      "card",
		Door:            deviceInfo.DoorName,
		DoorId:          deviceInfo.DoorId,
		Unit:            deviceInfo.UnitName,
		UnitId:          deviceInfo.UnitId,
		Building:        deviceInfo.BuildingName,
		BuildingId:      deviceInfo.BuildingId,
		Compound:        deviceInfo.CompoundName,
		CompoundId:      deviceInfo.CompoundId,
	}
	err = saveAccessLog(accessLog)
	return err
}
func CreatePasswordAccessLog(code string) error {
	deviceInfo, err := GetDeviceInfo()
	if err != nil {
		fmt.Println("err:", err)
	}
	accessLog := AccessLog{
		AccessTimestamp: time.Now().Unix(),
		AccessType:      "password",
		Door:            deviceInfo.DoorName,
		DoorId:          deviceInfo.DoorId,
		Unit:            deviceInfo.UnitName,
		UnitId:          deviceInfo.UnitId,
		Building:        deviceInfo.BuildingName,
		BuildingId:      deviceInfo.BuildingId,
		Compound:        deviceInfo.CompoundName,
		CompoundId:      deviceInfo.CompoundId,
	}
	err = saveAccessLog(accessLog)
	return err
}

func saveAccessLog(accessLog AccessLog) error {
	db, err := bolt.Open("bolt.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	defer db.Close()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(DB_ACCESS_LOGS))
		enAccessLog, err := json.Marshal(accessLog)
		err = b.Put([]byte(strconv.FormatInt(time.Now().Unix(), 10)), enAccessLog)
		if err != nil {
			fmt.Print("保存出错：", err)
		}
		return nil
	})
	// 关闭数据库释放锁
	if err := db.Close(); err != nil {
		fmt.Println(err)
	}
	return err
}

func GetAccessLog() ([]byte, [][]byte, error) {
	db, err := bolt.Open("bolt.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	defer db.Close()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	var accessLogs []AccessLog
	var keys [][]byte
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(DB_ACCESS_LOGS))
		if err := b.ForEach(func(k, v []byte) error {
			keys = append(keys, k)
			var accessLog AccessLog
			err := json.Unmarshal(v, &accessLog)
			accessLogs = append(accessLogs, accessLog)
			return err
		}); err != nil {
			return err
		}
		return nil
	})
	// 关闭数据库释放锁
	if err := db.Close(); err != nil {
		fmt.Println(err)
	}
	btAccessLogs, err := json.Marshal(accessLogs)
	return btAccessLogs, keys, err
}

type SuccessBody struct {
	Success bool `json:"success"`
}

func UploadAccessLog() error {
	btAccessLogs, keys, err := GetAccessLog()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(btAccessLogs))
	logUrl := fmt.Sprintf("%s/access_logs?device_identifier=%s&signature=%s", config.Host, config.GetIdentifier(), config.GetSignature())
	res, err := http.PostForm(logUrl, url.Values{"logs": {string(btAccessLogs)}})
	fmt.Println("upload access log :", res.StatusCode)
	if res.StatusCode == 200 {
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		var data SuccessBody
		err = json.Unmarshal(body, &data)
		fmt.Println("日志上传：", data.Success)
		if data.Success {
			RemoveAccessLogs(keys)
		}
		return err
	} else {
		return err
	}
}

func RemoveAccessLogs(keys [][]byte) error {
	db, err := bolt.Open("bolt.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	defer db.Close()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	if err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(DB_ACCESS_LOGS))
		var err error
		for _, key := range keys {
			err = b.Delete(key)
			fmt.Println("删除key：", string(key))
		}
		return err
	}); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
