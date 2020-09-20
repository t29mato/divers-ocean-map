## Test
.PHONY: test
test:
	cd hello-world && go test hello-world/...

## Build
.PHONY: build
build: test
	sam build

## Run Local
.PHONY: run
run: build
	sam local invoke -n ./env.json