# Contributing to Deployadactyl

Deployadactyl is an open source project and we welcome all contributions!

## Requirements

If this is your first contribution we require you to [sign our Contributing License Agreement](https://compozed-cla.cfapps.io/agreements/compozed/deployadactyl "Compozed CLA").

> Note: Do not make commits through the GitHub web interface due to issues with the automated CLA management.

Ensure tests have been added for your changes. If you need help writing tests, send us a pull request with your changes and we will help you out. We use [Ginkgo](https://github.com/onsi/ginkgo) and [Gomega](https://github.com/onsi/gomega) to write our tests using [behavior driven development](https://en.wikipedia.org/wiki/Behavior-driven_development).

## Requesting Features

[Make an issue](https://github.com/compozed/deployadactyl/issues/new)

## Making changes

Following these steps will help you get your pull request accepted:

- [Fork Deployadactyl](https://github.com/compozed/deployadactyl/compare#fork-destination-box). Pull and checkout to the `develop` branch to ensure you have the latest commits

- Create a topic branch where you want to base your work: `git checkout -b fix_some_issue`

- Ensure your code follows Go best practices found at [Effective Go](https://golang.org/doc/effective_go.html) and [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)

- Use `git rebase` (not `git merge`) to sync your work with the latest version: `git fetch upstream` && `git rebase upstream/master`

- Run **all** the tests to assure nothing else is broken: `ginkgo -r` or `go test ./...`

- [Create a pull request](https://github.com/compozed/deployadactyl/compare) against the `develop` branch

- Run `go fmt ./...`

- Ensure all pull request checks (such as continuous integration) pass
