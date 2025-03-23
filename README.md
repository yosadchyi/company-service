# Company Service

A microservice for managing companies using Go with Clean Architecture.

## Features
- Create, update, delete, and fetch company records
- JWT authentication (via `golang-jwt`)
- PostgreSQL storage
- Kafka event publishing
- Swagger documentation

## Setup

### Prerequisites
- Docker
- Go 1.23+
- `swag` CLI for generating Swagger docs

### Run with Docker
```bash
make docker-up
```

### API Docs
After starting, open: [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

## Usage Example

### Authorize

```bash
curl -X 'POST' \
  'http://localhost:8080/auth/login' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "password": "password",
  "username": "admin"
}'
```

### Create Company
```bash
curl -X POST http://localhost:8080/companies \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiOTAyYjA2MWItOWZlOS00ODJkLWEyMzAtZTc2MGIwM2JiZTM5IiwiZXhwIjoxNzQyNzM0NDc5fQ.SVCnqx8r8BxCUZ7l3QK62TyaIbJtHi_mYiiWVZNs1oI" \
  -H "Content-Type: application/json" \
  -d '{"name":"AcmeCo","amount_of_employees":10,"registered":true,"type":"Corporation"}'
```

### Get Company
```bash
curl http://localhost:8080/companies/92092713-3bb5-415e-ba19-e1483359b132 \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiOTAyYjA2MWItOWZlOS00ODJkLWEyMzAtZTc2MGIwM2JiZTM5IiwiZXhwIjoxNzQyNzM0NDc5fQ.SVCnqx8r8BxCUZ7l3QK62TyaIbJtHi_mYiiWVZNs1oI" \
  -H "Content-Type: application/json" \
  -d '{"name":"AcmeCo","amount_of_employees":10,"registered":true,"type":"Corporation"}'
```

### Update Company
```bash
curl -X PATCH http://localhost:8080/companies/92092713-3bb5-415e-ba19-e1483359b132 \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiOTAyYjA2MWItOWZlOS00ODJkLWEyMzAtZTc2MGIwM2JiZTM5IiwiZXhwIjoxNzQyNzM0NDc5fQ.SVCnqx8r8BxCUZ7l3QK62TyaIbJtHi_mYiiWVZNs1oI" \
  -H "Content-Type: application/json" \
  -d '{"description":"AcmeCo company"}'
```

### Delete Company
```bash
curl -X DELETE http://localhost:8080/companies/92092713-3bb5-415e-ba19-e1483359b132 \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiOTAyYjA2MWItOWZlOS00ODJkLWEyMzAtZTc2MGIwM2JiZTM5IiwiZXhwIjoxNzQyNzM0NDc5fQ.SVCnqx8r8BxCUZ7l3QK62TyaIbJtHi_mYiiWVZNs1oI" \
  -H "Content-Type: application/json" \
```
