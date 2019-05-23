package main

import (
	"fmt"
	"syscall"

	"./handler"
	"github.com/hyperledger/sawtooth-sdk-go/processor"
)

const valUrl = "tcp://localhost:4004"

func main() {
	processor := processor.NewTransactionProcessor(valUrl)
	processor.AddHandler(handler.CreateHandler("5b7349"))
	processor.ShutdownOnSignal(syscall.SIGINT, syscall.SIGTERM)
	err := processor.Start()
	if err != nil {
		fmt.Println("Connction error: ", err)
	}
}
