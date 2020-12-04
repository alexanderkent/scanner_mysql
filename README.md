### Intro
This repo is part of a technical assessment to build a scanner to detect MySQL running on a port on a host.

### Assessment
Using Go, write a scanner to detect MySQL running on a port on a host. It should connect to a single port that is running MySQL, and output some information about the MySQL instance’s configuration. The scanner should detect as much as it can from a single MySQL handshake, without logging in.

### Getting started
This project has a docker-compose.yml file, which will start the test chassis.