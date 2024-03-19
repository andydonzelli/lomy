# Lomy Chat

This is a simple peer-to-peer chat program that lives inside your Terminal. It's a small toy
project to refamiliarise myself with Go.

![screenshot.png](images%2Fscreenshot.png)

### Running it

You can build the binary with

```shell
go build lomy
```

... and then run it with

```shell
./lomy [-peerAddress=<peer hostname/ip:port>] [-listeningPort=<port>] [-encryptionSecret=<secret>]
```

At least one of `peerAddress` or `listeningPort` must be provided. If `-peerAddress` is set the
program will attempt connecting to the target. If `-listeningPort` is set the program will sit
listening for a connection on that port. If both are set one attempt will be made to reach
`peerAddress` and if that fails the program will listen on `listeningPort` instead.

Your peer must be waiting in listening mode for the `peerAddress` connection attempt to succeed.

If both parties provide the `encryptionSecret` argument then communication will be encrypted with
their shared secret.

### How it's structured

The entrypoint, as you would expect, is [main.go](main.go). This starts two goroutine loops: one
for monitoring user input and sending it to the peer, the other for monitoring messages from the
peer and writing it to the screen. These call out to three packages: `connection` and `tui`.

[tui.go](tui/tui.go) encapsulates the terminal user interface (written using the
[tview](https://github.com/rivo/tview) library). Here too the API is simple:
`WriteToTextView(string)` writes a line to the chat; `ReadInputLine() string` reads user input.

[connection.go](connection/connection.go) encapsulates the network connection with the peer. The
API is simple: `SendMessage(string)` sends a message through the connection;
`RetrieveMessage() string` returns any recieved messages. Both inbound and outbound messages are
held in buffered channels to improve performance.

[encrypted-connection.go](connection/encrypted-connection.go) provides an implementation of the
Connection interface, but with on-the-fly encryption applied to outbound messages and decryption applied to inbound messages.

[encryption.go](encryption/encryption.go) encapsulates the message encryption operations used in the
EncryptedConnection implementation.

### Missing features

This is never going to be a fully fledged chat application. But here are a few things I'd like to
add sometime soon.

- [ ] Improve exit behaviour. At the moment, when your peer disconnects your program exits without
      clear messaging.
- [ ] Move the initial user messaging into the TUI. Currently, we log some information to stdout
      until a Connection is made. That could be shown as a modal within the TUI instead (hovering
      over the chat UI).
