# Olu

A Nigerian election voting platform that allows citizens to cast votes for political candidates via SMS or web. The system is built around voter privacy — phone numbers are never stored, only a salted SHA-256 hash — and includes duplicate vote prevention, a full audit trail, and an async SMS queue.

## Apps

| App | Description |
|---|---|
| `apps/api` | Go REST API — core voting logic, candidate management, admin panel |
| `apps/web` | Next.js frontend — voter-facing web interface |
| `apps/sms-mock` | Lightweight Go SMS mock server for local development |

## Tech Stack

**API**
- Go + Gin
- PostgreSQL (pgx/v5 with connection pooling)
- Redis (OTP storage, results cache, rate limiting)
- JWT (HMAC-signed, separate secrets for voter OTP and admin sessions)
- zerolog for structured logging

**Web**
- Next.js 16 + React 19
- TypeScript
- Tailwind CSS v4

**SMS Mock**
- Go + Gin
- In-memory store with OTP extraction and channel detection
- Browser dashboard at `http://localhost:3001`

## How Voting Works

1. Voter submits their phone number
2. An OTP is sent via SMS and stored in Redis (10-minute TTL)
3. Voter verifies the OTP and receives a short-lived JWT
4. Voter submits their candidate code (format: `A1`, `B12`, etc.)
5. The API hashes the phone number, checks for duplicates, records the vote in a transaction, and writes an audit log
6. A confirmation SMS is queued for async delivery

## API Overview

```
POST   /api/v1/auth/send-otp        — request an OTP
POST   /api/v1/auth/verify-otp      — verify OTP, receive JWT
POST   /api/v1/vote                 — cast a vote (requires JWT)
GET    /api/v1/candidates           — list active candidates
GET    /api/v1/candidates/:id       — candidate detail
GET    /api/v1/results              — live vote tally (60s cache)
GET    /health                      — health check (pings DB + Redis)

POST   /api/v1/admin/login          — admin login (rate-limited: 5 req/min)
GET    /api/v1/admin/candidates     — all candidates including inactive
POST   /api/v1/admin/candidates     — create candidate
PATCH  /api/v1/admin/candidates/:id — update candidate
DELETE /api/v1/admin/candidates/:id — deactivate candidate
GET    /api/v1/admin/stats          — vote stats breakdown
```

## Running Locally

**API**
```bash
cd apps/api
go run .
# runs on :4006
```

**SMS Mock**
```bash
cd apps/sms-mock
go run .
# runs on :3001 — dashboard at http://localhost:3001
```

**Web**
```bash
cd apps/web
npm install
npm run dev
# runs on :3000
```

## Environment Variables

| Variable | Required | Description |
|---|---|---|
| `DATABASE_URL` | PostgreSQL connection string |
| `REDIS_URL` | Redis connection string |
| `JWT_SECRET`  | Secret for voter OTP tokens |
| `ADMIN_JWT_SECRET` | Secret for admin session tokens (8h TTL) |
| `PORT` | API port (default: `4006`) |
| `ENVIRONMENT` | `development` or `production` (default: `development`) |

## Supported Parties

APC, PDP, LP, NNPP, APGA, APM, ADP, ADC, AAC, APP, Accord, AA, BP, DLA, NRM, PRP, SDP, YPP, YP, ZLP
