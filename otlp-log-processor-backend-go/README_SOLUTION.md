# Notes about the solution

This document is a collection of notes about the solution, such as things that I consider worth adding, but that I
didn't add due to time constraints.

## GitHub setup

Following is needed:

1. Dependabot configuration, together with automatically merging PRs that perform minor upgrades.
2. PR checks (tests, linters)
3. Automatic releases (following a tag push for example)

## Shutdown/startup behavior

This service is responsible for counting records it received, not the records that were sent (for whatever reason). From
the counting service's point of view, it would make sense to handle shutdown and startup behavior (by pausing counting
when the service shuts down, and resuming when it starts up).

There are also arguments against handling this behavior, i.e. simply counting from 0 whenever the service is started
again.

It might give the wrong impression that everything is OK. Let's assume the counting window is 10 seconds,
and the service was down for 7 seconds. If we handle the shutdown, the counting service will show counts for the 3
seconds it was up. That might be a problem if the person using the counting service isn't aware of that.

Even if that person is aware of that, the counts are partial, and hence do not provide the full picture.

Summarizing everything, I'd rather not handle this behavior for now.