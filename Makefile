VERSION=$(shell date +"%Y%m%d%H%M%S")
APP_GROUP=foundation
APP_NAME=auth
IMAGE_REPO=ccr.ccs.tencentyun.com/fflow/fflow

auth:
	make build APP_NAME=auth APP_GROUP=foundation

faas:
	make build APP_NAME=faas APP_GROUP=foundation

timer:
	make build APP_NAME=timer APP_GROUP=foundation

engine:
	make build APP_NAME=engine APP_GROUP=workflow-app

build:
	rm -f $(APP_NAME)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o $(APP_NAME) service/cmd/$(APP_GROUP)/$(APP_NAME)/main.go
	docker build --platform linux/amd64 -t $(APP_NAME):$(VERSION) -f deployer/dockerfile/$(APP_GROUP)/$(APP_NAME).Dockerfile .
	docker tag $(APP_NAME):$(VERSION) $(IMAGE_REPO):$(APP_NAME)-$(VERSION)
	docker push $(IMAGE_REPO):$(APP_NAME)-$(VERSION)
	rm -f $(APP_NAME)

