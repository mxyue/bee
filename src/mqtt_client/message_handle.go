package mqtt_client

import (
	"config"
	"db"
	"driver"
	"encoding/json"
	"fmt"
)

type CardBody struct {
	Cards     []db.Card `json:"cards"`
	Operation string    `json:"operation"`
}

type PasscodeBody struct {
	Passcodes []db.Passcode `json:"passcodes"`
	Operation string        `json:"operation"`
}

func cardOperation(strParams []byte) error {
	var cardBody CardBody
	err := json.Unmarshal(strParams, &cardBody)
	if cardBody.Operation == "append" {
		err = db.AddCards(cardBody.Cards)
	} else if cardBody.Operation == "remove" {
		err = db.DeleteCards(cardBody.Cards)
	} else if cardBody.Operation == "batch" {
		go db.ResetCards()
	}
	return err
}

func passcodeOperation(strParams []byte) error {
	var passcodeBody PasscodeBody
	err := json.Unmarshal(strParams, &passcodeBody)
	if passcodeBody.Operation == "replaceAll" {
		err = db.AddPasscodes(passcodeBody.Passcodes)
	} else if passcodeBody.Operation == "replace" {
		err = db.AddPasscodes(passcodeBody.Passcodes)
		if err == nil {
			go db.DonePasscodesSync()
		}
	}
	return err
}

func openGate() error {
	if config.InDevice {
		go driver.OpenDoor()
	} else {
		fmt.Println("模拟开门成功")
	}
	return nil
}
