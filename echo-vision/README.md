# Ent commands

### Install atlas
```bash
curl -sSf https://atlasgo.sh | sh
```

### Create schema
```bash
ent new --target internal/infra/postgres/schema User
```

### Generate schema
```bash
go run internal/infra/postgres/generated/entc.go
```

### Create migrations
```bash
go run internal/infra/postgres/new.go users
```

### Update migration hash
```bash
docker compose exec echo-docs atlas migrate hash --dir file://./internal/infra/postgres/migrations/
```

# Migrations

### Migrate UP
```bash
docker compose exec echo-docs migrate -path ./internal/infra/postgres/migrations -database "postgres://postgres:postgres@db:5432/echo-docs?sslmode=disable" up
```


```bash
docker compose exec echo-docs sh -c "yes | migrate -path ./internal/infra/postgres/migrations -database 'postgres://postgres:postgres@db:5432/echo-docs?sslmode=disable' down"
```


