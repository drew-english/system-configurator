alias r := run
alias t := test
alias gt := ginkgo-test

default: run

run *ARGS:
    go run main.go {{ARGS}}

test:
    go test ./...

ginkgo-test:
    ginkgo ./...
