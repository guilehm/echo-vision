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
docker compose exec echo-vision go run internal/infra/postgres/generated/entc.go
```

### Create migrations
```bash
docker compose exec echo-vision go run internal/infra/postgres/new.go users
```

### Update migration hash
```bash
atlas migrate hash --dir file://./internal/infra/postgres/migrations
```

# Migrations

### Migrate UP
```bash
docker compose exec echo-vision migrate -path ./internal/infra/postgres/migrations -database "postgres://postgres:postgres@db:5432/echo-vision?sslmode=disable" up
```


```bash
docker compose exec echo-vision sh -c "yes | migrate -path ./internal/infra/postgres/migrations -database 'postgres://postgres:postgres@db:5432/echo-vision?sslmode=disable' down"
```

# Tests
```bash
go test -v ./...
```
