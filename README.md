# Passgen-svc: a minimal https service for password generation

### Server Usage
To run the service manually:
```
go run main.go
```

To run the service in a container (recommended):
```
docker run -d -p 8443:8443 quay.io/gbsalinetti/passgen-svc:latest
```

### Client Usage
To get a standard default password of length 32:
```
$ curl -k https://localhost:8443/passwd
```

To print a password of custom length:
```
$ curl -k 'https://localhost:8443/passwd?length=64'
```

To iterate more times and print a custom number passwords:
```
$ curl -k https://localhost:8443/passwd?iterations=5
```

Custom request with 64 runes, 8 digits, 4 symbols, allowed uppercase, allowed
repetitions an 10 iterations:
```
$ curl -k 'https://localhost:8443/passwd?length=64&digits=8&symbols=4&noupper=false&allowrepeat=true&iterations=10'
```

### Build
To build the image:
```
make build
```

To tag and push the image to the proper repository (Adjust Makefile to your personal
repository):
```
make tag && make push
```

### Self signed certificate
This project is a proof of concept. Self signed certificate and the associated key 
have been generated with the following command:
```
openssl req -newkey rsa:2048 -nodes -keyout server.key -x509 -days 365 -out server.crt
```

The script *hack/genselfsigned.sh* can be used to regenerate new certificates.
