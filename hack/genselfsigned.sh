# Self-signed certificates generation script
# Disclaimer: THIS TOOL IS NOT SUITABLE FOR PRODUCTION USE
#!/bin/bash

set -o errexit                                                                  
set -o nounset                                                                  
set -o pipefail                                                                 
                                                                                
cleanup() {
    printf "Cleaning up old certificates and keys.\n\n"
    printf "##########################################\n"
    rm -rf $CERTS_DIR/*
}

CERTS_DIR=$(dirname "${BASH_SOURCE}")/../certs                                   

SERVER_KEY=$CERTS_DIR/server.key
SERVER_CRT=$CERTS_DIR/server.crt
DURATION=365
CIPHER=rsa
BITS=2048

if [ -f $SERVER_KEY ] || [ -f $SERVER_CRT ]; then
    cleanup
fi
    
openssl req -newkey ${CIPHER}:${BITS} \
        -nodes -keyout ${SERVER_KEY} -x509 \
        -days ${DURATION} -out ${SERVER_CRT}
if [ $? != 0 ]; then
    exit 1
fi

exit 0

