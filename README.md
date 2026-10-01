# 🚀 Valve Test Bench Report Generator

A web-based platform for **hydraulic valve testing**, **real-time telemetry monitoring**, and **automated PDF report generation**.

Built for industrial test bench environments with native **OPC UA** integration, enabling seamless communication with controllers such as **Bosch Rexroth ctrlX CORE**.

---

## ✨ Features

### 🔐 Authentication & Authorization
- JWT-based authentication
- Admin & Operator roles
- Protected API endpoints

### 🗄️ Valve Management
- Manage valve master data
- Upload valve images and datasheets
- Maintain a centralized valve database

### 🔌 OPC UA Integration
- Connect to OPC UA servers
- Configurable authentication & security policies
- Live connection status monitoring

### 📈 Real-Time Monitoring
- Live pressure monitoring
- Live flow monitoring
- Command & feedback tracking
- Interactive telemetry charts

### 🧪 Test Recording
- Record valve test sessions
- Store telemetry history
- Review historical test results

### 📄 PDF Report Generation
Generate professional valve test reports containing:

- Valve specifications
- Test metadata
- Telemetry charts
- Test results
- Operator approval section

---

## 🏗️ Architecture

```text
┌─────────────────┐
│   Vue 3 Client  │
└────────┬────────┘
         │ REST API
         ▼
┌─────────────────┐
│  Go Gin Server  │
└───────┬─────────┘
        │
   ┌────┴────┐
   ▼         ▼
 SQLite    OPC UA
Database   Server
```

---

## 🛠️ Tech Stack

### Backend
- Go
- Gin
- GORM
- SQLite
- JWT
- OPC UA Client

### Frontend
- Vue 3
- Vite
- PrimeVue
- Pinia
- Axios
- Chart.js

---

## 📁 Project Structure

```text
.
├── backend/
│   ├── controllers/
│   ├── database/
│   ├── middleware/
│   ├── models/
│   ├── utils/
│   ├── images/
│   ├── datasheets/
│   └── main.go
│
├── frontend/
│   ├── public/
│   ├── src/
│   │   ├── components/
│   │   ├── layouts/
│   │   ├── router/
│   │   ├── services/
│   │   ├── stores/
│   │   └── views/
│   └── package.json
│
└── README.md
```

---

## ⚙️ Quick Start

### 1. Backend

```bash
cd backend

go mod download

go run main.go
```

Backend runs at:

```text
http://localhost:8080
```

### 2. Frontend

```bash
cd frontend

npm install

npm run dev
```

Frontend runs at:

```text
http://localhost:5173
```

---

## 🔧 Environment Variables

Create a `.env` file inside the `backend` directory:

```env
APP_PORT=8080

DB_PATH=./data.db

JWT_SECRET=your-secret-key

ADMIN_PASSWORD=adminpassword
OPERATOR_PASSWORD=operatorpassword
```

---

## 📂 File Storage

Uploaded assets are stored in:

```text
backend/images/
backend/datasheets/
```

---

## 🔒 Security

- JWT authentication
- Role-based authorization
- Upload validation
- Configurable OPC UA security settings

---

## 📋 Typical Workflow

```text
Login
  ↓
Select Valve
  ↓
Connect to OPC UA
  ↓
Run Test
  ↓
Monitor Telemetry
  ↓
Save Test Record
  ↓
Generate PDF Report
```

---

## 🎯 Main Use Cases

- Hydraulic valve testing
- Factory Acceptance Test (FAT)
- Repair & maintenance verification
- Quality assurance documentation
- Engineering test report generation

---

## 📜 License

created by Ahmadzakir Hanif (DCEA/SVC4-AS)