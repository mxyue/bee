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

const DEVICE_INFO_KEY = "device_info"

type DeviceInfo struct {
	Name                 string   `json:"name"`
	CompoundId           int32    `json:"compound_id"`
	CompoundName         string   `json:"compound_name"`
	BuildingId           int32    `json:"building_id"`
	BuildingName         string   `json:"building_name"`
	UnitId               int32    `json:"unit_id"`
	UnitName             string   `json:"unit_name"`
	DoorId               int32    `json:"door_id"`
	DoorName             string   `json:"door_name"`
	DoorType             string   `json:"door_type"`
	ServerAddress        string   `json:"server_address"`
	SnmpServerAddress    string   `json:"snmp_server_address"`
	ReconnectionInterval int32    `json:"reconnection_interval"`
	PingInterval         int32    `json:"ping_interval"`
	ConnectionTimeout    int32    `json:"connection_timeout"`
	SendVideoBitRate     int32    `json:"send_video_bit_rate"`
	SendAudioBitRate     int32    `json:"send_audio_bit_rate"`
	City                 string   `json:"city"`
	CityCode             string   `json:"city_code"`
	Topics               []string `json:"topics"`
}

func AddDeviceInfo(info DeviceInfo) {
	db, err := bolt.Open("bolt.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	defer db.Close()
	if err != nil {
		fmt.Println("AddDeviceInfo:", err)
		panic(err)
	}
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(DB_DEVICE_INFOS))
		fmt.Println("添加设备信息", info.Name)
		enInfo, err := json.Marshal(info)
		err = b.Put([]byte(DEVICE_INFO_KEY), enInfo)
		if err != nil {
			fmt.Print("保存出错：", err)
		}
		return nil
	})
	// 关闭数据库释放锁
	if err := db.Close(); err != nil {
		fmt.Println("close db", err)
	}
}

func GetDeviceInfo() (DeviceInfo, error) {
	db, err := bolt.Open("bolt.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	defer db.Close()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	var deviceInfo DeviceInfo
	// 从只读事务中读取值.
	if err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(DB_DEVICE_INFOS))
		bt_device := b.Get([]byte(DEVICE_INFO_KEY))
		err := json.Unmarshal(bt_device, &deviceInfo)
		return err
	}); err != nil {
		fmt.Println(err)
	}
	// 关闭数据库释放锁
	if err := db.Close(); err != nil {
		fmt.Println(err)
	}
	fmt.Println("数据库门名字：", deviceInfo.DoorName)
	return deviceInfo, err
}

func GetRemoteDeviceInfo() DeviceInfo {
	url := fmt.Sprintf("%s/configuration?device_identifier=%s&signature=%s", config.Host, config.GetIdentifier(), config.GetSignature())
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
		var data DeviceInfo
		err = json.Unmarshal(body, &data)
		if err != nil {
			fmt.Println("device info:", err)
		}
		return data
	} else {
		fmt.Println("获取远程device info 数据出错 code:", resp.StatusCode)
		var info DeviceInfo
		return info
	}
}

func DeviceInfoStart() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("DeviceInfoStart:", err)
			time.Sleep(5 * time.Second)
			DeviceInfoStart()
		}
	}()
	AddDeviceInfo(GetRemoteDeviceInfo())
}
