package main

import (
	"fmt"
	"github.com/bugst/go-serial"
	"github.com/davecgh/go-spew/spew"
	"log"
)

func main() {
	scanPorts()
	read()
}

func read() {
	mode := &serial.Mode{
		BaudRate: 9600,
		DataBits: 8,
		StopBits: serial.TwoStopBits,
		Parity:   serial.NoParity,
	}
	//mode := &serial.Mode{
	//	BaudRate: 57600,
	//	Parity: serial.EvenParity,
	//	DataBits: 7,
	//	StopBits: serial.OneStopBit,
	//}

	port, err := serial.Open("COM3", mode)
	if err != nil {
		log.Fatal(err)
	}

	buff := make([]byte, 16)
	i := 0
	for {
		n, err := port.Read(buff)
		if err != nil {
			log.Fatal(err)
			break
		}
		if n == 0 {
			fmt.Println("\nEOF")
			break
		}
		//fmt.Printf("%s\n", string(buff[:n]))
		spew.Dump(buff[:n])

		//if i > 100 {
		//	break
		//}
		i++

	}
}

func scanPorts() {
	ports, err := serial.GetPortsList()
	if err != nil {
		log.Fatal(err)
	}
	if len(ports) == 0 {
		log.Fatal("No serial ports found!")
	}
	for _, port := range ports {
		fmt.Printf("Found port: %v\n", port)
	}
}
