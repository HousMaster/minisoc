.PHONY: all build run clean 

NAME=minisoc

all: build run

build:
	go build -o ./bin/$(NAME) ./cmd/$(NAME)

run:
	CONFIG_FILE="config/local.yaml" ./bin/$(NAME)

clean:
	rm ./bin/$(NAME)