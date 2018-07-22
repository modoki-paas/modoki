BINARY_NAME = modoki
DOCKER_IMAGE_NAME = tsuzu/modoki

SRCS = $(wildcard *.go)

.PHONY: all build_only build dep clean
all: build docker

build: $(SRCS)
	go build -o $(BINARY_NAME) 

docker: $(SRCS) Dockerfile
	docker build -t $(DOCKER_IMAGE_NAME) .

clean:
	rm $(BINARY_NAME)