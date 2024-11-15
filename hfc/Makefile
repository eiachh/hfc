# Default TAG if not set
TAG ?= 0.1.1
IMAGE_TO_RM ?= none
# Include the .env file if it exists
-include .env

run:
	@go run ./

docker-build: increment-tag
	eval $$(minikube docker-env) && \
	echo "Building for TAG: $(TAG)" && \
	echo docker images before build: && \
	docker images && \
	docker build -t eiachh/home-foods-core:$(TAG) ./ && \
	docker tag eiachh/home-foods-core:$(TAG) eiachh/home-foods-core:latest && \
	if [ -n "$$(docker images --filter "dangling=true" -q --no-trunc)" ]; then \
		docker rmi $$(docker images --filter "dangling=true" -q --no-trunc); \
	fi && \
	docker images 

install-mongo:
	robot -t "Actually Setup" -v TEST_NAMESPACE:"mantest" -v _mansetup:"true" test/tests/common.robot

increment-tag:
	$(eval TAG := $(shell echo $(TAG) | awk -F. -v OFS=. '{$$NF += 1 ; print}'))
	@echo "New TAG: $(TAG)"

	@echo "TAG=$(TAG)" > .env

remove-docker-image:
	eval $$(minikube docker-env) && \
	docker rmi $(IMAGE_TO_RM)