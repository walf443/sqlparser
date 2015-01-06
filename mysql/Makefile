test: *.go parser.go
	go test -v ./...

clean:
	rm parser.go *.output

.PHONY: test clean

parser.go: parser.go.y
	go tool yacc -o parser.go -v parser.output parser.go.y
