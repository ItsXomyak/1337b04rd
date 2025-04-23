# 1337b04rd 📂

![Go](https://img.shields.io/badge/Go-1.23-blue) ![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-blue) ![Docker](https://img.shields.io/badge/Docker-✓-blue)

An anonymous image board backend written in **Go**. Supports threads, comments, image uploads to S3-compatible storage (MinIO), and session tracking via cookies. Built using **Hexagonal Architecture**, follows clean coding principles, and includes a frontend with Go templates.

## ✨ Features

- 🛸 Anonymous sessions with Rick & Morty avatars
- 📍 Thread and comment posting
- 📷 Image upload support (MinIO / S3-compatible)
- ⌛ Auto-cleanup of threads without comments
- ⚖️ Moderation-ready architecture
- ✨ Clean Go codebase with layered separation


## 📚 Tech Stack

- **Go** 1.23+
- **PostgreSQL**
- **MinIO** (S3-compatible object storage)
- **HTML templates** (server-rendered)
- **TailwindCSS** (frontend styling)
- **Docker** (virtualization)

## 🧬 Project Structure

```
1337b04rd/
├── cmd/                     # CLI entrypoint
├── config/                  # Manual .env parsing & configuration
├── db/                      # SQL init scripts
├── internal/                # Core business logic and adapters
│   ├── adapters/            # Infrastructure adapters
│   │   ├── http/            # HTTP handlers & middleware
│   │   ├── postgres/        # PostgreSQL repositories
│   │   ├── rickmorty/       # External avatar client
│   │   └── s3/              # S3/MinIO storage client
│   ├── app/
│   │   ├── common/
│   │   │   ├── logger/      # Structured logging
│   │   │   └── utils/       # Helpers (UUID, etc)
│   │   ├── ports/           # Interfaces for domain-driven services
│   │   └── services/        # Business logic
│   └── domain/              # Core domain models and rules
│       ├── avatar/
│       ├── comment/
│       ├── errors/
│       ├── session/
│       └── thread/
├── test/                    # Tests and testdata
│   ├── integration/
│   ├── testdata/
│   └── unit/
└── web/                     # Frontend
    ├── static/              # Static assets (img)
    │   └── img/
    └── templates/           # HTML templates
```

## ╰┈➤ Getting Started (Local)

### 1. Clone & Init

```bash
git clone git@git.platform.alem.school:mromanul/1337b04rd.git
cd 1337b04rd
```

### 2. Setup `.env`

Create `.env` based on `.env.example`:

```env
# HTTP server
PORT=8080

# PostgreSQL
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_NAME=1337b04rd
DB_SSLMODE=disable

# S3-compatible storage (MinIO or triple-s)
S3_ENDPOINT=minio:9000
S3_ACCESS_KEY=your_s3_key
S3_SECRET_KEY=your_s3_secret
S3_BUCKET_THREADS=1337-threads
S3_BUCKET_COMMENTS=1337-comments
S3_REGION=us-east-1
S3_USE_SSL=false

# Rick and Morty API
AVATAR_API_BASE_URL=https://rickandmortyapi.com/api

# Session settings
SESSION_COOKIE_NAME=1337session
SESSION_DURATION_DAYS=7

# App mode (for logging, etc.)
APP_ENV=development
```

### 3. Run MinIO

Use:
```bash
docker-compose up --build
```

## 🎨 Frontend (Python)

Open with VSCode Live Server or Python:

```bash
cd web/templates
python -m SimpleHTTPServer 'port'
```

Then open:
```
http://localhost:'port'/catalog.html or what do you want to open
```

## 📑 Tests

```bash
go test ./...
```


## 👨🏻‍💻 Authors

- [![Status](https://img.shields.io/badge/alem-azhaxyly-success?logo=github)](https://platform.alem.school/git/azhaxyly) <a href="https://t.me/hmlssdeus" target="_blank"><img src="https://img.shields.io/badge/telegram-@hmlssdeus-blue?logo=Telegram" alt="Status" /></a>
- [![Status](https://img.shields.io/badge/alem-mromanul-success?logo=github)](https://platform.alem.school/git/mromanul)  <a href="https://t.me/scrameee" target="_blank"><img src="https://img.shields.io/badge/telegram-@scrameee-blue?logo=Telegram" alt="Status" /></a>

## 🎉 Acknowledgements <a name = "acknowledgement"></a>

This project has been created by:

- Zhaxylykov A, ***"kazakh painter"***
- Romanuly M, ***"FullStack overflow"***

## 📜 License

Apache License Version 2.0

