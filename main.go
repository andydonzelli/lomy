package main

import (
	"flag"
	"lomy/connection"
	"lomy/tui"
)

func main() {
	targetAddressPtr := flag.String("address", "localhost:9000", "Target address")
	flag.Parse()

	conn := connection.NewConnection(*targetAddressPtr)
	defer conn.Close()

	tuiApp := tui.CreateTui()

	go func() {
		for {
			receivedMessage := <-conn.RecvQueue
			tuiApp.WriteToTextView("Them: " + receivedMessage)
		}
	}()

	go func() {
		for {
			messageToSend := <-tuiApp.InputFieldQueue
			conn.SendQueue <- messageToSend
		}
	}()

	tuiApp.RunAppAndBlock()
}
