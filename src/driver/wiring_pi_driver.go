package driver

// #cgo LDFLAGS: -lwiringPi
// #include <wifingPi.h>
/*
void open(void){
        wiringPiSetup () ;
        pinMode (0, OUTPUT) ;
        digitalWrite (0, HIGH);
        delay (5000);
        digitalWrite (0, LOW);
}
*/
import "C"

func OpenDoor() {
	C.open()
}
