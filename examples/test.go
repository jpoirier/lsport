/*
Tested using an arduino with a simple loopback program,
see loopback.ino

*/

package main

import (
	"fmt"
	"github.com/kezl/lsport"
	"time"
)

func main() {

	// rxBuf size should not exceed the capacity of int32
	rxBuf := make([]byte, 256)

	s := lsport.Conf{}
	lsport.Init(&s, "/dev/ttyUSB0")
	lsport.SetParams(&s, 115200, 8, 1)

	go PollingRead(s.Port, rxBuf, 100, 100)
	// Allow some settling time, tweak as required.
	time.Sleep(time.Duration(2000) * time.Millisecond)

	lsport.Write(s.Port, []byte("Hello "))
	time.Sleep(time.Duration(200) * time.Millisecond)

	lsport.Write(s.Port, []byte("World "))
	time.Sleep(time.Duration(200) * time.Millisecond)

	lsport.Close(&s)

	fmt.Println("\nList ports:")
	ports, err := lsport.PortsSlice()
	if err == nil {
		for _, port := range ports {
			fmt.Printf("%s\n", port)
		}
	} else {
		fmt.Printf("%s\n", err.Error())
	}
}

func PollingRead(port lsport.Port, rxBuf []byte, sleep time.Duration, timeout uint) {

	for {
		waiting, err := lsport.Waiting(port)
		if waiting > 0 && err == nil {
			lsport.Read(port, rxBuf)
			fmt.Println(string(rxBuf))

		}
		time.Sleep(time.Millisecond * sleep)
	}
}
