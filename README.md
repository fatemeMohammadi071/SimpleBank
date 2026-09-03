# Simple Bank — Backend Service (Go)

A production-style banking backend that manages user accounts, tracks balances, and performs
money transfers between accounts with proper database transactions and locking. Built to
practice and demonstrate real-world backend engineering with **Go**, **PostgreSQL**, **gRPC**,
and a full **Docker / CI** workflow.

> This repository is a learning + portfolio project. It follows the architecture of a real
> microservice: layered code, generated type-safe SQL, dual gRPC + REST interfaces,
> token-based auth, database migrations, unit tests with mocks, and continuous integration.

---

## What this service does

- **User management** — register, login, update profile, with hashed passwords.
- **Bank accounts** — create accounts per currency, list and fetch accounts (owner-scoped).
- **Money transfers** — transfer funds between two accounts inside a single ACID transaction,
  writing transfer + ledger entries and updating both balances without deadlocks.
- **Sessions & token renewal** — refresh tokens stored server-side, access tokens renewed
  without re-login.
- **Dual API** — the same business logic is exposed as **gRPC** and as a **RESTful JSON API**
  (via gRPC-Gateway), plus auto-generated **Swagger** documentation.

---

## Tech stack

### Language & core
| Area | Technology |
|------|-----------|
| Language | **Go 1.26** |
| HTTP framework (REST) | **Gin** |
| RPC framework | **gRPC** + **Protocol Buffers** |
| REST-from-gRPC | **grpc-gateway** (single proto definition serves both gRPC and HTTP/JSON) |
| API docs | **OpenAPI / Swagger** generated from proto, served as static UI |

### Database
| Area | Technology |
|------|-----------|
| Database | **PostgreSQL 18** |
| Query layer | **sqlc** — generates type-safe Go code from raw SQL (no ORM) |
| Migrations | **golang-migrate** — versioned up/down SQL migrations, run automatically on startup from an embedded filesystem (`embed.FS`) |
| Driver | `lib/pq` |
| Schema design | **DBML** → `dbdocs` (hosted ERD) and `dbml2sql` (SQL schema) |
| Transactions | Hand-written `execTx` wrapper; transfer logic uses ordered updates to avoid deadlocks; foreign keys `DEFERRABLE INITIALLY IMMEDIATE` |

### Authentication & security
| Area | Technology |
|------|-----------|
| Password hashing | **bcrypt** (`golang.org/x/crypto`) |
| Access tokens | **PASETO** (v2, ChaCha20-Poly1305) with a JWT implementation also available behind a `token.Maker` interface |
| Auth model | Stateless access tokens + server-side **refresh-token sessions** table |
| Authorization | Per-RPC / per-route middleware; users can only touch their own accounts |
| Input validation | `go-playground/validator` + custom validators (e.g. currency codes) |

### Configuration & observability
| Area | Technology |
|------|-----------|
| Config | **Viper** — reads `app.env` and/or environment variables (12-factor friendly, works in CI with no file) |
| Logging | **zerolog** — structured JSON logs, human-readable console writer in development |
| Middleware | Custom gRPC unary interceptor and HTTP logger middleware (method, status, duration) |

### Testing
| Area | Technology |
|------|-----------|
| Test framework | Go `testing` + **testify** (assert/require) |
| DB mocking | **gomock / mockgen** — generated mock of the `Store` interface for handler unit tests |
| Coverage | `go test -v -cover ./...` |
| Concurrency tests | Transfer transaction tested with many parallel goroutines to prove no deadlock / correct balances |

### Build, packaging & CI/CD
| Area | Technology |
|------|-----------|
| Containerization | **Docker** multi-stage build (build on `golang:alpine`, run on bare `alpine`) |
| Local orchestration | **docker-compose** — API + Postgres with health checks and dependency ordering |
| CI | **GitHub Actions** — spins up Postgres service, runs migrations, runs the full test suite on every push / PR to `main` |
| Task runner | **Makefile** — `postgres`, `createdb`, `migrateup/down`, `sqlc`, `proto`, `mock`, `test`, `server`, `evans` |

---

## Architecture

```
                 ┌─────────────┐        ┌──────────────┐
   gRPC clients ─┤  gapi (gRPC) │        │ api (Gin/REST)│─ legacy REST handlers
                 └──────┬──────┘        └──────┬───────┘
   HTTP/JSON ──► grpc-gateway ──► gapi ────────┤
                        │                      │
                        ▼                      ▼
                 ┌───────────────────────────────────┐
                 │  db/sqlc  — type-safe queries       │
                 │  Store interface + TransferTx (ACID)│
                 └────────────────┬──────────────────┘
                                  ▼
                            PostgreSQL
```

- **`proto/`** – Protocol Buffer service & message definitions (source of truth for the API).
- **`pb/`** – generated Go gRPC + gateway + swagger stubs.
- **`gapi/`** – gRPC server implementation: handlers, auth, converters, interceptors.
- **`api/`** – Gin-based REST implementation (earlier iteration, kept for reference).
- **`db/migration/`** – versioned schema; embedded and applied on boot.
- **`db/query/`** + **`db/sqlc/`** – raw SQL and its generated type-safe Go.
- **`token/`** – pluggable token makers (PASETO / JWT) behind one interface.
- **`util/`** – config loading, password hashing, random test helpers.
- **`doc/`** – DBML schema, generated SQL, Swagger UI assets.

---

## Data model

`users` → `accounts` → (`entries`, `transfers`) and `sessions`.
Ledger-style design: every balance change is recorded as an immutable `entry`, and every
transfer creates a `transfer` row plus two entries, all within one transaction. A unique
constraint on `(owner, currency)` prevents duplicate accounts.

---

## Running locally

```bash
# 1. Start Postgres
make postgres && make createdb

# 2. Run the service (migrations run automatically at startup)
make server
#   gRPC   -> 0.0.0.0:9090
#   REST   -> 0.0.0.0:8080
#   Swagger-> http://localhost:8080/swagger/

# or run everything in containers
docker compose up
```

Useful commands:

```bash
make sqlc     # regenerate DB code from SQL
make proto    # regenerate gRPC/gateway/swagger from .proto
make mock     # regenerate the Store mock
make test     # run all tests with coverage
make evans    # interactive gRPC client (REPL)
```

---

## Backend concepts demonstrated

- Designing a normalized relational schema and evolving it with reversible migrations
- Writing raw SQL and keeping it type-safe without an ORM
- ACID transactions, row locking, and deadlock avoidance under concurrency
- Defining an API in Protocol Buffers and serving it as both gRPC and REST from one definition
- Interface-driven design (`Store`, `token.Maker`) enabling mock-based unit tests
- Secure auth: password hashing, PASETO tokens, refresh-token sessions, per-request authorization
- 12-factor configuration, structured logging, and request-logging middleware
- Multi-stage Docker builds and a CI pipeline that provisions a real database for tests
