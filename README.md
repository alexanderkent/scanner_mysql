### Intro
This repo is part of a technical assessment to build a scanner to detect MySQL running on a port on a host.

### Assessment
Using Go, write a scanner to detect MySQL running on a port on a host. It should connect to a single port that is running MySQL, and output some information about the MySQL instance’s configuration. The scanner should detect as much as it can from a single MySQL handshake, without logging in.

### Getting started
This project has a docker-compose.yml file, which will start the test chassis.

```
docker-compose build && docker-compose up
```

The test chassis hosts various versions of MySQL and MariaDB servers. In addition, a very simple server/fuzzer creates random responses to help test scanner robustness.

| Port  | Description  |
|---|---|
| 3306 | mysql:latest |
| 3307 | mysql:5.7 |
| 3308 | mysql:5.6 |
| 3309 | mariadb:latest |
| 3310 | mariadb:10.1 |
| 3311 | server/fuzzer |

### Compiling 
For convenience a Makefile has been provided.
```
make build
```

### Running the scanner

To run the scanner against the test chassis

```
./bin/scanner
```

To run the scanner against a single target
```
./bin/scanner host port
```


### Methodology
By way of background, a MySQL server responds to a client connection with a handhake packet. Depending on the server version and configuration options different variants of the initial packet are sent. 

| MySQL Version  | Handshake  |
|---|---|
| >= 3.21.0 | HandshakeV10 |
|  < 3.21.0  | HandshakeV9 | 

Having received a valid handshake packet, a MySQL client will reply with either a `SSL Connection Request Packet` and then a `Handshake Response Packet` or merely a `Handshake Response Packet` when SSL is not used. This creates an additional opportunity to gather asset information such as server capabilities and supported cipher suites. However, for the purpose of this assessment, no data is ever sent to the server beyond the initial TCP connection.