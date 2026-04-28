_i am reading Systems Design Interview: insiders guide and saw this problem. i asked GPT to make an assignment for me where i implemet this problem. that is the extent of ai used in this project_

# Assignment: Build a Rate Limiter From Scratch

## Context

This assignment is inspired by the rate limiter system design problem from *System Design Interview: An Insider's Guide*. The goal is to move from the high-level interview prompt into a working, tested implementation that demonstrates algorithmic reasoning, API design, concurrency control, and practical software engineering discipline.

You will implement the project yourself. This document intentionally contains requirements, expectations, and test cases only. It does not include application code.

## Chosen Stack

Use **Go** for this assignment.

Recommended tool stack:

- Go 1.22 or newer
- Go standard library only for the core implementation
- `net/http` for the demo HTTP server or middleware
- `testing` for unit and integration tests
- `sync` primitives for concurrency safety
- Optional: `go test -race` for race-condition checks

Do not use an existing rate limiter package. The point of the assignment is to implement the rate limiting logic yourself.

## Learning Objectives

By the end of this assignment, you should be able to:

- Explain the difference between fixed window, sliding window, leaky bucket, and token bucket rate limiters.
- Implement a correct in-memory token bucket rate limiter.
- Apply rate limits per client, user, API key, or other request key.
- Make the limiter safe under concurrent access.
- Expose rate limiting behavior through a small HTTP middleware or handler.
- Write deterministic tests for time-dependent behavior.
- Analyze tradeoffs between correctness, memory usage, latency, and distributed scalability.

## Required Algorithm

Implement a **token bucket** rate limiter.

The token bucket should support these concepts:

- **Capacity**: maximum number of tokens the bucket can hold.
- **Refill rate**: how quickly tokens are added over time.
- **Consumption**: each allowed request consumes one token.
- **Rejection**: requests are rejected when no tokens are available.
- **Per-key buckets**: each client or key gets an independent bucket.
- **Maximum refill cap**: the bucket must never exceed its configured capacity.

Example configuration:

- Capacity: 5 tokens
- Refill rate: 1 token per second
- Meaning: a client can make 5 requests immediately, then 1 additional request per second after that.

## Functional Requirements

Your implementation must provide a reusable rate limiter component with the following behavior:

- A caller can ask whether a request for a given key is allowed.
- An allowed request consumes one token.
- A rejected request does not consume a token.
- Each key has an independent bucket.
- Buckets refill based on elapsed time.
- Buckets never refill past capacity.
- The limiter is safe to call from multiple goroutines concurrently.
- The limiter can report useful metadata such as remaining tokens or retry delay.

## HTTP Requirements

Build a small HTTP demonstration around the limiter.

The HTTP layer should:

- Apply rate limiting before calling the protected handler.
- Identify the client using a simple keying strategy, such as IP address or a request header like `X-API-Key`.
- Return HTTP `200 OK` for allowed requests.
- Return HTTP `429 Too Many Requests` for rejected requests.
- Include a `Retry-After` header or equivalent response metadata when a request is rejected.
- Optionally include rate limit headers such as `X-RateLimit-Limit`, `X-RateLimit-Remaining`, and `X-RateLimit-Reset`.

## Non-Functional Requirements

Your solution should satisfy these engineering expectations:

- The core rate limiter logic should be separate from the HTTP layer.
- Tests should cover the core limiter without needing to start an HTTP server.
- Time-dependent tests should be deterministic. Prefer an injectable clock or equivalent testing strategy instead of relying heavily on real sleeps.
- Shared state must be protected against data races.
- The implementation should be simple and readable before being highly optimized.
- Error cases and invalid configuration should be handled clearly.
- The project should include a short README explaining how to run tests and start the demo server.

## Constraints

- Do not use third-party rate limiter libraries.
- Do not use Redis, databases, queues, or external services for the required implementation.
- Keep the required implementation in-memory.
- You may use the Go standard library freely.
- You may add optional third-party tools only for development convenience, such as linters, but not for the limiter itself.

## Required Test Cases

Your implementation should pass tests equivalent to the cases below. You may name the tests however you like.

### 1. Allows Requests Within Capacity

Configuration:

- Capacity: 5
- Refill rate: 1 token per second
- Key: `client-a`

Expected behavior:

- The first 5 immediate requests are allowed.
- The 6th immediate request is rejected.
- Remaining tokens after the 5th allowed request should be 0.

### 2. Refills Tokens Over Time

Configuration:

- Capacity: 5
- Refill rate: 1 token per second
- Key: `client-a`

Scenario:

- Make 5 allowed requests immediately.
- Confirm the 6th immediate request is rejected.
- Advance time by 1 second.
- Make another request.

Expected behavior:

- The request after 1 second is allowed.
- Remaining tokens after that request should be 0.

### 3. Does Not Exceed Capacity After Long Idle Period

Configuration:

- Capacity: 5
- Refill rate: 1 token per second
- Key: `client-a`

Scenario:

- Make 5 allowed requests immediately.
- Advance time by 100 seconds.
- Make requests until rejection.

Expected behavior:

- Exactly 5 requests are allowed after the idle period.
- The 6th request after the idle period is rejected.
- The bucket does not accumulate more than 5 tokens.

### 4. Maintains Separate Buckets Per Key

Configuration:

- Capacity: 2
- Refill rate: 1 token per second
- Keys: `client-a`, `client-b`

Scenario:

- Make 2 allowed requests for `client-a`.
- Confirm the 3rd request for `client-a` is rejected.
- Make 2 requests for `client-b`.

Expected behavior:

- `client-a` is limited independently.
- `client-b` still has its full capacity and both requests are allowed.

### 5. Does Not Consume Tokens On Rejection

