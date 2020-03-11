package main

import (
	"github.com/devplayg/hippo"
	"github.com/devplayg/scale"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	config := &hippo.Config{
		Debug: true,
	}

	//server:=scale.NewServer("ST, .,+0000721kg\r\n")
	//server:=scale.NewServer("ST, .,+0000721kg\r\n")
	server := scale.NewServer("ABCDEFGHIJKLMNOP\r\n")
	engine := hippo.NewEngine(server, config)
	if err := engine.Start(); err != nil {
		panic(err)
	}
}
