alias r := run
alias t := test
alias gt := ginkgo-test
alias gb := ginkgo-bootstrap

default: run

run *ARGS:
    go run main.go {{ARGS}}

test:
    go test ./...

ginkgo-test:
    ginkgo ./...

ginkgo-bootstrap DIR:
    #!/usr/bin/env bash
    cd {{DIR}} && ginkgo bootstrap
