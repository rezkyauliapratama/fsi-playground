# FSI Playground

Experiments for financial-services architectures: IAM/RBAC/ABAC with Ory Keto, transaction handling, and load-test tooling. Go workspace.

## Structure

```
libs/
├── error/          # shared error helpers
└── helper/         # transaction helpers
services/
├── intenal-iam-service/   # IAM service: RBAC vs ABAC (Ory Keto), seed + service + handlers
└── loadtest/               # load testing scripts
```

## IAM service (Ory Keto)

Demonstrates relationship-based access control (ReBAC):

- `doc/rbac_vs_abac.md` — design notes comparing RBAC vs ABAC
- Ory Keto via docker-compose (`docker/keto.yml`)
- Seed command populates `relations.json`

```bash
cd services/intenal-iam-service
docker compose up -d
go run cmd/seed/main.go
go run cmd/service/main.go
```

## Status

Playground — active experiments for banking/fintech security patterns.
