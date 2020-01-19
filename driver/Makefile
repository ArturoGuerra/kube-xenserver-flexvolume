.PHONY: all build clean

GOBUILD = go build

all: clean build

clean:
	rm -rf bin

build: clean
	COOS=linux $(GOBUILD) -o bin/driver ./cmd
