package connection

import (
	"lomy/encryption"
)

// An EncryptedConnection implements the Connection interface, but adds
// on-the-fly encryption and decryption atop of `BaseConnection`.
type EncryptedConnection struct {
	baseConnection    BaseConnection
	encryptionSession encryption.EncryptionSession
}

func NewEncryptedConnection(peerAddress string, listeningPort uint, secret string) EncryptedConnection {
	baseConnection := NewConnection(peerAddress, listeningPort)

	encryptionSession := encryption.NewEncryptionSession(secret)

	return EncryptedConnection{
		encryptionSession: encryptionSession,
		baseConnection:    baseConnection,
	}
}

func (ec EncryptedConnection) RetrieveMessage() (string, error) {
	encryptedMessage, _ := ec.baseConnection.RetrieveMessage()
	return ec.encryptionSession.Decrypt(encryptedMessage)
}

func (ec EncryptedConnection) SendMessage(message string) {
	encryptedMessage := ec.encryptionSession.Encrypt(message)
	ec.baseConnection.SendMessage(encryptedMessage)
}

func (ec EncryptedConnection) Close() {
	ec.baseConnection.Close()
}
