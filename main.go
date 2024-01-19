package main

import (
	"flag"
	"fmt"
	"lomy/connection"
	"lomy/tui"
	"os"
)

func main() {
	peerAddress, listeningPort := parseCliArgs()

	conn := connection.NewConnection(peerAddress, listeningPort)
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

func parseCliArgs() (string, uint) {
	helpPtr := flag.Bool("help", false, "Print usage information")
	peerAddressPtr := flag.String("peerAddress", "", "Peer's address and port. E.g., "+
		"'192.168.0.18:9000'")
	listeningPortPtr := flag.Uint("listeningPort", 0,
		"Port on which to listen for incoming connections. E.g., '9000'")

	flag.Usage = func() {
		w := flag.CommandLine.Output()
		fmt.Fprintf(w, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintf(w,
			"\nAt least one of -peerAddress and -listeningPort must be set."+
				" If both are provided an attempt is made to dial peerAddress;"+
				" if that fails the program begins listening on listeningPort.\n")
	}

	flag.Parse()

	if *helpPtr {
		flag.Usage()
		os.Exit(0)
	}

	// If both arguments were not specified print Usage and exit
	if *peerAddressPtr == "" && *listeningPortPtr == 0 {
		flag.Usage()
		os.Exit(1)
	}

	return *peerAddressPtr, *listeningPortPtr
}
