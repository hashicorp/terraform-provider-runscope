GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)

default: build test testacc

test: fmtcheck
	go test -v . ./plugin/providers/runscope

testacc: fmtcheck
	@test "${RUNSCOPE_ACCESS_TOKEN}" || (echo '$$RUNSCOPE_ACCESS_TOKEN required' && exit 1)
	@test "${RUNSCOPE_TEAM_ID}" || (echo '$$RUNSCOPE_TEAM_ID required' && exit 1)

	go test -v ./plugin/providers/runscope -run="TestAcc" -timeout 20m

build: fmtcheck vet testacc
	@go install
	@mkdir -p ~/.terraform.d/plugins/
	@cp $(GOPATH)/bin/terraform-provider-runscope ~/.terraform.d/plugins/terraform-provider-runscope
	@echo "Build succeeded"

build-gox: deps fmtcheck vet
	gox -osarch="linux/amd64 windows/amd64 darwin/amd64" \
	-output="pkg/{{.OS}}_{{.Arch}}/terraform-provider-runscope" .

release:
	go get github.com/goreleaser/goreleaser; \
	goreleaser; \

deps:
	go get -u github.com/mitchellh/gox

clean:
	rm -rf pkg/
fmt:
	gofmt -w $(GOFMT_FILES)

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

vet:
	@echo "go vet ."
	@go vet $$(go list ./... | grep -v vendor/) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

.PHONY: build test testacc vet fmt fmtcheck errcheck vendor-status test-compile
