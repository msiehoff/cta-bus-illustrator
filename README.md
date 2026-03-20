# CTA Bus Illustrator

An interactive map that visualizes CTA bus ridership across Chicago's route network. Line color and thickness scale with average daily ridership, making it easy to see which corridors carry the most riders and how demand compares across day types and months.

![Map showing Chicago bus routes colored by ridership](https://github.com/user-attachments/assets/93f41032-fa18-48c8-a7a3-93811a6a538e)

## Features

- **Ridership heat map** — routes colored yellow → orange → deep red by average daily ridership, with line width also scaling proportionally
- **Express/local corridors** — paired routes (e.g. Route 9 + X9 Ashland Express) are grouped and shown together with combined and individual ridership breakdowns
- **Hover tooltip** — displays route name, ID, and average ridership; shows a rank badge for top-10 routes
- **Top Routes panel** — ranked list of the 10 highest-ridership corridors; clicking a route triggers a blinking highlight on the map
- **Filters** — switch between weekday / Saturday / Sunday ridership, and select any available month from the database
- **Color legend** — gradient bar in the filter panel keeps the scale visible at all times

## Tech Stack

| Layer | Technology |
|---|---|
| Frontend | React 18, TypeScript, Vite 5 |
| Map | MapLibre GL JS, react-map-gl |
| Styling | Tailwind CSS |
| Backend | Go, Gin |
| Database | PostgreSQL, GORM |
| Migrations | goose v3 |
| Dev reloading | air |

## Project Structure

```
.
├── frontend/          # React + Vite app
│   └── src/
│       ├── components/   # RouteMap, FilterBar, RouteTooltip
│       ├── lib/          # Shared color scale constants
│       └── types/        # API response types
└── backend/
    ├── api/           # Gin HTTP handlers and DTOs
    ├── app/           # Application services and repository interfaces
    ├── business/      # Domain models (Route, RidershipRecord)
    ├── migrations/    # goose SQL migration files
    └── storage/       # PostgreSQL and fake repository implementations
```

The backend follows an Onion Architecture: the `business` domain layer has no external dependencies, `app` defines repository interfaces and orchestrates use cases, and `storage` provides the concrete implementations.

## Getting Started

### Prerequisites

- Go 1.22+
- Node.js 20+
- PostgreSQL (a local instance or Docker)
- [`air`](https://github.com/cosmtrek/air) for Go hot reloading: `go install github.com/air-verse/air@latest`

### 1. Clone the repo

```bash
git clone https://github.com/msiehoff/cta-bus-illustrator.git
cd cta-bus-illustrator
```

### 2. Configure the backend

```bash
cp backend/.env.example backend/.env
```

Edit `backend/.env` with your database URL and optional CTA API key:

```env
DATABASE_URL=postgres://postgres:postgres@localhost:5432/bus_tracker?sslmode=disable
CTA_API_KEY=your_cta_api_key_here
```

The app runs without a database — it falls back to in-memory fake data when `DATABASE_URL` is not set.

### 3. Install frontend dependencies

```bash
cd frontend && npm install
```

### 4. Run the app

From the repo root:

```bash
make serve
```

This starts both the Go backend (with hot reloading via `air`) and the Vite dev server concurrently. The app is available at [http://localhost:5173](http://localhost:5173). The frontend proxies `/api` requests to the backend on port 8080.

### Database migrations

Migrations run automatically on startup. To run them manually:

```bash
cd backend && make serve
```

## Importing Data

### Route segments

Fetches route geometry from the CTA Bus Tracker API and stores it in the database:

```
POST /api/v1/routes/import-segments
```

Requires `CTA_API_KEY` to be set in the backend environment.

### Ridership data

Upload a CSV export from the [CTA Ridership dataset](https://data.cityofchicago.org/Transportation/CTA-Ridership-Bus-Routes-Monthly-Day-Type-Averages/bynn-gwxy):

```
POST /api/v1/ridership/import
Content-Type: multipart/form-data
file: <csv file>
```

Expected columns: `route`, `routename`, `Month_Beginning`, `Avg_Weekday_Rides`, `Avg_Saturday_Rides`, `Avg_Sunday-Holiday_Rides`, `MonthTotal`.

## Planned Enhancements

See [`docs/ROADMAP.md`](docs/ROADMAP.md) for the full roadmap. Upcoming work includes:

- Speed overlay by segment (color = average speed, thickness = ridership)
- Passenger delay index — ridership × delay to surface where improvements help the most people
- Scheduled vs actual frequency analysis
- Corridor congestion analysis
