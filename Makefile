install: # install dependencies
	@go mod tidy

test: # run all unit tests	
	@go test ./... -timeout 5s -cover -coverprofile=cover.out	