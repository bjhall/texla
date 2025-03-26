main: main.go enums
	go build -C ~/dev/texla -o texla main.go

test: tests/runner.go
	cd tests && go run runner.go && cd ..

enums: parser/enums
	generators/gen_enums parser/enums > parser/enums.go
