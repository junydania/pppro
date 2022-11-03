SHELL := /bin/bash -o pipefail
PPRO_SERVICE=pppro-hello
PPRO_REPO=""
BRANCH=$(shell git rev-parse --abbrev-ref HEAD | sed 's/\//-/g') # get branch name and replace forward-slashes

.PHONY: clean lint build start-dynamodb

clean:
	rm -rf ./build

local-build: test
	CGO_ENABLED=0 go build -ldflags="-w -s" -o $(BUILD_DIR)/$(APP_NAME) cmd/app/main.go

docker.build:
	docker build --platform amd64 -t helloapp .

start-dynamodb:
	docker run --rm -d \
		--name dynamodb \
		-p 8000:8000 \
		amazon/dynamodb-local:latest

build-container: ecr-login .PHONY

	# pull (previously) built image as cache
	# build and tag, specify cache source
	time docker build --platform linux/amd64 --file=Dockerfile . \
		--pull \
		--build-arg=pppro_SERVICE=$(PPPRO_SERVICE) \
		--tag=$(PPPRO_SERVICE):latest \
		--tag=$(PPPRO_REPO):)/$(PPPRO_SERVICE):last_build \
		--tag=$(PPPRO_REPO)/$(PPPRO_SERVICE):branch-$(BRANCH) \
		--tag=$(PPPRO_REPO)/$(PPPRO_SERVICE):latest \
		--cache-from=$(PPPRO_SERVICE):latest \
		--cache-from=$(PPPRO_REPO)/$(PPPRO_SERVICE):branch-$(BRANCH) \
		--cache-from=$(PPPRO_REPO)/$(PPPRO_SERVICE):branch-master \
		--cache-from=$(PPPRO_REPO)/$(PPPRO_SERVICE):last_build

push-container: ecr-login .PHONY
	# push images as cache for future builds
	docker push $(PPPRO_REPO)/$(PPPRO_SERVICE):last_build
	docker push $(PPPRO_REPO)/$(PPPRO_SERVICE):branch-$(BRANCH)
	docker push $(PPPRO_REPO)/$(PPPRO_SERVICE):latest

ecr-login: .PHONY
	aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin $(PPPRO_REPO)