Configuration:

- Capacity: 1
- Refill rate: 1 token per second
- Key: `client-a`

Scenario:

- Make 1 allowed request.
- Make 3 rejected requests immediately after.
- Advance time by 1 second.
- Make another request.

Expected behavior:

- The first request is allowed.
- The next 3 requests are rejected.
- The request after 1 second is allowed.
- Rejected requests do not push the bucket into a negative state.

### 6. Calculates Retry Delay For Rejected Requests

Configuration:

- Capacity: 1
- Refill rate: 1 token per second
- Key: `client-a`

Scenario:

- Make 1 allowed request.
- Immediately make another request.

Expected behavior:

- The second request is rejected.
- The limiter reports a retry delay of approximately 1 second.
- The retry delay should never be negative.

### 7. Handles Concurrent Requests Safely

Configuration:

- Capacity: 10
- Refill rate: 1 token per second
- Key: `client-a`

Scenario:

- Start 100 concurrent requests for the same key at the same logical time.

Expected behavior:

- Exactly 10 requests are allowed.
- Exactly 90 requests are rejected.
- Running the test with `go test -race` reports no data races.

### 8. Rejects Invalid Configuration

Scenarios:

- Capacity is 0.
- Capacity is negative.
- Refill rate is 0.
- Refill rate is negative.

Expected behavior:

- The limiter refuses invalid configuration.
- The failure mode is explicit, such as returning an error from a constructor.

### 9. HTTP Allows Requests Under Limit

Configuration:

- Capacity: 2
- Refill rate: 1 token per second
- Key source: `X-API-Key`
- Header value: `demo-key`

Scenario:

- Send 2 HTTP requests with `X-API-Key: demo-key`.

Expected behavior:

- Both responses return HTTP `200 OK`.
- The protected handler is called exactly twice.

### 10. HTTP Rejects Requests Over Limit

Configuration:

- Capacity: 2
- Refill rate: 1 token per second
- Key source: `X-API-Key`
- Header value: `demo-key`

Scenario:

- Send 3 immediate HTTP requests with `X-API-Key: demo-key`.

Expected behavior:

- The first 2 responses return HTTP `200 OK`.
- The 3rd response returns HTTP `429 Too Many Requests`.
- The rejected response includes retry metadata, preferably a `Retry-After` header.
- The protected handler is not called for the rejected request.

### 11. HTTP Separates Different Clients

Configuration:

- Capacity: 1
- Refill rate: 1 token per second
- Key source: `X-API-Key`

Scenario:

- Send 1 request with `X-API-Key: client-a`.
- Send 1 request with `X-API-Key: client-b`.
- Send another immediate request with `X-API-Key: client-a`.

Expected behavior:

- The first `client-a` request returns HTTP `200 OK`.
- The `client-b` request returns HTTP `200 OK`.
- The second `client-a` request returns HTTP `429 Too Many Requests`.

## Suggested Project Structure

You may choose your own structure, but a clean submission should separate concerns clearly.

Suggested layout:

- `limiter/`: core token bucket logic
- `http/` or `middleware/`: HTTP integration
- `cmd/server/`: optional demo server entry point
- `README.md`: setup and usage instructions
- Tests near the packages they exercise

This is a suggested organization, not a strict requirement.

## Design Questions To Answer In Your README

Include short answers to these questions:

- Why did you choose token bucket instead of fixed window or sliding window?
- What happens when many clients use the limiter at the same time?
- How does your implementation avoid data races?
- How would your design need to change for multiple application servers?
- What metadata should a rejected client receive?
- What are the memory growth risks of per-key buckets?
- How could idle buckets be cleaned up?

## Evaluation Rubric

Total: 100 points

- Correct token bucket behavior: 25 points
- Per-key isolation: 10 points
- Time-based refill correctness: 15 points
- Concurrency safety: 15 points
- HTTP integration: 10 points
- Test coverage and determinism: 15 points
- Code organization and readability: 5 points
- README design explanation: 5 points

## Optional Extensions

After completing the required version, you may add one or more extensions:

- Implement a fixed window limiter and compare its behavior to token bucket.
- Implement a sliding window log or sliding window counter limiter.
- Add cleanup for idle client buckets.
- Add configurable costs, where some requests consume more than one token.
- Add a small benchmark suite.
- Add structured logs for allowed and rejected requests.
- Add Prometheus-style metrics counters.
- Sketch a distributed design using Redis, but do not make this part of the required implementation.

## Suggested Timeline

Estimated total time: **12 to 18 focused hours**.

Suggested schedule:

- Day 1, 1 to 2 hours: Review token bucket, fixed window, sliding window, and leaky bucket algorithms.
- Day 2, 2 to 3 hours: Design the core API, data structures, configuration validation, and time strategy.
- Day 3, 3 to 4 hours: Implement the core in-memory token bucket and basic unit tests.
- Day 4, 2 to 3 hours: Add deterministic time-based tests and per-key bucket tests.
- Day 5, 2 to 3 hours: Add concurrency safety and verify with normal tests and race detection.
- Day 6, 2 to 3 hours: Build the HTTP middleware or handler and add HTTP tests.
- Day 7, 1 to 2 hours: Write the README, answer design questions, and clean up naming and organization.
- Day 8, 1 hour: Final verification, run all tests, and review edge cases.

If you are new to Go or concurrent programming, budget closer to **20 to 25 hours**.

## Completion Criteria

The assignment is complete when:

- All required test cases pass.
- The core limiter is independent from the HTTP layer.
- The implementation is safe under concurrent access.
- `go test ./...` passes.
- `go test -race ./...` passes.
- The README explains how to run and reason about the project.
- You can explain the design tradeoffs without reading directly from the code.
