GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
PKG_NAME=runscope
WEBSITE_REPO=github.com/hashicorp/terraform-website

default: lint build test testacc

test: goimportscheck
	go test -v . ./runscope

testacc: goimportscheck
	@test "${RUNSCOPE_ACCESS_TOKEN}" || (echo '$$RUNSCOPE_ACCESS_TOKEN required' && exit 1)
	@test "${RUNSCOPE_TEAM_ID}" || (echo '$$RUNSCOPE_TEAM_ID required' && exit 1)

	go test -v -count=1 ./runscope -run="TestAcc" -timeout 20m -sweep "eu-west-1"

build: goimportscheck vet
	@go install
	@mkdir -p ~/.terraform.d/plugins/
	@cp $(GOPATH)/bin/terraform-provider-runscope ~/.terraform.d/plugins/terraform-provider-runscope
	@echo "Build succeeded"

build-gox: deps goimportscheck vet
	gox -osarch="linux/amd64 windows/amd64 darwin/amd64" \
	-output="pkg/{{.OS}}_{{.Arch}}/terraform-provider-runscope" .

release:
	GO111MODULE=off go get github.com/goreleaser/goreleaser; \
	goreleaser; \

deps:
	GO111MODULE=off go get -u github.com/mitchellh/gox

clean:
	rm -rf pkg/
	
fmt:
	gofmt -w $(GOFMT_FILES)

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

errcheck:
	@sh -c "'$(CURDIR)/scripts/errcheck.sh'"


goimports:
	goimports -w $(GOFMT_FILES)

goimportscheck:
	@sh -c "'$(CURDIR)/scripts/goimportscheck.sh'"

vet:
	@echo "go vet ."
	@go vet $$(go list ./... | grep -v vendor/) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

lint:
	@echo "go lint ."
	@golint $$(go list ./... | grep -v vendor/) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Lint found errors in the source code. Please check the reported errors"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi
	
test-compile:
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package. For example,"; \
		echo "  make test-compile TEST=./$(PKG_NAME)"; \
		exit 1; \
	fi
	go test -c $(TEST) $(TESTARGS)

website:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), get-ting..."
	git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
endif
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)

website-test:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), get-ting..."
	git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
endif
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider-test PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)

.PHONY: build test testacc vet goimports goimportscheck fmt fmtcheck errcheck errcheck lint test-compile website website-test
