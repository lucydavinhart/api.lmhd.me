.PHONY: build clean deploy gomodgen

.DEFAULT_GOAL := default
default: clean build deploy-dev test-dev



build: gomodgen
	@echo ========================================
	@echo Building
	@echo ========================================
	export GO111MODULE=on
	env GOOS=linux go build -ldflags="-s -w" -o bin/hello hello/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/v1 v1/*.go

clean:
	@echo ========================================
	@echo Cleaning
	@echo ========================================
	rm -rf ./bin ./vendor Gopkg.lock

deploy: deploy-dev
deploy-dev:
	@echo ========================================
	@echo Deploying to Dev
	@echo ========================================
	sls deploy --verbose

test-dev:
	@echo ========================================
	@echo Running Dev Tests
	@echo ========================================
	@./test.sh

deploy-prod:
	@echo ========================================
	@echo Deploying to Prod
	@echo ========================================
	sls deploy --verbose --stage=prod

test-prod:
	@echo ========================================
	@echo Running Prod Tests
	@echo ========================================
	@echo Not Implemented

gomodgen:
	@echo ========================================
	@echo GoModGen
	@echo ========================================
	chmod u+x gomod.sh
	./gomod.sh


# Install Dependencies
deps:
    ifeq ($(shell uname -s),Darwin)
		brew install npm
    endif
	npm i -g serverless@1.76.1
	npm install