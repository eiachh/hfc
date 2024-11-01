run:
	@go run ./

docker-build:
	eval $$(minikube docker-env) && \
	docker build -t eiachh/home-foods-core ./ && \
	if [ -n "$$(docker images --filter "dangling=true" -q --no-trunc)" ]; then \
		docker rmi $$(docker images --filter "dangling=true" -q --no-trunc); \
	fi && \
	docker images