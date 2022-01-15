# Go DDD API

This is example code for generating DDD scaffold files.

Too many boilerplate is the main pain of DDD(Layered Architecture) in Golang.

The [scaffold](./SCAFFOLD.md) creates domain layer and infrastructure(mysql) code from given name parameter.

[sqlx](github.com/jmoiron/sqlx) and [squirrel](github.com/Masterminds/squirrel) are used for mysql.

1. Run `manaita -p name=company`
2. `domain/service.go`, `domain/repository.go`, `domain/types.go` and `infra/company.go` files are generated.
