## Test
.PHONY: test
test:
	cd scraping && go test scraping/...

## Build
.PHONY: build
build: test
	sam build

## Run Local
.PHONY: run
run: build
	sam local invoke -n ./env.json