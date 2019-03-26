# G(RPC) Chat

A very simple chat implementation in Go using a single bi-directional gRPC stream.

## Getting started

Clone the repository:

```bash
git clone https://www.github.com/sebnyberg/gchat.git

cd gchat
```

Run the server

```bash
go run main.go server start
```

Connect as a client and start chatting. You can join as many users in different terminals:

```bash
go run main.go client connect --username bob
```

Another option is to install the service and run the binary directly:

```bash
go install

# $GOPATH/bin has to be in your path for this to work
gchat --help
```

## Structure

The app uses [Cobra](https://github.com/spf13/cobra) to provide CLI commands and flags - in this case simply the `--username` flag.

The gRPC service is defined in `pkg/pb/gchat.proto`. To generate code for the service you will need `protoc` installed:

```bash
protoc pkg/pb/gchat.proto --go_out=plugins=grpc:.
```

The service consists of a single bi-directional stream that expects a stream of `ChatSessionRequest` to go from the client to the server, and a stream of `ChatSessionResponse` to be sent back from the server to the client. Each request and response contains a message with an author and content.

Logic for handling client-server connections is in `pkg/client/chat.go` and `pkg/server/chat.go`. When a client connects to the server (server-side), it will receive a shared chat channel. This channel is where all messages from clients to the server end up. Then, the client subscribes to the broadcast (`pkg/server/broadcast.go`), which will create a client-specific channel that receives messages whenever the shared chat channel does. Each of these client-specific, listening channels, are used to relay messages back to its respective clients.

There's a lot that can be improved from this example, it simply serves as an example of how easy it is to set up a chat with gRPC and Go.