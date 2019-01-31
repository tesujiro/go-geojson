.PHONY: test
test:
	go vet .
	go test -v .

.PHONY: cover
cover:
	go test -v . -coverpkg . -coverprofile=cover.out && go tool cover -html=cover.out -o cover.html
