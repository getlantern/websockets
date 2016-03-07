This project is a simplistic benchmark of WebSocket performance. It includes a
client and server.

To run server:

```sh
go run server.go -logtostderr
```

To run client with 10,000 connections:

```sh
go run client.go 10000
```
