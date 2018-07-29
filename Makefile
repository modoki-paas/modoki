BINARY_NAME = modoki
DOCKER_IMAGE_NAME = modokipaas/modoki

SRCS = $(wildcard *.go)

.PHONY: all  build  clean
all: build docker

build: $(SRCS)
	go build -o $(BINARY_NAME) 

docker: $(SRCS) Dockerfile
	docker build -t $(DOCKER_IMAGE_NAME) .

clean:
	rm $(BINARY_NAME)