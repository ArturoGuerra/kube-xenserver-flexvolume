.PHONY: all build clean

GOBUILD = go build

all: clean build

clean:
	rm -rf bin

build: clean
	CGO_ENABLED=0 COOS=linux $(GOBUILD) -a -ldflags '-extldflags "-static"' -o bin/driver ./cmd

docker-build: build
	test $(DOCKERREPO)
	docker build . -t $(DOCKERREPO)

docker-push: docker-build
	test $(DOCKERREPO)
	docker push $(DOCKERREPO)

docker: build docker-build docker-push clean
