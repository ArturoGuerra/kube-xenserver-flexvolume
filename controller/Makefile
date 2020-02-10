.PHONY: all build clean

GOBUILD = go build

all: clean build

clean:
	rm -rf bin

build: clean
	COOS=linux $(GOBUILD) -o bin/controller ./cmd

docker-build:
	test $(DOCKERREPO)
	docker build . -t $(DOCKERREPO)

docker-push:
	test $(DOCKERREPO)
	docker build $(DOCKERREPO)

docker: docker-build docker-push clean
