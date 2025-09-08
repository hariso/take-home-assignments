# Notes about the solution

This document is a collection of notes about the solution, such as things that I consider worth adding, but that I
didn't add due to time constraints.

## GitHub setup

The following is needed:

1. Dependabot configuration, together with automatically merging PRs that perform minor upgrades.
2. PR checks (tests, linters).
3. Automatic releases (following a tag push, for example).

## Makefile

This solution adds a `Makefile` with a few useful targets, such as `make test` for running tests, and `make build`, for
building the app.

## Test app

In the [`test_app`](/test_app) folder, there's a test application that can be used to test the counting service. To run
it, execute `go run main.go`

## Observability

1. More metrics can be added so that we can watch the number of log records received, per resource/scope.
2. Similarly, we should track the time needed to count records (per request/resource/scope).

## Tests

More tests should be added. Examples would be:

1. More benchmark tests (e.g., tests with many requests).
2. More tests for the counter, printer, config.

## Shutdown/startup behavior

This is about handling the counting service's shutdown and start-up, and how it affects counting logs. After the service
starts, we may choose to continue counting or reset counts to zero.

### Pros

This service is responsible for counting records it received, not the records that were sent (for whatever
reason). From the counting service's point of view, it would make sense to handle shutdown and startup behavior (by
pausing counting when the service shuts down, and resuming when it starts up).

### Cons 

It might give the wrong impression that everything is OK. Let's assume the counting window is 10 seconds,
and the service was down for 7 seconds. If we handle the shutdown, the counting service will show counts for the 3
seconds it was up. That might be a problem if the person using the counting service isn't aware of that.

Even if that person is aware of that, the counts are partial, and hence do not provide the full picture.

Summarizing everything, I'd rather not handle this behavior for now.

