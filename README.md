# Tin Can Communicator
Command line chat app built in Golang with tcp sockets and a custom binary protocol.

Currently a work in progress.

## Features
* Chat with other people!
* Has 255 possible chat rooms to join. (Room 0 is reserved for server-wide messages).
These are numbered 1 through 255.
* chat messages are only broadcast to members of the same chat room.
* Utilizes a TLV style binary protocol to transfer messages.

## Basic Usage
### Server
1. Build the executable using `go build server_or_client`
2. Run the executable with `./server_or_client -mode=server`
6. When you want to stop the server, simply type ctrl+c.

### Client
1. Build the executable using `go build server_or_client`
2. Run the executable with `./server_or_client -mode=client`
3. Enter a username.
4. Enter a desired chat room number (1-255).
5. Chat!
6. When you're finished, simply type ctrl+c to quit.

## To Do
0. Add in unit testing and adopt TDD into the process to help reduce bugs.
1. Clean up code where TLV encoding and decoding happen to make it more obvious whats going on.
2. Make error handling more robust
3. Allow users to easily join a different room or multiple rooms.
4. Use room 0 better for messages originating from the server. (Shutdown message, or general news and tips)
5. Setup auth system to create users that are persistent.
6. Create database for persistent storage of chat logs.
7. Fix echoing of messages on client side when that user enters a piece of chat.
EX: If user#1 types in "Hello Everyone!", then they will see it twice, whereas everyone else in that room sees it once.
8. Go from golang CLI client to nodeJS or Python client. Alternatively, have a VueJS frontend and turn it into a web app. 

### References
These are places that I had read from and that I found very helpful in building this.

* [https://levelup.gitconnected.com/binary-encoding-of-variable-length-options-with-golang-4481ff59e767](https://levelup.gitconnected.com/binary-encoding-of-variable-length-options-with-golang-4481ff59e767) 
* [https://www.thepolyglotdeveloper.com/2017/05/network-sockets-with-the-go-programming-language/](https://www.thepolyglotdeveloper.com/2017/05/network-sockets-with-the-go-programming-language/)
* [https://medium.com/johnshenk77/create-a-simple-chat-application-in-go-using-websocket-d2cb387db836](https://medium.com/@johnshenk77/create-a-simple-chat-application-in-go-using-websocket-d2cb387db836)
* [https://www.cs.dartmouth.edu/campbell/cs60/socketprogramming.html](https://www.cs.dartmouth.edu/~campbell/cs60/socketprogramming.html)
* [https://stackoverflow.com/questions/2681267/what-is-the-fundamental-difference-between-websockets-and-pure-tcp?lq=1](https://stackoverflow.com/questions/2681267/what-is-the-fundamental-difference-between-websockets-and-pure-tcp?lq=1)
* [https://ipfs.io/ipfs/QmfYeDhGH9bZzihBUDEQbCbTc5k5FZKURMUoUvfmc27BwL/architecture/protocol_layers.html](https://ipfs.io/ipfs/QmfYeDhGH9bZzihBUDEQbCbTc5k5FZKURMUoUvfmc27BwL/architecture/protocol_layers.html)
* [https://www.embeddedrelated.com/showthread/comp.arch.embedded/178636-1.php](https://www.embeddedrelated.com/showthread/comp.arch.embedded/178636-1.php)

