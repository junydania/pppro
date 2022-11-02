SHELL := /bin/bash -o pipefail
.PHONY: clean lint build start-dynamodb

clean:
	rm -rf ./build

build: test
	CGO_ENABLED=0 go build -ldflags="-w -s" -o $(BUILD_DIR)/$(APP_NAME) cmd/app/main.go

docker.build:
	docker build --platform amd64 -t stack-management-api .

start-dynamodb:
	docker run --rm -d \
		--name dynamodb \
		-p 8000:8000 \
		amazon/dynamodb-local:latest



PPRO_SERVICE=welcome-app
PPRO_REPO=""
BRANCH=$(shell git rev-parse --abbrev-ref HEAD | sed 's/\//-/g') # get branch name and replace forward-slashes

.PHONY:

build-container: ecr-login .PHONY

	# pull (previously) built image as cache
	# build and tag, specify cache source
	time docker build --platform linux/amd64 --file=Dockerfile . \
		--pull \
		--build-arg=XAPO_SERVICE=$(XAPO_SERVICE) \
		--tag=$(XAPO_SERVICE):latest \
		--tag=$(XAPO_REPO):)/$(XAPO_SERVICE):last_build \
		--tag=$(XAPO_REPO)/$(XAPO_SERVICE):branch-$(BRANCH) \
		--tag=$(XAPO_REPO)/$(XAPO_SERVICE):latest \
		--cache-from=$(XAPO_SERVICE):latest \
		--cache-from=$(XAPO_REPO)/$(XAPO_SERVICE):branch-$(BRANCH) \
		--cache-from=$(XAPO_REPO)/$(XAPO_SERVICE):branch-master \
		--cache-from=$(XAPO_REPO)/$(XAPO_SERVICE):last_build

push-container: ecr-login .PHONY
	# push images as cache for future builds
	docker push $(XAPO_REPO)/$(XAPO_SERVICE):last_build
	docker push $(XAPO_REPO)/$(XAPO_SERVICE):branch-$(BRANCH)
	docker push $(XAPO_REPO)/$(XAPO_SERVICE):latest

ecr-login: .PHONY
	aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin $(XAPO_REPO)

	