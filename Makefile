
.PHONY: deps-list deps-upgrade
deps-list:
	go list -m -u -mod=readonly all

deps-upgrade:
	go get -u -v ./...
	go mod tidy

deps-cleancache:
	go clean -modcache

# -count=1 is needed to ignore cached test results
test:
	go test ./... -count=1
	staticcheck -checks=all ./...

test-verbose:
	go test -v ./... -count=1 -coverprofile=coverage.out -covermode=atomic
	staticcheck -checks=all ./...
