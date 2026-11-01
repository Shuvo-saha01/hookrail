# Hookrail

**A reliable webhook delivery engine built with Go and PostgreSQL.**

Hookrail is a backend infrastructure project designed to accept events, queue webhook deliveries, and reliably send them to external HTTP endpoints.

The focus is on learning and implementing real-world backend engineering concepts such as asynchronous processing, job queues, concurrency, retries, idempotency, and fault tolerance.

## 1. The Problem

Imagine you're building an e-commerce application.

When a customer completes a payment, several things might need to happen:

- Notify an inventory management service.
- Send an event to an analytics platform.
- Notify a third-party accounting system.

One approach is to make HTTP requests to all three services directly from the payment handler.

But what happens if one service is offline? What if a request times out? What if a server responds with an error?

We don't want the customer's payment request to fail just because the analytics service is unavailable.

We also don't want to lose important events.

**This is the problem Hookrail aims to solve.**

## 2. The Solution

Hookrail acts as an intermediary between an application that generates events and the external services that need to receive them.

Instead of delivering webhooks directly, an application sends an event to Hookrail.

Hookrail takes responsibility for delivering it reliably.

### Example workflow

```text
    E-commerce Application
             |
             | POST /events
             | "order.paid"
             v
       +------------+
       |  Hookrail  |
       |  Go API    |
       +------------+
             |
             | Persist event
             v
       +------------+
       | PostgreSQL |
       | Event Queue|
       +------------+
             |
             | Background workers
             v
       +------------------+
       | Webhook Delivery |
       +------------------+
             |
       +-----+------+
       |            |
       v            v
   Analytics     Inventory
   Service       Service
       |            |
    Success       Failure
                    |
                    v
               Retry Later
```

### What happens internally?

1. An application sends an event to Hookrail.
2. Hookrail validates the event and saves it in PostgreSQL.
3. Hookrail creates delivery jobs for the registered endpoints interested in that event.
4. The API responds with `202 Accepted` after the event and its delivery jobs are committed.
5. Background workers process delivery jobs independently.
6. Workers send HTTP POST requests to the destination URLs.
7. Successful deliveries are marked as completed.
8. Failed deliveries are retried with exponential backoff.
9. After the maximum number of attempts, a delivery is marked as permanently failed.

**The important idea: accepting an event and delivering it are two separate operations.**

## 3. Project Goals

The primary goal is not to build another CRUD application. It is to understand how reliable asynchronous backend systems work.

### Core engineering objectives

| Concept | What we want to learn |
|---|---|
| Event ingestion | Validating and accepting incoming events |
| Durable storage | Preventing event loss after crashes |
| Job queues | Using PostgreSQL as a persistent queue |
| Background workers | Processing jobs without blocking API requests |
| Concurrency | Running multiple workers safely |
| Retries | Recovering from temporary failures |
| Exponential backoff | Avoiding overwhelming failing services |
| Idempotency | Handling repeated requests and deliveries safely |
| Fault tolerance | Recovering from worker or server crashes |
| Observability | Tracking delivery attempts, failures, and performance |

## 4. Technology Stack

| Technology | Purpose |
|---|---|
| Go | HTTP API and background worker implementation |
| PostgreSQL | Persistent events, delivery jobs, and attempt history |
| Docker Compose | Local development environment |
| pgx | PostgreSQL connectivity and connection pooling |
| SQL | Database operations and concurrency control |
| Go `net/http` | HTTP server and outbound webhook requests |

The initial implementation intentionally avoids Kafka, Redis, and other external queue infrastructure.

PostgreSQL will serve as both the primary database and the job queue.

## 5. Planned Features

### Phase 1 — Foundation

- Set up the Go HTTP server.
- Connect to PostgreSQL.
- Implement database migrations.
- Build a health-check endpoint.
- Implement event ingestion.

### Phase 2 — Webhook Delivery

