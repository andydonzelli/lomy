package connection

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

type Connection interface {
	RetrieveMessage() (string, error)
	SendMessage(string)
	Close()
}

// BaseConnection represents an active TCP connection.
// Messages can be sent through this connection with SendMessage(string)`.
// Received messages can be read with `RetrieveMessage() string`.
// The implementation uses buffered channels for both sending and receiving queues to improve performance.
type BaseConnection struct {
	sendQueue chan string
	recvQueue chan string
	netConn   net.Conn
}

// NewConnection Create and initialise a new BaseConnection object.
// Note that this object must be closed with BaseConnection.Close().
//
// A connection attempt is first made against `peerAddress` (if it's not zero valued).
// If that fails listen for incoming connections on `listeningPort` (if it's not zero-valued).
func NewConnection(peerAddress string, listeningPort uint) BaseConnection {
	netConn := connectToPeer(peerAddress, listeningPort)

	sendQueue := make(chan string, 4)
	recvQueue := make(chan string, 4)

	connection := BaseConnection{
		netConn:   netConn,
		sendQueue: sendQueue,
		recvQueue: recvQueue,
	}
	go connection.handleSendQueue()
	go connection.handleRecvQueue()

	return connection
}

func (c BaseConnection) RetrieveMessage() (string, error) {
	message := <-c.recvQueue
	return message, nil
}

func (c BaseConnection) SendMessage(message string) {
	c.sendQueue <- message
}

func (c BaseConnection) Close() {
	c.netConn.Close()
}

func connectToPeer(peerAddress string, listeningPort uint) net.Conn {
	if peerAddress != "" {
		conn := dialPeerOrNil(peerAddress)
		if conn != nil {
			return conn
		}
		log.Println("Connection attempt failed because no Listener was found on the" +
			" other end.")

		if listeningPort != 0 {
			log.Println("Switching to Listening mode.")
		} else {
			os.Exit(1)
		}
	}

	return listenForPeer(listeningPort)
}

func dialPeerOrNil(peerAddress string) net.Conn {
	log.Printf("Dialing '%s'...\n", peerAddress)
	conn, _ := net.Dial("tcp", peerAddress)
	return conn
}

func listenForPeer(listeningPort uint) net.Conn {
	listeningAddress := fmt.Sprintf(":%d", listeningPort)
	log.Printf("Listening for incoming Connection at '%s'...\n", listeningAddress)
	listener, err := net.Listen("tcp", listeningAddress)
	if err != nil {
		log.Fatalln("Error: ", err)
	}
	conn, err := listener.Accept()
	if err != nil {
		log.Fatalln("Error: ", err)
	}
	return conn
}

func (c *BaseConnection) handleRecvQueue() {
	for {
		recvMessage, err := bufio.NewReader(c.netConn).ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		c.recvQueue <- recvMessage
	}
}

func (c *BaseConnection) handleSendQueue() {
	for {
		messageToSend := <-c.sendQueue

		_, err := c.netConn.Write([]byte(messageToSend + "\n"))
		if err != nil {
			log.Fatalln(err)
		}
	}
}
