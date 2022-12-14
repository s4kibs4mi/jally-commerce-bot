.PHONY: deps build run all

all: deps build run

export CONFIG_DIR=.

deps:
	@echo "Running get deps..."
	@$(go get .)

build: deps
	@echo "Running build..."
	GO111MODULE=on GOARCH="amd64" CGO_ENABLED=0 go build -v -o jally-commerce-bot

run:
	@echo "Running server..."
	./jally-commerce-bot serve
