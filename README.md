# Passgen: a minimal https service for password generation

This project is a minimial web service for password generations. It doesn't
recordy any data about the customer (except for the caller IP address). 
It provides a safe and lightweight tool to generate random passwords of 
variable lenght and format using the package 
[github.com/sethvargo/go-password](https://github.com/sethvargo/go-password).

### Server Usage
To run the service manually:
```
go run main.go
```

To run the service in a container (recommended):
```
docker run -d -v <path_to_certs>:/etc/passgen/certs -p 8443:8443 quay.io/gbsalinetti/passgen:latest
```

To run Passgen on Kubernetes/OpenShift an Helm chart is provided. To install
using the default values:
```
helm install passgen ./helm/passgen
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
To generate the certificates:
```
$ make gencerts
```

To build the image:
```
# make build
```

To tag and push the image to the proper repository (Adjust Makefile to your personal
repository):
```
# make tag && make push
```

To build in OpenShift using the **oc new-app** tool with the Docker strategy:
```
$ oc new-app https://github.com/giannisalinetti/passgen --strategy=docker
```

This will create the following objects:
- imagestream.image.openshift.io "golang"
- imagestream.image.openshift.io "passgen"
- buildconfig.build.openshift.io "passgen"
- deploymentconfig.apps.openshift.io "passgen"
- service "passgen" created

The route won't be created automatically. To create a passthough route use the
command **oc create route passthough** or the manifest provided in the 
*manifests/openshift/route.yaml* after adjusting it to the correct service name.

### Self signed certificate
This project is a proof of concept. Self signed certificate and the associated key 
have been generated with the following command:
```
$ openssl req -newkey rsa:2048 -nodes -keyout server.key -x509 -days 365 -out server.crt
```

The script *hack/genselfsigned.sh* can be manually used to regenerate new certificates.
