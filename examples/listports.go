package main

import (
	"fmt"
	"github.com/kezl/lsport"
)

func main() {

	fmt.Println("List ports:")
	ports, err := lsport.PortsSlice()
	if err == nil {
		for _, port := range ports {
			fmt.Printf("%s\n", port)
		}
	} else {
		fmt.Printf("%s\n", err.Error())
	}
}
