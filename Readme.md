# Workshop Registration Service

A full-stack application developed for the Kenora Full Stack Developer Technical Assessment. It helps a community training centre manage workshops, staff access, attendee registrations, and cancellations without overbooking.

## Tech Stack

- **Backend:** Go, Gin, GORM
- **Database:** PostgreSQL
- **Frontend:** React, Vite, Axios
- **Security:** JWT authentication, bcrypt password hashing, role-based access control (RBAC)

## Main Features

- **Admin:** Create Manager and Staff accounts; no public signup.
- **Manager:** Create/edit workshops, view workshops, register/cancel attendees, and view registration history.
- **Staff:** View workshops, register/cancel attendees, and view registration history.
- **Workshop search:** Filter by date range, status, or available seats.
- **Registration history:** Cancelling a registration marks it as cancelled rather than deleting it; the staff IDs and timestamps are recorded.
- **Capacity protection:** PostgreSQL transactions and workshop row-level locks (`SELECT FOR UPDATE`) coordinate concurrent registration attempts.

Authorization is enforced by the backend, not just by hiding interface controls.

## Project Layout

```text
workshop-registration/
├── backend/             # Go API, configuration, services and database access
│   ├── cmd/server/
│   ├── internal/
│   └── .env.example
├── frontend/            # React application
│   └── src/
└── README.md
```

## Prerequisites

Go, Node.js/npm, and a running PostgreSQL server.

## Run Locally

**1. Create the database** (in pgAdmin or psql):

```sql
CREATE DATABASE workshop_db;
```

**2. Configure the backend.** From `backend/`, copy `.env.example` to `.env` and fill in the PostgreSQL credentials, `JWT_SECRET` (at least 32 characters), and `ADMIN_NAME`, `ADMIN_EMAIL`, `ADMIN_PASSWORD`.

PowerShell:

```powershell
cd backend
Copy-Item .env.example .env
# Edit .env with your local credentials and secrets
```

**3. Start the Go API** (from `backend/`):

```powershell
go mod tidy
go run ./cmd/server
```

Backend: `http://localhost:8080` — health check: `GET /health`.

**4. Start the frontend** in another terminal (from the project root):

```powershell
cd frontend
npm install
npm run dev
```

Frontend: `http://localhost:5173`.

If needed, set `VITE_API_URL=http://localhost:8080/api` in `frontend/.env`.

## First Login and Sample Data

On startup, the backend seeds the initial Admin account from `backend/.env` if it does not already exist. Log in using the configured `ADMIN_EMAIL` and `ADMIN_PASSWORD`.

1. As **Admin**, create one Manager and one Staff account.
2. As **Manager**, create a workshop, e.g. `WS001` — *Introduction to Go*, capacity `20`, status `scheduled`.
3. As **Staff**, add attendees, view registrations, and cancel one to verify that a seat is released.

**Note:** The sample workshop above is created manually; automatic sample-workshop seeding is not included in the implementation described here.

## Main API Endpoints

| Method | Endpoint | Role |
|---|---|---|
| POST | `/api/auth/login` | Public |
| GET | `/api/auth/me` | Authenticated |
| POST | `/api/users` | Admin |
| GET | `/api/workshops` | Manager, Staff |
| GET | `/api/workshops/:id` | Manager, Staff |
| POST | `/api/workshops` | Manager |
| PUT | `/api/workshops/:id` | Manager |
| POST | `/api/workshops/:id/registrations` | Manager, Staff |
| GET | `/api/workshops/:id/registrations` | Manager, Staff |
| POST | `/api/registrations/:id/cancel` | Manager, Staff |

Example filtering: `GET /api/workshops?status=scheduled&availableOnly=true`.

Protected endpoints require `Authorization: Bearer <JWT_TOKEN>`.

## Important Design Notes

- Active registrations must never exceed workshop capacity, including when requests arrive simultaneously.
- The registration service checks capacity inside a database transaction after locking the workshop row. Workshop capacity edits use the same lock.
- Cancelled registrations remain available in history, including who cancelled them and when.
- Workshop codes are unique; workshop statuses are `scheduled`, `completed`, and `cancelled`.
- Workshop deletion, waitlists, and broader audit logging are outside the initial scope.

## Before Submission

Verify login and role restrictions, workshop CRUD, registration/cancellation, seat counts, filtering, and **simultaneous last-seat requests**. Include the requested one-page technical design document. Tests have not been independently verified in this README.

**Security:** Never commit `backend/.env`, real passwords, JWT secrets, or `node_modules/` to Git.
