FROM docker.io/golang

MAINTAINER Gianni Salinetti <gbsalinetti@extraordy.com>

# Copy files for build
COPY server.crt server.key main.go /go/src/passgen-svc/

# Set the working directory
WORKDIR /go/src/passgen-svc

# Download dependencies
RUN go get -d -v ./...

# Install the package
RUN go install -v ./...

EXPOSE 443

CMD ["passgen-svc"]
