package connection

import (
	"bufio"
	"log"
	"net"
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
func NewConnection(targetAddress string) Connection {
	netConn := initiateConnection(targetAddress)
	if netConn == nil {
		// No one was listening on the other end. Let's become listeners ourselves.
		netConn = listenForConnection(targetAddress)
	}

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

func initiateConnection(targetAddress string) net.Conn {
	conn, err := net.Dial("tcp", targetAddress)
	if err != nil {
		log.Println("Connection attempt failed because no Listener was found on the" +
			" other end. Switching to Listening mode.")
		return nil
	}
	return conn
}

func listenForConnection(targetAddress string) net.Conn {
	log.Println("Listening for incoming Connection...")
	ln, err := net.Listen("tcp", targetAddress)
	if err != nil {
		log.Fatalln("Error: ", err)
	}
	conn, err := ln.Accept()
	if err != nil {
		log.Fatalln("Error: ", err)
	}
	log.Println("Connection acquired!")
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
