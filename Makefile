echo:
	@go build -o bin/echo ./cmd/echo

cat:
	@go build -o bin/cat ./cmd/cat

ls:
	@go build -o bin/ls ./cmd/ls

all: echo cat ls

clean:
	@rm -f bin/echo bin/cat bin/ls

test:
	@go test ./... -v