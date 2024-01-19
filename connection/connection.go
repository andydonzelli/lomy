package connection

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

// Connection represents an active TCP connection.
// Users can send messages by writing to buffered channel Connection.SendQueue,
// and can receive messages by writing to buffered channel Connection.RecvQueue.
type Connection struct {
	SendQueue chan<- string
	RecvQueue <-chan string

	internalSendQueue chan string
	internalRecvQueue chan string
	netConn           net.Conn
}

// NewConnection Create and initialise a new Connection object.
// Note that this object must be closed with Connection.Close().
//
// A connection attempt is first made against `peerAddress` (if it's not zero valued).
// If that fails listen for incoming connections on `listeningPort` (if it's not zero-valued).
func NewConnection(peerAddress string, listeningPort uint) Connection {
	netConn := connectToPeer(peerAddress, listeningPort)

	internalSendQueue := make(chan string, 4)
	internalRecvQueue := make(chan string, 4)

	connection := Connection{
		netConn:           netConn,
		internalSendQueue: internalSendQueue,
		internalRecvQueue: internalRecvQueue,
		SendQueue:         internalSendQueue,
		RecvQueue:         internalRecvQueue,
	}
	go connection.handleSendQueue()
	go connection.handleRecvQueue()

	return connection
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

func (c *Connection) handleRecvQueue() {
	for {
		recvMessage, err := bufio.NewReader(c.netConn).ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		c.internalRecvQueue <- recvMessage
	}
}

func (c *Connection) handleSendQueue() {
	for {
		messageToSend := <-c.internalSendQueue

		_, err := c.netConn.Write([]byte(messageToSend + "\n"))
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func (c *Connection) Close() {
	c.netConn.Close()
}
