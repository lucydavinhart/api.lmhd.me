.PHONY: build clean deploy gomodgen

.DEFAULT_GOAL := default
default: deploy

build: gomodgen
	@echo ========================================
	@echo Building
	@echo ========================================
	export GO111MODULE=on
	env GOOS=linux go build -ldflags="-s -w" -o bin/hello hello/main.go

clean:
	@echo ========================================
	@echo Cleaning
	@echo ========================================
	rm -rf ./bin ./vendor Gopkg.lock

deploy: clean build
	@echo ========================================
	@echo Deploying
	@echo ========================================
	sls deploy --verbose

gomodgen:
	@echo ========================================
	@echo GoModGen
	@echo ========================================
	chmod u+x gomod.sh
	./gomod.sh
