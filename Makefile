.PHONY: build clean deploy gomodgen

.DEFAULT_GOAL := default
default: clean build deploy-dev

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

deploy-prod:
	@echo ========================================
	@echo Deploying to Prod
	@echo ========================================
	sls deploy --verbose --stage=prod

gomodgen:
	@echo ========================================
	@echo GoModGen
	@echo ========================================
	chmod u+x gomod.sh
	./gomod.sh
