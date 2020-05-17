# Tin Can Communicator

Command line chat app built in Golang 1.13 with tcp sockets and a custom binary protocol.

Currently a work in progress.

## Features
* Chat with other people!
* Has 255 possible chat rooms to join. (Room 0 is reserved for server-wide messages).
These are numbered 1 through 255.
* chat messages are only broadcast to members of the same chat room.
* Utilizes a TLV style binary protocol to transfer messages.

## Basic Usage
### Server
1. Build the executable using `go run config.go -mode=server -port=####`
2. When you want to stop the server, simply press ctrl+c.

### Client
1. Build the executable using `go run config.go -mode=client -address=your.address.goes.here:####`
2. Enter a username.
3. Enter a desired chat room number (1-255).
4. Chat!
5. When you're finished, simply press ctrl+c to quit.

## To Do
* Add in additional tests
* Make error handling more robust
* Allow users to easily join a different room or multiple rooms.
* Use room 0 better for messages originating from the server. (Shutdown message, or general news and tips)
* Setup auth system to create users that are persistent.
* Create database for persistent storage of chat logs. 

### References
These are places that I had read from and that I found very helpful in building this.
* [gobyexample.com](gobyexample.com)
* [golangbot.com](golangbot.com)
* [https://searchnetworking.techtarget.com/tutorial/Protocols-Lesson-2-Binary-and-the-Internet-Protocol](https://searchnetworking.techtarget.com/tutorial/Protocols-Lesson-2-Binary-and-the-Internet-Protocol)
* [https://levelup.gitconnected.com/binary-encoding-of-variable-length-options-with-golang-4481ff59e767](https://levelup.gitconnected.com/binary-encoding-of-variable-length-options-with-golang-4481ff59e767) 
* [https://www.thepolyglotdeveloper.com/2017/05/network-sockets-with-the-go-programming-language/](https://www.thepolyglotdeveloper.com/2017/05/network-sockets-with-the-go-programming-language/)
* [https://medium.com/johnshenk77/create-a-simple-chat-application-in-go-using-websocket-d2cb387db836](https://medium.com/@johnshenk77/create-a-simple-chat-application-in-go-using-websocket-d2cb387db836)
* [https://www.cs.dartmouth.edu/campbell/cs60/socketprogramming.html](https://www.cs.dartmouth.edu/~campbell/cs60/socketprogramming.html)
* [https://stackoverflow.com/questions/2681267/what-is-the-fundamental-difference-between-websockets-and-pure-tcp?lq=1](https://stackoverflow.com/questions/2681267/what-is-the-fundamental-difference-between-websockets-and-pure-tcp?lq=1)
* [https://ipfs.io/ipfs/QmfYeDhGH9bZzihBUDEQbCbTc5k5FZKURMUoUvfmc27BwL/architecture/protocol_layers.html](https://ipfs.io/ipfs/QmfYeDhGH9bZzihBUDEQbCbTc5k5FZKURMUoUvfmc27BwL/architecture/protocol_layers.html)
* [https://www.embeddedrelated.com/showthread/comp.arch.embedded/178636-1.php](https://www.embeddedrelated.com/showthread/comp.arch.embedded/178636-1.php)

