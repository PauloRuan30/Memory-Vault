# Memory Vault - PS2 Nostalgic File Manager

A web-based file manager that completely reimagines the user experience as the iconic PlayStation 2 Memory Card interface. Combines 2000s gaming nostalgia with modern, high-performance, asynchronous microservice architecture.

## Architecture Overview

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Frontend      │    │  API Gateway    │    │  AI Worker      │
│   React + R3F   │◄──►│     Go          │◄──►│   Python        │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│ PostgreSQL DB   │    │    Redis        │    │   S3 Storage    │
│ (Metadata)      │    │ (Job Queue)     │    │ (Files/Icons)   │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## 🚀 Quick Start

### Prerequisites
- Docker & Docker Compose
- Node.js 18+ (for local frontend development)
- Go 1.22+ (for local backend development)

### Local Development Setup

1. **Clone the repository:**
```bash
git clone https://github.com/PauloRuan30/Memory-Vault.git
cd memory-vault
```

2. **Set up environment variables:**
Create `.env` files in each service directory with appropriate configuration (see `.env.example` templates).

3. **Start all services with Docker Compose:**
```bash
docker-compose up --build
```

4. **Access the application:**
- Frontend: `http://localhost:3000`
- API Gateway: `http://localhost:8080`
- PostgreSQL: `localhost:5432`
- Redis: `localhost:6379`

### Manual Setup (Alternative)

**Frontend:**
```bash
cd frontend
npm install
npm run dev
```

**API Gateway:**
```bash
cd services/api-gateway
go mod tidy
go run cmd/server/main.go
```

**Icon Worker:**
```bash
cd services/icon-worker
pip install -r requirements.txt
python main.py
```

## 🗄️ Database Schema

```sql
CREATE TABLE "users" (
  "id" SERIAL PRIMARY KEY,
  "username" VARCHAR(50) UNIQUE NOT NULL,
  "password_hash" VARCHAR(255) NOT NULL,
  "created_at" TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE "files" (
  "id" SERIAL PRIMARY KEY,
  "user_id" INT NOT NULL REFERENCES "users"("id") ON DELETE CASCADE,
  "parent_folder_id" INT REFERENCES "files"("id") ON DELETE CASCADE,
  "is_folder" BOOLEAN NOT NULL DEFAULT false,
  "name" VARCHAR(255) NOT NULL,
  "s3_path" VARCHAR(1024),
  "texture_path" VARCHAR(1024),
  "size_kb" INT NOT NULL DEFAULT 0,
  "processing_status" VARCHAR(20) DEFAULT 'PENDING',
  "metadata" JSONB,
  "created_at" TIMESTAMPTZ DEFAULT NOW()
);
```

## 🎯 Core Features

1. **PS2-Style Navigation**: 4x4 grid with stepped cursor movement
2. **Real-time AI Processing**: Asynchronous texture generation for file icons
3. **Auth System**: JWT-based authentication and authorization
4. **File Management**: CRUD operations for files and folders
5. **CRT Effects**: Authentic PS2 visual experience with shaders
6. **Retro Audio**: Authentic UI interaction sounds

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

*Built with ❤️ for PS2 nostalgia and modern full-stack architecture*