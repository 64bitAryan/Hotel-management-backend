build:
	@go build -o bin/api

run: build
	@./bin/api

seed:
	@go run scripts/seed.go

test:
	@go test -v ./...

ARG ?= default_value

# Use the MAKECMDGOALS variable to get the command line argument
ifeq ($(MAKECMDGOALS),)
  # No arguments were passed
else
  ARG = $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  $(eval $(ARG):;@:)
endif

commit:
	@git add --all
	@git commit -m "$(ARG)"