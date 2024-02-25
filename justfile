# Justfile for common project tasks
# Docs: https://github.com/casey/just

alias r := run
alias t := test
alias gt := ginkgo-test
alias gb := ginkgo-bootstrap
alias gg := ginkgo-generate

default: run

run *ARGS:
    go run main.go {{ARGS}}

test:
    go test ./...

fmt:
    go fmt ./...

ginkgo-test:
    ginkgo ./...

ginkgo-bootstrap DIR:
    #!/usr/bin/env bash
    cd {{DIR}} && ginkgo bootstrap

ginkgo-generate DIR TEST_NAME:
    #!/usr/bin/env bash
    cd {{DIR}} && ginkgo generate {{TEST_NAME}}
