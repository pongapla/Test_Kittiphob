# Test_Kittiphob — Queue Management System

**Full Stack Developer Technical Assessment**  
**Candidate:** Kittiphob  
**Project:** Queue Management System (IT 05-1, IT 05-2, IT 05-3)

## Overview

A queue-ticket web application based on the provided assessment screens. Users can request a queue ticket, view its number and issue date/time, and reset the current queue.

## Tech Stack

| Layer | Technology |
| --- | --- |
| Frontend | Vue 3, TypeScript, Vite, Vue Router |
| Backend | Go, Chi v5 |
| Database | SQLite (`modernc.org/sqlite`) |
| Testing | Go tests, Go race detector, Playwright |

## Features

- **IT 05-1 — Queue reception:** Request a new ticket or navigate to queue reset.
- **IT 05-2 — Ticket display:** Show issued queue number, date, and time; return to reception. Reloading this page does not issue a new ticket.
- **IT 05-3 — Queue reset:** Show the current queue number and reset the sequence. The display becomes `00` after a successful reset.
- Sequential numbering from `A0` through `Z9` (260 tickets per cycle). After `Z9`, no new tickets are issued until a reset.
- Ticket identifiers and cycle numbers distinguish tickets issued before and after a reset.

## Architecture

```text
Browser
  |
Vue 3 + TypeScript (Vite)
  | REST / JSON (/api via Vite proxy)
Go + Chi Router
  |
HTTP Handler
  |
SQLite Repository (transactions)
  |
SQLite Database
```

The backend assigns ticket numbers. SQLite transactions and a unique `(cycle, queue_number)` constraint help maintain consistent ticket issuance. The current implementation serializes database access within its connection pool; multi-process deployment is not covered by the current tests.

## Project Structure

```text
Test_Kittiphob/
├── backend/
│   ├── cmd/api/main.go
│   ├── internal/handler/
│   ├── internal/repository/
│   │   └── schema.sql
│   ├── internal/service/
│   └── data/                  # Local SQLite file (not committed)
├── frontend/
│   ├── src/router/
│   ├── src/services/
│   ├── src/views/
│   └── tests/e2e/
├── docs/
│   ├── REQUIREMENTS.md
│   └── TEST_CASES.md
└── README.md
```

## Prerequisites

- Go (use a version compatible with `backend/go.mod`)
- Node.js and npm (use versions compatible with `frontend/package.json` / Vite)
- Git

No standalone SQLite server is required.

## Run Locally

### 1. Clone

```bash
git clone https://github.com/pongapla/Test_Kittiphob.git
cd Test_Kittiphob
```

### 2. Start the backend

In terminal 1:

```bash
cd backend
go mod download
go run ./cmd/api
```

The API listens at **http://localhost:4000**. The SQLite database is initialized in `backend/data/` at startup.

Check its health:

```bash
curl http://localhost:4000/api/health
```

Expected response:

```json
{"status":"ok"}
```

### 3. Start the frontend

In terminal 2, from the project root:

```bash
cd frontend
npm ci
npm run dev
```

Open **http://localhost:5173/**. Vite proxies `/api` to the backend at `http://localhost:4000`.

## API Endpoints

| Method | Endpoint | Purpose |
| --- | --- | --- |
| `GET` | `/api/health` | Health check |
| `POST` | `/api/queues/` | Issue a new ticket |
| `GET` | `/api/queues/current` | Get current queue number (`00` after reset) |
| `GET` | `/api/queues/{id}` | Retrieve an issued ticket |
| `POST` | `/api/queues/reset` | Reset sequence and start a new cycle |

Example ticket request:

```bash
curl -X POST http://localhost:4000/api/queues/
```

An example response (IDs and timestamps vary):

```json
{
  "id": "8ed3a1369bc03fc4b4367f9be9707406",
  "cycle": 1,
  "queueNumber": "A0",
  "issuedAt": "2026-10-14T08:53:21Z"
}
```

Reset queue:

```bash
curl -X POST http://localhost:4000/api/queues/reset
```

## Testing

### Backend tests

From `backend/`:

```bash
go test ./... -v
go test -race ./...
```

Tests cover queue-number formatting, valid ranges, uniqueness across the 260 possible numbers, issuance/reset, capacity, and 100 concurrent goroutines issuing tickets.

### Browser end-to-end test

Start the Go API before running the browser test. From `frontend/`:

```bash
npm ci
npx playwright install chromium
npx playwright test
```

The queue E2E test covers receiving consecutive tickets, ticket-page refresh, checking the current number, resetting the queue, and restarting at `A0`.

**Important:** Browser tests reset and modify the SQLite queue used by the running backend. Run them only against a disposable development/test database, not against real usage data.

### Frontend production build

From `frontend/`:

```bash
npm run build
```

## Design Notes and Assumptions

- Queue numbering is `A0, A1, ..., A9, B0, ..., Z9`.
- Reset starts a new cycle, allowing `A0` to be used again in that new cycle.
- Previously issued tickets can be retrieved by their unique IDs.
- No login or authorization is implemented for the reset action; it is outside the assignment's stated scope.
- A real production deployment would require additional considerations such as access control for resets, deployment configuration, and broader concurrency/restart testing.
- The application contains **no voice-announcement feature**.

## Supporting Documents

- [Requirements](docs/REQUIREMENTS.md)
- [Test cases](docs/TEST_CASES.md)

---

Prepared as a Full Stack Developer technical assessment submission.
