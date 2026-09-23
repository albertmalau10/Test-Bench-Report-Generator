# Test Bench App

Web-based application for generating, managing, and exporting valve test bench reports.

## Features

- User authentication
- Valve master data management
- Test report generation
- Datasheet management
- Image management
- System settings configuration

## Technology Stack

### Backend
- Go
- Gin Framework
- GORM
- SQLite

### Frontend
- Vue.js
- Vite
- JavaScript

## Project Structure

```text
backend/
frontend/
```

## Backend Setup

```bash
cd backend

go mod download

go run main.go
```

## Frontend Setup

```bash
cd frontend

npm install

npm run dev
```

## Environment Variables

Create file:

```text
backend/.env
```

Example:

```env
PORT=8080
```

Update the values according to your environment.

## Database

This application uses SQLite.

Database tables are created automatically through GORM AutoMigrate when the application starts.

## Notes

- `.env` is not included in the repository.
- SQLite database file is not included in the repository.
- Frontend dependencies are installed using `npm install`.