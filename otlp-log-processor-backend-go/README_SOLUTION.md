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
it, execute `go run main.go`. It runs a loop that emits records with a few random attributes.

## Observability

The following is worth adding:

1. More metrics can be added so that we can watch the number and size of log records received, per resource/scope.
2. Time needed to count records (per request/resource/scope).
3. A health check endpoint.

## Tests

More tests should be added. Examples would be:

1. More benchmark tests (e.g., tests with many requests). Tests should be added for both, `dash0LogsServiceServer` and
   the `counter` struct.
2. More tests for the counter, printer, config.

## Shutdown/startup behavior

This is about handling the counting service's shutdown and start-up, and how it affects counting logs. Generally
speaking, a service will want to save its progress when shutting down and resume it when starting up. When the logs'
service shuts down, we can save the attribute value counts and resume counting when it starts up.

In this particular case, that would mean that partial counts are shown to a user. That's not always desirable. In cases
when it is possible to use partial counts, a user needs to be able to know that's the case. The logs counting service
can do that (e.g., by printing "these counts might be partial", or setting a "partial" flag in the response).

Not having any counts at all makes it very clear that something happened to the service.

Summarizing everything: IMHO, it's better not show any partial counts, since more often than not, the user will want to
see the full counts.