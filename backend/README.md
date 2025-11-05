# 🧩 CostEasy Backend

CostEasy is a modular accounting and business management backend written in **Go (Golang)**.  
It is designed for **multi-tenant** and **per-client deployments** (e.g., UAE clinic setups)  
and focuses on accuracy, auditability, and modular expansion.

---

## 🧭 Overview

### 💼 Core Features
- 🔐 Authentication & RBAC (Roles, Permissions, MFA)
- 📘 General Ledger (GL) Core with journal entries & accounts
- ⚙️ Settings & Configuration Management
- 💰 Payroll Management
- 📊 Costing Ledger (Cost Accounting)
- 📈 Reporting & Analytics (Financial & Operational)
- 🧨 Subscription Kill Switch (Client Access Control)
- 🔑 Multi-Factor Authentication (Email / SMS / TOTP)
- ☁️ Single-instance or cloud deployment per client

---

## 🧩 System Architecture

┌──────────────────────────────┐
│ API Gateway (Go) │
│ • JWT Auth / Routing / RBAC │
│ • REST + gRPC Interface │
└──────────────┬───────────────┘
│
┌────────┼────────────────────────────────────┐
▼ ▼ ▼ ▼
┌────────┐ ┌────────────┐ ┌────────────┐ ┌──────────────┐
│ GL Core│ │ Payroll │ │ Costing │ │ Reporting │
│ (Ledger│ │ Module │ │ Ledger │ │ Engine │
│ Engine) │ │ (Salaries)│ │ (COGS) │ │ (Summaries) │
└────────┘ └────────────┘ └────────────┘ └──────────────┘
│ │ │ │
└──────────────┴──────────────┴───────────────┘
PostgreSQL + Redis + Audit Logs


---

## ⚙️ Quick Setup Guide

### 1️⃣ Requirements

| Tool | Recommended Version |
|-------|----------------------|
| Go | 1.22+ |
| PostgreSQL | 14+ |
| Redis | (optional, for MFA/OTP cache) |
| Git | Latest |

---

### 2️⃣ Clone Repository

```bash
git clone https://github.com/yourname/costeasy.git
cd costeasy/backend

3️⃣ Environment Configuration

Create your .env:

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=costing_admin
DB_PASSWORD=secret
DB_NAME=costeasy

# JWT
JWT_SECRET=my_super_secret_key
JWT_ACCESS_EXPIRY=15m
JWT_REFRESH_EXPIRY=7d

# App
HOST=0.0.0.0
PORT=8080
LOG_LEVEL=debug
ENVIRONMENT=development

4️⃣ Database Migration

The project has a built-in SQL migration tool.

Run:

go run ./app/cmd/migrate -command up -path database/migrations

Rollback:

go run ./app/cmd/migrate -command down -steps 1

5️⃣ Start API Gateway

go run ./app/cmd/api-gateway

Check:

http://localhost:8080/api/v1/health

✅ You should see:

{"status": "ok"}

📂 Project Structure

backend/
├── app/
│   ├── cmd/
│   │   ├── api-gateway/          # Main entry point
│   │   └── migrate/              # Migration CLI
│   ├── config/                   # Config loader
│   └── middleware/               # Common middlewares
│
├── internal/
│   ├── auth/                     # Auth, MFA, Roles, Permissions
│   ├── gl-core/                  # General Ledger Engine
│   ├── payroll/                  # Payroll module (Salaries, Payslips)
│   ├── costing/                  # Costing ledger module (COGS, expenses)
│   ├── reporting/                # Reporting and analytics
│   └── settings/                 # Organization, users, etc.
│
├── pkg/
│   ├── crypto/                   # Password hashing
│   ├── jwt/                      # JWT generation/validation
│   ├── migrator/                 # SQL migrations
│   ├── email/                    # SMTP utilities
│   ├── sms/                      # SMS gateway utils (Twilio / Nexmo)
│   └── cache/                    # Redis connection (for OTP/MFA)
│
└── database/
    └── migrations/               # .up.sql / .down.sql files

💰 Payroll Module

Handles all salary-related transactions.

Features:

    Employee master & contracts

    Salary computation (basic, allowances, deductions)

    Payroll journal entries to GL

    Pay slips generation

    Integration with UAE WPS format (optional)

Planned Tables:

employees
payroll_periods
salary_components
payroll_entries
payroll_journal_links

📊 Costing Ledger

Dedicated sub-ledger for cost accounting and expense tracking.

Features:

    COGS (Cost of Goods Sold)

    Department-level costing

    Allocation to projects or branches

    Journal integration with GL

Example Journal:

Dr  Cost of Sales (Expense)
Cr  Inventory (Asset)

📈 Reporting Module

Purpose:
Generate financial and operational insights.

Features:

    Trial Balance

    Profit & Loss

    Balance Sheet

    Payroll Summary Reports

    Cost Center Analysis

Reports use pre-aggregated SQL views for performance and can export to:

    CSV / Excel / PDF

    Email-based auto-scheduling

🔐 MFA & Security

Supports multiple MFA methods:

    Email-based verification (SMTP)

    SMS OTP (via Twilio / Nexmo)

    TOTP (Google Authenticator)

    Backup codes for recovery

MFA secrets and OTPs are cached in Redis for speed and expiry tracking.
☁️ Multi-Instance (UAE Deployment)

Each clinic or client runs as a single instance:

    Own PostgreSQL DB

    Optional Redis cache

    Same binary deployment

At startup:

m := migrator.NewMigrator(pool, "backend/database/migrations")
_ = m.Up(context.Background())

Each instance can update its schema independently.
🧠 Roadmap

✅ Auth + RBAC
✅ GL Core (Accounts, Journal)
✅ Custom SQL Migrator
✅ MFA (in progress)
🔲 Payroll (Q1 2026)
🔲 Costing Ledger (Q1 2026)
🔲 Reporting Engine (Q2 2026)
🔲 Redis-based Session Cache
🔲 Audit Logs & Data Archiving
🧰 Developer Commands
Task	Command
Run API	go run ./app/cmd/api-gateway
Run migrations	go run ./app/cmd/migrate -command up
Rollback migration	go run ./app/cmd/migrate -command down -steps 1
Show current migration	go run ./app/cmd/migrate -command version
Test	go test ./...
🤝 Contributing

    Fork the repo

    Create a branch (feature/payroll-module)

    Commit and push

    Open a PR

🧑‍💻 Author

Krishna Chaitanya
Vitasta Code Crafts
“Code with flow, scale with trust.”
🪪 License

MIT License © 2025 Vitasta Code Crafts