# Go Blockchain 
A proof of concept of Blockchain system following the series of [MyCoralHealth](https://mycoralhealth.medium.com). This blockchain is meant for a healthcare company to store Pulse rate info of users.

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
  "Index": 2,
  "Timestamp": "2021-10-19 22:47:16.877277 +0700 +07 m=+19.383792232",
  "BPM": 10,
  "Hash": "0c4f602b3d2dec1918949fcccea78934f58761f5cbee4fd46ecf7985c50124f2",
  "PrevHash": "0c4eb43a2ea0eda7529de470d2ba6bb0f30882700c1307dbf4600986c6befc9e",
  "Difficulty": 1,
  "Nonce": "f"
}
```

- Get current blockchain info
```
curl -X GET http://localhost:8080

Response:
[
 {
  "Index": 0,
  "Timestamp": "2021-10-19 22:46:57.49489 +0700 +07 m=+0.001598502",
  "BPM": 0,
  "Hash": "f1534392279bddbf9d43dde8701cb5be14b82f76ec6607bf8d6ad557f60f304e",
  "PrevHash": "",
  "Difficulty": 1,
  "Nonce": ""
 },
 {
  "Index": 1,
  "Timestamp": "2021-10-19 22:47:01.342396 +0700 +07 m=+3.849066168",
  "BPM": 10,
  "Hash": "0c4eb43a2ea0eda7529de470d2ba6bb0f30882700c1307dbf4600986c6befc9e",
  "PrevHash": "f1534392279bddbf9d43dde8701cb5be14b82f76ec6607bf8d6ad557f60f304e",
  "Difficulty": 1,
  "Nonce": "2"
 },
 {
  "Index": 2,
  "Timestamp": "2021-10-19 22:47:16.877277 +0700 +07 m=+19.383792232",
  "BPM": 10,
  "Hash": "0c4f602b3d2dec1918949fcccea78934f58761f5cbee4fd46ecf7985c50124f2",
  "PrevHash": "0c4eb43a2ea0eda7529de470d2ba6bb0f30882700c1307dbf4600986c6befc9e",
  "Difficulty": 1,
  "Nonce": "f"
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

## Proof Of Work
Solve a hard math problem: find the Nonce so that the hash of the next block has the number of leading Os equals to Difficulty

## Proof Of Stake
Instead of nodes competing with each other to solve hashes, in Proof of Stake, blocks are “minted” or “forged” based on the amount of tokens each node is willing to put up as collateral. These nodes are called validators

## Notes
- mining bitcoin = solve a hard math problem
- SHA-256: cryptographic hash, idempotency
- Netcat (or nc ) is a command-line utility that reads and writes data across network connections, using the TCP or UDP protocols

## Packages
- github.com/joho/godotenv: pretty print structs and slices
- github.com/gorilla/mux: write web handlers
- github.com/joho/godotenv: read environment variables from a .env file
