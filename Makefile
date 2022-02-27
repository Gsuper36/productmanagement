run:
	docker run -it -p 50051:50051 \
	--add-host=host.docker.internal:host-gateway \
	--env-file=./service/.env.example \
	grpc-productmanagement-service

build:
	docker image build \
	-f service/Dockerfile \
	-t grpc-productmanagement-service:latest \
	. 