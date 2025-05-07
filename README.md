## Comment
This project uses a simplified `Proof-of-Work (POW)` algorithm inspired by [Hashcash](https://en.wikipedia.org/wiki/Hashcash). 
The server issues a challenge containing a random seed and a difficulty level. 
The client must find a nonce such that the `SHA-256 hash` of `seed + nonce` has a given number of leading
zero bits.

This approach helps mitigate denial-of-service (DDoS) attacks by forcing clients to spend computational effort 
before the server responds. 
It’s lightweight, easy to implement, and does not rely on external state or synchronization, 
making it suitable for TCP-based challenge-response systems.

### First run
```
1. cp .env.dist .env

# docker-compose up
make up

# docker-compose down
make down
```

