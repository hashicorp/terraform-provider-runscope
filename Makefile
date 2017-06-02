test:
	go test -v . ./plugin/providers/runscope

testacc:
	@test "${RUNSCOPE_ACCESS_TOKEN}" || (echo '$$RUNSCOPE_ACCESS_TOKEN required' && exit 1)
	@test "${RUNSCOPE_TEAM_ID}" || (echo '$$RUNSCOPE_TEAM_ID required' && exit 1)

	TF_ACC=1 go test -v ./plugin/providers/runscope -run="TestAcc" -timeout 20m

build: deps
	gox -osarch="linux/amd64 windows/amd64 darwin/amd64" \
	-output="pkg/{{.OS}}_{{.Arch}}/terraform-provider-runscope" .

deps:
	go get -u github.com/mitchellh/gox

clean:
	rm -rf pkg/
