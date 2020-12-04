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
