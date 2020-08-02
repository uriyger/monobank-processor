
include .env
export $(shell sed 's/=.*//' .env)

DANGLING_IMAGES = $(shell docker images -q -f "dangling=true")

build:
	GOOS=linux go build -o app_linux cmd/main.go
	docker build -t ${REGISTRY_IMAGE}:latest .
# Cleanup dangling images

ifneq ($(DANGLING_IMAGES),)
	docker rmi -f $(DANGLING_IMAGES)
endif

publish:build
	docker push ${REGISTRY_IMAGE}:latest

deploy:publish
	gcloud compute instances update-container ${GCE_INSTANCE} \
	--container-image ${REGISTRY_IMAGE}:latest \
	--container-env-file .env

