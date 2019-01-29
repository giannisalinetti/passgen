FROM docker.io/golang

MAINTAINER Gianni Salinetti <gbsalinetti@extraordy.com>

# Define a volume for certificates
VOLUME /etc/passgen/certs/

# Copy files for build
COPY certs/server.crt certs/server.key /etc/passgen/certs/
COPY main.go /go/src/passgen/

# Set the working directory
WORKDIR /go/src/passgen

# Download dependencies
RUN go get -d -v ./...

# Install the package
RUN go install -v ./...

EXPOSE 8443

CMD ["passgen"]
