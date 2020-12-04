#!/bin/sh

#
# Simple server accepts TCP 3306 and sends 256 bytes of random data
#
while true ; do nc -l -p 3306 -c 'echo -e `dd if=/dev/urandom count=256` '; done