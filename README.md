# Test Bench Report Generator
 
A web-based application for hydraulic valve testing, telemetry monitoring, report generation, and valve master data management.
 
The system is designed for industrial test bench environments and integrates with OPC UA-enabled controllers such as Bosch Rexroth ctrlX CORE.
 
---
 
## Features
 
### Authentication & Authorization
 
- JWT-based authentication
- User login system
- Role-based access control
- Admin
- Operator
 
### Valve Database Management
 
- Create, update, delete, and view valve specifications
- Store:
- Manufacturer
- Part Number
- Component Series
- Valve Type
- Size (NG)
- Flow Ratings
- Pressure Ratings
- Command Parameters
- Upload valve images
- Upload PDF datasheets
 
### OPC UA Integration
 
- Connect to OPC UA servers
- Configurable OPC UA endpoint
- Support for:
- Anonymous authentication
- Username/password authentication
- Security modes:
- None
- Sign
- Sign & Encrypt
- Security policies:
- None
- Basic256Sha256
- Aes128_Sha256_RsaOaep
- Aes256_Sha256_RsaPss
 
### Test Bench Workspace
 
- Live telemetry monitoring
- Real-time charts:
- Pressure
- Flow
- Command
- Feedback
- HPU output control
- Start/Stop valve testing
- Session recording
- CSV export
 
### Test Records
 
- Save telemetry logs into database
- Historical test session storage
- Record retrieval
- Valve-linked test history
 
### Report Generator
 
Generate professional valve test reports including:
 
- Valve specifications
- Test metadata
- Parameter verification
- Pressure graphs
- Flow graphs
- Command graphs
- Feedback graphs
- Operator approval section
- Print-ready PDF output
 
### System Configuration
 
- Configure OPC UA connection settings
- Configure OPC UA node IDs
- Live node connectivity status monitoring
 
---
 
## Architecture
 
```text
┌──────────────┐
│ Frontend │
│ Vue 3 + Vite │
└──────┬───────┘
│ REST API
▼
┌──────────────┐
│ Backend API │
│ Gin + GORM │
└──────┬───────┘
│
┌─────┴─────┐
│ │
▼ ▼
SQLite OPC UA
Database Server
(ctrlX CORE)
```
 
---
 
## Technology Stack
 
### Backend
 
- Go
- Gin Framework
- GORM
- SQLite
- OPC UA Client
- JWT Authentication
 
### Frontend
 
- Vue 3
- Vite
- PrimeVue
- Pinia
- Vue Router
- Chart.js
- Axios
- Vue ChartJS
 
---
 
## Repository Structure
 
```text
backend/
├── controllers/
├── database/
├── middleware/
├── models/
├── utils/
├── images/
├── datasheets/
├── main.go
 
frontend/
├── public/
├── src/
│ ├── assets/
│ ├── components/
│ ├── layouts/
│ ├── router/
│ ├── services/
│ ├── stores/
│ └── views/
├── package.json
 
README.md
```
 
---
 
## Database Models
 
### User
 
```text
ID
Username
Password
Role
```
 
### Valve
 
```text
ID
Manufacturer
PartNumber
ComponentSeries
ValveType
SizeNG
Weight
RatedFlow
MaxFlow
CommandValue
CommandType
MaxPressure
ImagePath
DatasheetPath
```
 
### TestRecord
 
```text
ID
ValveID
CreatedAt
DataPayload
```
 
### Settings
 
```text
OPC UA Address
Node IDs
Username
Password
Security Policy
Security Mode
```
 
---
 
## API Overview
 
### Authentication
 
```http
POST /api/login
```
 
### Valves
 
```http
GET /api/valves
GET /api/valves/:id
POST /api/valves
PUT /api/valves/:id
DELETE /api/valves/:id
```
 
### Uploads
 
```http
POST /api/valves/:id/image
POST /api/valves/:id/datasheet
```
 
### Test Records
 
```http
GET /api/records
POST /api/records
```
 
### OPC UA
 
```http
GET /api/opcua/data
GET /api/opcua/status
POST /api/ctrlx/start
POST /api/ctrlx/stop
```
 
### Settings
 
```http
GET /api/settings
PUT /api/settings
```
 
---
 
## Environment Variables
 
### Backend
 
Create:
 
```bash
backend/.env
```
 
Example:
 
```env
APP_PORT=8080
 
DB_PATH=./data.db
 
JWT_SECRET=change_this_secret
 
ADMIN_PASSWORD=adminpassword
OPERATOR_PASSWORD=operatorpassword
 
OPCUA_CERT_PATH=./client_cert.pem
OPCUA_KEY_PATH=./client_key.pem
```
 
---
 
## Running the Backend
 
```bash
cd backend
 
go mod download
 
go run main.go
```
 
Server:
 
```text
http://localhost:8080
```
 
---
 
## Running the Frontend
 
```bash
cd frontend
 
npm install
 
npm run dev
```
 
Frontend:
 
```text
http://localhost:5173
```
 
---
 
## Default Users
 
If users do not exist, the application automatically seeds:
 
### Admin
 
```text
Username: admin
Password: ADMIN_PASSWORD value
```
 
### Operator
 
```text
Username: operator
Password: OPERATOR_PASSWORD value
```
 
---
 
## File Storage
 
Uploaded files are stored in:
 
```text
backend/images/
backend/datasheets/
```
 
Static access:
 
```text
/images/*
/datasheets/*
```
 
---
 
## Security Notes
 
- JWT-protected endpoints
- Role-based access control
- File type validation for uploads
- OPC UA credentials stored in Settings
- Sensitive files excluded via `.gitignore`
 
---
 
## Use Case
 
1. Login
2. Register or load valve specifications
3. Configure OPC UA settings
4. Connect to ctrlX CORE
5. Start valve test
6. Monitor telemetry
7. Record session
8. Save test record
9. Generate PDF report
10. Export documentation
 
---
 
## License
 
created by Ahmadzakir Hanif (DCEA/SVC4-AS)