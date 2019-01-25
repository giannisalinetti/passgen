#!/bin/bash

SERVER_KEY=../server.key
SERVER_CRT=../server.crt
DURATION=365

openssl req -newkey rsa:2048 -nodes -keyout $SERVER_KEY -x509 -days $DURATION -out $SERVER_CRT
if [ $? != 0 ]; then
    exit 1
fi

exit 0

