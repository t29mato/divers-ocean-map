## Test
.PHONY: test
test:
	cd functions/scraping && go test scraping/...
	cd functions/api && go test api/...

## Build
.PHONY: build
build: test
	sam build

## Run Scraping Function on Local
.PHONY: run-scraping
run-scraping: build
	echo '{"id":"66936b3e-08e3-404b-815d-ddbccfb03cc9"}' | sam local invoke ScrapingFunction -n ./env.json

## Run API on Local
.PHONY: run-api
run-api: build
	sam local start-api -n ./env.json

## Deploy
.PHONY: deploy
deploy: build
	sam deploy --profile t2k