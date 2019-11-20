.PHONY: build dependencies doc fmt install lint push server test watch

build:
	go build

dependencies:
	go get -u github.com/tools/godep
	godep restore
	rm -rf Godeps
	godep save ./...

doc:
	godoc -http=:6060

fmt:
	for package in $$(go list ./... | grep -v /vendor/); do go fmt $$package; done

install: dependencies

lint:
	for package in $$(go list ./... | grep -v /vendor/); do golint $$package; done

push:
	cf login
	cf push

server:
	go run server.go

test:
	ginkgo -r

watch:
	ginkgo watch -r
