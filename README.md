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

#### Response Samples

A typical `HandshakeV10` response:
```
nc localhost 3306 | xxd
00000000: 4a00 0000 0a38 2e30 2e32 3200 0900 0000  J....8.0.22.....
00000010: 7b76 1b3e 7269 5638 00ff ffff 0200 ffc7  {v.>riV8........
00000020: 1500 0000 0000 0000 0000 0001 1e61 4863  .............aHc
00000030: 5643 2a3e 5b1c 7e00 6361 6368 696e 675f  VC*>[.~.caching_
00000040: 7368 6132 5f70 6173 7377 6f72 6400       sha2_password.
```

For `db3`, the `MYSQL_ROOT_HOST` configuration produces an quasi-unexpected response:
```
nc localhost 3308 | xxd
00000000: 4300 0000 ff6a 0448 6f73 7420 2731 3732  C....j.Host '172
00000010: 2e32 302e 302e 3127 2069 7320 6e6f 7420  .20.0.1' is not 
00000020: 616c 6c6f 7765 6420 746f 2063 6f6e 6e65  allowed to conne
00000030: 6374 2074 6f20 7468 6973 204d 7953 514c  ct to this MySQL
00000040: 2073 6572 7665 72                         server
```

Lastly, `server1` responds with random data:
```
nc localhost 3311 | xxd
00000000: 2d65 2097 b09d 14cf 806b da0c bde0 828b  -e ......k......
00000010: 5613 2a9d 8b80 f6b4 53fa 9c2f b0b0 cf0b  V.*.....S../....
00000020: 7ba3 5ff2 2cf6 f433 9c08 7f4d 859d 164f  {._.,..3...M...O
00000030: 3b17 17f4 56dd 71f3 530e cde2 b62e 30b1  ;...V.q.S.....0.
00000040: f76c acca 0d71 d6ea 0135 a853 1264 2559  .l...q...5.S.d%Y
// omitted for clarity 
```

#### Protocol

The `HandshakeV10` response provides the following information:

| Field  | Size  | Description|
|---|---|---|
| protocol_version | 1 | 0x0a protocol_version|
| server_version   | (string.NUL) | human-readable server version|
| connection_id | 4 | connection id |
| auth_plugin_data_part_1| (string.fix_len) [len=8] | first 8 bytes of the auth-plugin data
| filler_1 | 1 | 0x00 |
| capability_flag_1 | 2 | lower 2 bytes of the Protocol::CapabilityFlags (optional) |
| character_set | 1 | default server character-set, only the lower 8-bits Protocol::CharacterSet (optional) |
| status_flags |  2 | Protocol::StatusFlags (optional) | 
| capability_flags_2 | 2 | upper 2 bytes of the Protocol::CapabilityFlags |
| auth_plugin_data_len | 1 | length of the combined auth_plugin_data, if auth_plugin_data_len is > 0|
| auth_plugin_name |  (string.NUL) | ame of the auth_method that the auth_plugin_data belongs to |

Per, MySQL docs, Bug#59453 the auth-plugin-name is missing the terminating NUL-char in versions prior to 5.5.10 and 5.6.2. 

#### Scanner Output
The scanner app displays the MySQL banner information as follows:
```
localhost:3306
ProtocolVersion: 10
ServerVersion: 8.0.22
ConnectionId: 10
AuthPluginName: caching_sha2_password
StatusFlags: 2
```

#### Interesting Observation
Whilst performing a few quick sanity checks against MySQL assets discovered on the internet, the `Connection ID` field caught my eye. Specifically, at least during testing, I noticed freshly started MySQL instances would have a low `Connection ID` whereas some systems I stumbled across yielded a sizable number.

The MySQL documentation states:
> Returns the connection ID (thread ID) for the connection. Every connection has an ID that is unique among the set of currently connected clients. The value returned by CONNECTION_ID() is the same type of value as displayed in the ID column of the INFORMATION_SCHEMA.PROCESSLIST table, the Id column of SHOW PROCESSLIST output, and the PROCESSLIST_ID column of the Performance Schema threads table.`

Assuming this holds true not just for TCP connections but UNIX sockets and windows named pipes as well; I briefly postulated whether or not this could have value beyond basic utilization insights. 

For example, given identical MySQL versions but sizeably different connection id information:


| System  | Connection ID  |
|---|---|
| system1 mysql:5.7 | 3244 |
| system2 mysql:5.7 | 105001700 |

One may infer (probabilistically):
* mysql service recently started on system1 (uptime)
* system1 utilization is lower than system2
* given a relatively new 0-day/CVE for which remediation requires a service restart, system1 might be freshly patched, whereas system2 might still be vulnerable



### References
* [MySQL Docs - Protocol::Handshake](https://dev.mysql.com/doc/internals/en/connection-phase-packets.html#packet-Protocol::Handshake)
* [MySQL Docs - Connection ID](https://dev.mysql.com/doc/refman/5.7/en/information-functions.html#function_connection-id)
* [Writing MySQL Proxy in GO for self-learning: Part 2 — decoding handshake packet](https://medium.com/@alexanderravikovich/writing-mysql-proxy-in-go-for-learning-purposes-part-2-decoding-connection-phase-server-response-7091d87e877e)