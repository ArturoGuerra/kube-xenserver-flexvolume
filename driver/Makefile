.PHONY: all build clean

GOBUILD = go build

all: clean build

clean:
	rm -rf bin

build: clean
	COOS=linux $(GOBUILD) -o bin/driver ./cmd

docker:
	test $(DOCKERREPO)

	docker build . -t $(DOCKERREPO)
	docker push $(DOCKERREPO)

