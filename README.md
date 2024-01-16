
# Lomy Chat

This is a simple peer-to-peer chat program that lives inside your Terminal. It's a 
small toy project allowing me to refamiliarise myself with Go.

![screenshot.png](images%2Fscreenshot.png)

### Running it
You can build the binary with
```shell
go build lomy
```

... and then run it with
```shell
./lomy --address=<peerAddress:port>
```


### How its structured

The entrypoint, as you would expect, is [main.go](main.go). It calls out to two 
packages, `connection.go` and `tui.go`.

[connection.go](connection/connection.go) encapsulates the network connection with the 
peer. The API is simple: the returned struct exposes two string channels, one for 
retrieving messages and one for sending them.

[tui.go](tui/tui.go) encapsulates the terminal user interface (written using the 
[tview](https://github.com/rivo/tview) library). Here too the API is simple. The
`WriteToTextView (string)` method writes a line to the chat; and the `InputFieldQueue` 
property provides a channel for receiving user input.


### Missing features

This is never going to be a fully fledged chat application. But here are a few things 
I'd like to add sometime soon.

- [ ] Encryption. Both peers provide a secret key when starting their chat instances, 
 and the application encrypts & decrypts their messages on the fly.
- [ ] Improve exit behaviour. At the moment, when your peer disconnects your program 
  exits without clear messaging.
- [ ] Move the initial user messaging into the TUI. Currently, we log some information 
  to stdout until a Connection is made. That could be shown as a modal within the TUI 
  instead (hovering over the chat UI).