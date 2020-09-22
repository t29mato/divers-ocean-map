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
	echo '{"id":"66936b3e-08e3-404b-815d-ddbccfb03cc9"}' | sam local invoke -n ./env.json