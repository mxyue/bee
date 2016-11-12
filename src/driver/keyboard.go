package driver

import (
	"db"
	"fmt"
	"os"
	"rpio"
	"strconv"
	"time"
)

var lmap = map[rpio.Pin]int{
	pin_l1: 0,
	pin_l2: 1,
	pin_l3: 2,
	pin_l4: 3,
}
var rmap = map[rpio.Pin]int{
	pin_r1: 0,
	pin_r2: 1,
	pin_r3: 2,
	pin_r4: 3,
}

const (
	pin_l1 = rpio.Pin(6)
	pin_l2 = rpio.Pin(13)
	pin_l3 = rpio.Pin(19)
	pin_l4 = rpio.Pin(26)

	pin_r1 = rpio.Pin(12)
	pin_r2 = rpio.Pin(16)
	pin_r3 = rpio.Pin(20)
	pin_r4 = rpio.Pin(21)

	door_pin = rpio.Pin(27)
)

var numberArr = [4][4]int{
	{1, 2, 3, 11},
	{4, 5, 6, 12},
	{7, 8, 9, 13},
	{14, 0, 15, 16},
}

var (
	larr = []rpio.Pin{pin_l1, pin_l2, pin_l3, pin_l4}
	rarr = []rpio.Pin{pin_r1, pin_r2, pin_r3, pin_r4}
)
var numberStr string

func KeyStart() {
	fmt.Println("keyboard >>")
	first_step()
}

func first_step() {
	rpio.Close()
	if err := rpio.Open(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, channel := range rarr {
		channel.Output()
		channel.High()
	}
	for _, channel := range larr {
		channel.Input()
		channel.PullDown()
	}
	loop_listen()
}

func loop_listen() {
	var yTemp rpio.Pin
	for true {
		flag := false
		for _, channel := range larr {
			if channel.Read() == 1 {
				flag = true
				yTemp = channel
				break
			}
		}
		if flag {
			break
		}
		time.Sleep(5 * time.Second / 100)
	}
	second_step(yTemp)
}

func second_step(pin rpio.Pin) {
	pin.Output()
	pin.High()
	for _, channel := range rarr {
		channel.Input()
		channel.PullDown()
	}
	var xtemp rpio.Pin
	for _, channel := range rarr {
		if channel.Read() == 1 {
			xtemp = channel
			break
		}
	}
	for true {
		if xtemp.Read() == 0 {
			break
		}
		time.Sleep(5 * time.Second / 100)
	}
	var yPin = lmap[pin]
	var xPin = rmap[xtemp]
	// fmt.Println("x:", xPin, ";", "y:", yPin)
	var number = numberArr[yPin][xPin]
	fmt.Println("number: ", number)
	rpio.Close()
	dealNumber(number)
}

func dealNumber(number int) {
	if number >= 0 && number < 10 {
		numberStr = numberStr + strconv.Itoa(number)
	} else if number == 11 {
		if len(numberStr) > 0 {
			newLen := len(numberStr) - 1
			numberStr = numberStr[0:newLen]
		}
	} else if number == 12 {
		if numberStr == "" {
			fmt.Println("密码为空")
		} else {
			valid_passcode := db.IsValidPasscode(numberStr)
			fmt.Println("密码是否正确：", valid_passcode)
			if valid_passcode {
				go OpenDoor()
			}
			numberStr = ""

		}
	}
	fmt.Println("当前数：", numberStr)
	first_step()
}

func OpenDoor() {
	rpio.Close()
	if err := rpio.Open(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	door_pin.Output()
	door_pin.High()
	time.Sleep(5 * time.Second)
	door_pin.Low()
}
