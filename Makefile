build:
	go build -o esql-ast cmd/esql-ast/main.go

run-sample:
	./esql-ast -f examples/sample.esql -pretty
