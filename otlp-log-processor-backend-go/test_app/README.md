# test-app

The program configures OpenTelemetry logging and exports logs to a local OTLP collector at `localhost:4317`. It sets up
a LoggerProvider with resource attributes (including service name, namespace, version, and a custom attribute,
`resource-foo`), then registers it as the global provider. Using the otelslog bridge, it creates a slog.Logger that
emits an info log every 100ms, randomly tagging each entry with either a "foo" or "bar" attribute and a small set of
random values. This simulates a steady stream of structured logs for testing the pipeline.
