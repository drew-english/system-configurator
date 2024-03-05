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
    go test -v ./...

@fmt:
    go fmt ./...

ginkgo-test:
    ginkgo ./...

ginkgo-bootstrap DIR:
    #!/usr/bin/env bash
    cd {{DIR}} && ginkgo bootstrap

ginkgo-generate DIR TEST_NAME:
    #!/usr/bin/env bash
    cd {{DIR}} && ginkgo generate {{TEST_NAME}}

gh-ci WORKFLOW='ci-multi-platform':
    #!/usr/bin/env bash
    gh workflow run {{WORKFLOW}} --ref $(git rev-parse --abbrev-ref HEAD)
    if [[ $? -ne 0 ]]; then
        echo "Failed to start workflow"
        exit 1
    fi

    echo "Waiting for workflow to spawn..."
    sleep 5
    gh run watch $(gh run list --workflow={{WORKFLOW}} --jq '.[0].databaseId' --json databaseId)
