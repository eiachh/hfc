TAG ?= 0.1.1

run:
	@go run ./

docker-build:
	eval $$(minikube docker-env) && \
	docker build -t eiachh/home-foods-core:$(TAG) ./ && \
	docker tag eiachh/home-foods-core:$(TAG) eiachh/home-foods-core:latest
	if [ -n "$$(docker images --filter "dangling=true" -q --no-trunc)" ]; then \
		docker rmi $$(docker images --filter "dangling=true" -q --no-trunc); \
	fi && \
	docker images

install-mongo:
	robot -t "Actually Setup" -v TEST_NAMESPACE:"mantest" -v _mansetup:"true" test/tests/common.robot