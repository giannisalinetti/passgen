# Builder image
FROM docker.io/golang AS builder

# Copy files for build
COPY go.mod /go/src/passgen/
COPY main.go /go/src/passgen/

# Set the working directory
WORKDIR /go/src/passgen

# Download dependencies
RUN go get -d -v ./...

# Install the package
RUN go build -v ./...

# Runtime image
FROM registry.access.redhat.com/ubi8/ubi-minimal:latest as bin
COPY --from=builder /go/src/passgen/ /

# Define a volume for certificates
VOLUME /etc/passgen/certs/

EXPOSE 8443

CMD ["/passgen"]
