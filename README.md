# **PennyPilot API**

The **PennyPilot API** is a backend service for a personal budgeting and saving app that helps users allocate their income into sub-accounts, set financial goals, and track expenses, savings, and investments. The API is built using **Go**, with **Chi router** for handling routes and **PostgreSQL** as the database.

---

## **Features**

- User authentication (sign up, log in, password reset)
- Income and sub-account management
- Dynamic percentage-based allocation of funds into sub-accounts
- Tracking of transfers between sub-accounts
- Lending and borrowing management
- Goal setting for savings and investments

---

## **Tech Stack**

- **Language**: Go
- **Framework**: Chi router
- **Database**: PostgreSQL
- **Architecture**: RESTful API
- **Authentication**: JWT-based authentication

---

## **Setup and Installation**

1. **Clone the repository**
   ```bash
   git clone https://github.com/felix-Asante/pennyPilot-go-api.git
   cd pennypilot-api
   ```
   2. **Install dependencies**
   ```bash
   go mod tidy
   ```
2. **Setup DB**
   Run this code to setup uuid for the database
   ```bash
   CREATE EXTENSION IF NOT EXISTS "uuid-ossp"
   ```
   4. **Run the server**
   ```bash
   go run main.go
   or
   air
   ```
