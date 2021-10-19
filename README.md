# Proof-of-Work Blockchain
A proof of concept of Blockchain system using Proof-of-Work consensus algorithm

## What we'll do
- Create your own blockchain
- Understand how hashing works in maintaining integrity of the blockchain
- See how new blocks get added
- See how tiebreakers get resolved when multiple nodes generate blocks
- View your blockchain in a web browser
- Write new blocks
- Get a foundational understanding of the blockchain so you can decide where your journey takes you from here!

## Getting started
- Start the app
```
go mod download
go run *.go
```

- Create a new block
```
curl -X POST -d '{"BPM":10}' http://localhost:8080

Response:
{
  "Index": 4,
  "Timestamp": "2021-10-19 08:00:58.741514 +0700 +07 m=+743.922474724",
  "BPM": 10,
  "Hash": "ce27d4811945d5495521eba45e8d1958828d9be453918d37a78c0b477ebf236e",
  "PrevHash": "a0475ea7b1ef6b5ab5c980c9299341525192babf5342c2eac729c2d25315c391"
}
```

- Get current blockchain info
```
curl -X GET http://localhost:8080

Response:
[
 {
  "Index": 0,
  "Timestamp": "2021-10-19 07:42:00.967173 +0700 +07 m=+0.001552370",
  "BPM": 0,
  "Hash": "",
  "PrevHash": ""
 },
 {
  "Index": 1,
  "Timestamp": "2021-10-19 07:56:54.090725 +0700 +07 m=+499.274069703",
  "BPM": 10,
  "Hash": "71e7cdbf2f126b40f4a1777a8a824271a32607a406ca832a98ed775317ad481e",
  "PrevHash": ""
 }
]
```

## Networking
- Start a TCP server at port 9000
```
go run *.go 
```
- Start a client to connect to the server using TCP
```
nc localhost 9000

Output:
Enter a new BPM: [enter a number]
```

## Notes
- mining bitcoin = solve a hard math problem
- SHA-256: cryptographic hash
- Netcat (or nc ) is a command-line utility that reads and writes data across network connections, using the TCP or UDP protocols

## Packages
- github.com/joho/godotenv: pretty print structs and slices
- github.com/gorilla/mux: write web handlers
- github.com/joho/godotenv: read environment variables from a .env file
