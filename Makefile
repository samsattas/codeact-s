.PHONY: build web agent test vet clean

# Build both binaries into bin/.
build:
	go build -o bin/agent ./cmd/agent
	go build -o bin/web ./cmd/web

# Build and run the web UI in one command.
# Extra flags: make web ARGS="-workdir ./some/project -addr :9090"
web:
	go build -o bin/web ./cmd/web
	./bin/web $(ARGS)

# Build and run the CLI in one command.
# Extra flags/task: make agent ARGS="-workdir . \"list the .go files\""
agent:
	go build -o bin/agent ./cmd/agent
	./bin/agent $(ARGS)

test:
	go test ./...

vet:
	go vet ./...

clean:
	rm -rf bin/
