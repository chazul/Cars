Cover report
 * https://about.codecov.io/blog/getting-started-with-code-coverage-for-golang/
 * https://go.dev/blog/cover

    `go test -race -covermode=atomic -coverprofile=coverage.out`
    `go tool cover -html=coverage.out`

    *generated report will be located in /tmp/ 