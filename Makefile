APP_NAME=passgen
REGISTRY=quay.io
NAMESPACE=gbsalinetti

gencerts:
	hack/genselfsigned.sh

build:
	docker build -t $(APP_NAME) .

tag:
	docker tag $(APP_NAME) $(REGISTRY)/$(NAMESPACE)/$(APP_NAME):latest

push:
	docker push $(REGISTRY)/$(NAMESPACE)/$(APP_NAME):latest

