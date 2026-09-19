# ADR 0001: Rewrite the CRM backend in Go instead of staying on Node/NestJS

## Status
Accepted

## Context
GloveBoxCRM's backend runs on NestJS today. As the platform grows toward higher-throughput needs (active policy monitoring, real-time carrier sync, analytics on tens of thousands of service transactions), the team needs to decide whether to keep optimizing the Node stack or invest in a unified compiled rewrite.

The job description explicitly frames codebase fragmentation as a top challenge and requires establishing high-performance backend standards.

## Decision
Migrate the backend to **Go (Golang)** incrementally using a strangler-fig pattern (see ADR 0003) rather than a big-bang rewrite. New modules are built in Go from day one; existing modules migrate module-by-module as their throughput or reliability requirements dictate.

## Alternatives & Language Stacks Considered

### 1. Stay on NestJS & Invest in Performance Tuning
- **Approach**: Keep the existing Node.js/TypeScript monolith. Tune event loop handlers, add Redis caching layers, and use worker threads for concurrency.
- **Pros**: Zero short-term language ramp-up; retains 100% existing TypeScript feature velocity.
- **Cons**: Single-threaded Event Loop stalls on CPU-heavy data transformations or carrier encryption payloads; runtime type erasure (`tsc` strips types) allows runtime bugs to slip past build steps; heavy deployment container sizes (`node_modules`).

### 2. Rewrite in Rust
- **Approach**: Rebuild the backend in Rust for maximum memory safety and bare-metal CPU performance.
- **Pros**: Industry-leading memory safety, zero-cost abstractions, and maximum CPU efficiency.
- **Cons**: Steep learning curve (borrow checker, lifetime syntax) slows engineering iteration speed; over-engineered for a network-bound CRM domain where DB and HTTP I/O dominate CPU time.

### 3. Rewrite in Java / Spring Boot
- **Approach**: Migrate to enterprise Java with Spring Boot.
- **Pros**: Massive enterprise library ecosystem, mature ORMs, and deep battle-tested integration frameworks.
- **Cons**: High JVM memory overhead (hundreds of MBs idle vs ~15MB for Go); heavier application ceremony and slower container cold start times.

## Real-World Industry Precedents
- **Uber & Twitch**: Replaced legacy Node.js/Python microservices with Go to handle high-concurrency real-time connections and achieve predictable p99 latency SLAs.
- **Tailscale & Cloudflare**: Standardized core control planes and edge proxy logic on Go for single-binary portability and lightweight concurrency.

## Consequences
+ **Lightweight Concurrency**: Goroutines (`~2KB` per stack) easily handle thousands of concurrent carrier sync tasks with `context.Context` cancellation.
+ **Static Typing & Compile-Time Safety**: Catches null pointers, structural type mismatches, and unhandled errors prior to deployment.
+ **Lean Deployments**: Compiles to a single static binary (~15MB), drastically reducing container sizes and cold-start latency.
- **Ramp-Up Curve**: Engineering team must adapt to Go idioms and explicit error handling (`if err != nil`).
- **Coexistence Surface Area**: Two application runtimes (Node and Go) live in production during the migration window (mitigated by ADR 0003's Gateway proxy).
