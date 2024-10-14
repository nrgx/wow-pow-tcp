# "Word of Wisdom" test task

A test task with the following requirements:

- Design and implement “Word of Wisdom” TCP server.
- TCP server should be protected from DDOS attacks with the Proof of Work (https://en.wikipedia.org/wiki/Proof_of_work), the challenge-response protocol should be used.
- The choice of the POW algorithm should be explained.
- After Proof Of Work verification, server should send one of the quotes from “word of wisdom” book or any other collection of the quotes.
- Docker file should be provided both for the server and for the client that solves the POW challenge

## Chosen PoW algorithm and justification of choice

`Hashcash` PoW system was chosen as PoW algorithm for its simplicity in implementation and resistant from DDoS attacks just by its nature.
Hashcash representation:

```
1:4:1728913453:172.21.0.3:35192::MTg4ODk1:146572

where

1                   - version of Hashcash 
4                   - number of leading zeroes after hashing  
1728913453          - date of creation  
172.21.0.3:35192    - resource (e.g. user who requested)
MTg4ODk1            - random number encoded in base64 
146572              - number of iteration to run to get hash with 4 leading zeroes of new hash (e.g. nonce)
```

This system is already used in many other real-world applications (core of Bitcoin blockchain or anti-spam systems). 
DDoS resistance: Hashcash protects servers from DDoS attacks by requiring users to solve a challenge, which forces users to run calculations on their side before receiving reward. For attackers it'll be costly operations with little to no profit. 

Users who want to get a reword must run some computations with brute force. The initial value of Hashcash is hashed via hashing algorithm (SHA-1,SHA-256,SHA-512,etc.). It's a challenge for user to solve. The second parameter of Hashcash defines number of leading zeroes in hash which is the solution for the challenge (e.g. hash must start with 0000). User runs a for loop hashing initial value by adding nonce so that `hashfunc(challenge+nonce)` can produce a hash with the required number of zeroes.

## Usage

```
docker-compose up
```

or via Makefile

```
make

Makefile commands:
help            Show available commands
server          Build and run the server using Docker
client          Build and run the client using Docker
up              Run both client and server using docker-compose
down            Stop both client and server using docker-compose
deps            Download Go dependencies
lint            Run the linter on the project

Usage examples:
make server      - to run the server only
make client      - to run the client only
make up          - to run both server and client with docker-compose
```
