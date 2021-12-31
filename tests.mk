#!/usr/bin/make -f

########################################
### Testing

BINDIR ?= $(GOPATH)/bin

## required to be run first by most tests
build_docker_test_image:
	docker build -t tester -f ./test/docker/Dockerfile .
.PHONY: build_docker_test_image

### coverage, app, persistence, and libs tests
test_cover:
	# run the go unit tests with coverage
	bash test/test_cover.sh
.PHONY: test_cover

test_apps:
	# run the app tests using bash
	# requires `msm-cli` and `augusteum` binaries installed
	bash test/app/test.sh
.PHONY: test_apps

test_msm_apps:
	bash msm/tests/test_app/test.sh
.PHONY: test_msm_apps

test_msm_cli:
	# test the cli against the examples in the tutorial at:
	# ./docs/msm-cli.md
	# if test fails, update the docs ^
	@ bash msm/tests/test_cli/test.sh
.PHONY: test_msm_cli

test_integrations:
	make build_docker_test_image
	make tools
	make install
	make test_cover
	make test_apps
	make test_msm_apps
	make test_msm_cli
	make test_libs
.PHONY: test_integrations

test_release:
	@go test -tags release $(PACKAGES)
.PHONY: test_release

test100:
	@for i in {1..100}; do make test; done
.PHONY: test100

vagrant_test:
	vagrant up
	vagrant ssh -c 'make test_integrations'
.PHONY: vagrant_test

### go tests
test:
	@echo "--> Running go test"
	@go test -p 1 $(PACKAGES) -tags deadlock
.PHONY: test

test_race:
	@echo "--> Running go test --race"
	@go test -p 1 -v -race $(PACKAGES)
.PHONY: test_race

test_deadlock:
	@echo "--> Running go test --deadlock"
	@go test -p 1 -v  $(PACKAGES) -tags deadlock 
.PHONY: test_race