- Register destination endpoints.
- Associate event types with destinations.
- Persist delivery jobs.
- Implement background workers.
- Send HTTP POST requests to registered endpoints.
- Record delivery results.

### Phase 3 — Reliability

- Retry failed deliveries.
- Implement exponential backoff.
- Configure maximum retry attempts.
- Recover jobs interrupted by worker crashes.
- Track permanently failed deliveries.

### Phase 4 — Concurrency and Safety

- Run multiple workers concurrently.
- Use PostgreSQL `FOR UPDATE SKIP LOCKED` for job claiming.
- Prevent simultaneous processing of the same job.
- Implement job leases and recovery.
- Support event idempotency keys.
- Add webhook signatures and outbound request security.

### Phase 5 — Observability

- Track event and delivery statuses.
- Store delivery attempt history.
- Add structured logging.
- Expose operational metrics.
- Provide basic delivery inspection APIs.

## 6. Planned API

These endpoints represent the target MVP, not necessarily functionality already implemented.

### Register a webhook endpoint

`POST /endpoints`

```json
{
  "url": "https://example.com/webhooks",
  "event_types": ["order.paid", "order.refunded"]
}
```

### Publish an event

`POST /events`

```json
{
  "type": "order.paid",
  "payload": {
    "order_id": "ord_123",
    "amount": 1500,
    "currency": "INR"
  }
}
```

Expected response:

```json
{
  "id": "event-uuid",
  "status": "accepted"
}
```

HTTP status: `202 Accepted`

This confirms that Hookrail has durably accepted the event, not that its webhook deliveries have completed.

### Inspect an event

`GET /events/{id}`

Retrieve an event and its delivery status.

### Inspect delivery attempts

`GET /deliveries/{id}`

Retrieve the status and attempt history of an individual webhook delivery.

## 7. Reliability Guarantees

Hookrail aims to provide the following guarantees:

**Durable acceptance**

Once an event is accepted, its data and corresponding delivery jobs are persisted in PostgreSQL.

**At-least-once delivery**

A webhook may be delivered more than once. This is intentional because retries and worker crashes can cause duplicate deliveries.

**Failure isolation**

An unavailable destination should not prevent independent destinations from receiving their webhooks.

**Crash recovery**

A worker crashing during delivery should not permanently lose the job.

**Bounded retries**

Failed deliveries are retried according to a configured policy rather than indefinitely.

These are target guarantees that will be implemented and tested incrementally.

## 8. What Hookrail Is Not

To keep the project focused, the initial version will not include:

- A complex frontend dashboard.
- Kafka or Redis-based queues.
- A distributed microservices architecture.
- Billing or user management.
- Enterprise-scale throughput guarantees.
- Exactly-once webhook delivery.

The goal is to build a small but well-engineered system before adding complexity.

## 9. Definition of Done (MVP)

The MVP is complete when we can demonstrate the following scenario:

1. Start Hookrail using Docker Compose.
2. Register a webhook destination.
3. Submit an event through the API.
4. Observe a worker delivering the event.
5. Shut down the receiving service.
6. Submit another event.
7. Observe failed attempts and scheduled retries.
8. Restart the receiving service.
9. Observe the pending delivery eventually succeed.
10. Verify the complete delivery history in PostgreSQL.

The system should also recover correctly when a Hookrail worker itself crashes during processing.

## 10. Why Build Hookrail?

Hookrail is a learning-focused infrastructure project intended to develop practical experience with:

- Production-oriented Go programming.
- PostgreSQL transactions and row-level locking.
- Queue-based asynchronous architectures.
- Concurrent worker execution.
- Distributed-systems failure scenarios.
- Reliable HTTP communication.
- Backend observability and testing.

The end result should demonstrate not just the ability to write an API, but the ability to design, implement, and test a reliable backend service.

---

**Project philosophy:** Start small, understand every component, test failure scenarios, and introduce complexity only when it solves a real problem.
