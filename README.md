# t-challenge

This project is an API service developed in Go. It is a simple banking system that allows users to send money to each other. The project is dockerized and can be run with a single command. The project uses Postgres as a database and has a scalable structure. The project has a user authentication system and uses JWT for this. The project has a simple structure and can be developed further.

### Features
- [x] General
    - [x] Go 1.16 >
    - [x] Postgres
    - [x] Dockerized
    - [x] Scalable (nginx, load balancing)
    - [x] Database indexes
    - [x] Database transactions
    - [x] User authentication (JWT)
- [x] Routes
  - [x] User login (/login)
  - [x] User registration (/register)
  - [x] User list accounts (/my-accounts)
  - [x] Sent money (/sent)

- [x] Optionals
  - [x] Transaction logs
  - [ ] Transaction fees
  - [ ] Email/SMS notifications

### Installation

1. Clone the project:
```
git clone https://github.com/sercanarga/t-challenge.git
```
2. Go to the project directory:
```
cd t-challenge
```
3. edit the env file:
```env
SERVER_PORT=3000

# Database
DATABASE_HOST=postgres
DATABASE_NAME=postgres
DATABASE_USER=postgres
DATABASE_PASSWORD=postgres
DATABASE_PORT=5432
DB_DSN="host=${DATABASE_HOST} port=${DATABASE_PORT} user=${DATABASE_USER} password=${DATABASE_PASSWORD} dbname=${DATABASE_NAME} sslmode=disable"
```
4. Stand up the project with Docker compose:
```
docker-compose up --build -d
```

### Endpoints
[postman.com/t-challenge](http://postman.com/sercanarga/workspace/t-challenge)

### Database Diagram
![Database Diagram](https://github.com/sercanarga/t-challenge/blob/main/assets/db_scheme.png?raw=true)

### Performance & Security Tests
AI-powered [deepsource](https://deepsource.com/) was used for security tests. The results are as follows.
![test result](https://github.com/sercanarga/t-challenge/blob/main/assets/deepsource.png?raw=true)

### License
No license.