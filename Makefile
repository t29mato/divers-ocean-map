## Test
.PHONY: test
test:
	cd hello-world && go test hello-world/...

## Build
.PHONY: build
build:
	sam build
