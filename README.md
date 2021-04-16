# Tin Can Communicator

Command line chat app built in Golang 1.13 with tcp sockets and a custom binary protocol.

Currently a work in progress.

## Features
* Chat with other people using the command line!
* 256 Chat rooms available. (numbered 0-255)
* 200 character limit per message
* Chat messages are only broadcast to members of the same chat room.
* Utilizes a TLV style binary protocol to encode and decode messages.

## Basic Usage
### Server
1. Build and run the executable using `go run config.go -mode=server -port=####`
2. When you want to stop the server, simply press ctrl+c. (This will disconnect all clients immediately!)

### Client
1. Build and run the executable using `go run config.go -mode=client -address=your.address.goes.here:####`
2. Enter a username (3+ characters).
3. Enter a desired chat room number (0-255).
4. Chat!
5. When you're finished, simply press ctrl+c or the escape key to quit.

## Extra Repos
Here are some other repositories that I created to test out new libraries, or explore some concepts
* https://github.com/Mecope1/serialization_comparison: Compares benchmarks for encoding and decoding the chat app's 
TLV protocol to Json and MessagePack for some sample payloads.

* https://github.com/Mecope1/go-tui-testing: Sandbox where I played around with the Termui library to see how I could 
 use it in the chat app.

## To Do
* Add in additional tests
* Make error handling even more robust
* Allow users to easily join a different room
* Integrate a database to store chat logs
* Setup auth system, allowing persistent users
* Allow users to send requests to my [unit-converter api/ web app](https://convert-this.herokuapp.com/) from the chat app



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

