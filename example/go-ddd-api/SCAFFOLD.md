---
Params:
- name
---

# domain/{{.Params.name}}/service.go

```golang
package {{.Params.name}}

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/igtm/example-go-api/server/gen/{{.Params.name}}"
)
{{ $Name := .Params.name | ToCamel}}
{{ $name := .Params.name }}

// Service
type Service interface {
	Create{{$Name}}(ctx context.Context, body *{{$name}}.V1{{$Name}}) error
	Update{{$Name}}(ctx context.Context, id string, body *{{$name}}.V1{{$Name}}) error
	Delete{{$Name}}(ctx context.Context, id string) error
	Get{{$Name}}(ctx context.Context, id string) (*{{$name}}.V1{{$Name}}, error)
	List{{$Name}}(ctx context.Context, params SearchParams) ({{$name}}.V1{{$Name}}Collection, error)
	GetTotalNum(ctx context.Context, params SearchParams) (uint64, error)
}

type service struct {
	repo Repository
}

// NewService
func NewService(repo Repository) Service {
	return &service{
		repo,
	}
}

func (s *service) Create{{$Name}}(ctx context.Context, body *{{$name}}.V1{{$Name}}) error {
	return s.repo.WithDBTransaction(ctx, func(tx *sqlx.Tx) error {
		return s.repo.Create{{$Name}}(tx, body)
	})
}
func (s *service) Update{{$Name}}(ctx context.Context, id string, body *{{$name}}.V1{{$Name}}) error {
	return s.repo.WithDBTransaction(ctx, func(tx *sqlx.Tx) error {
		return s.repo.Update{{$Name}}(tx, id, body)
	})
}

func (s *service) Delete{{$Name}}(ctx context.Context, id string) error {
	return s.repo.WithDBTransaction(ctx, func(tx *sqlx.Tx) error {
		return s.repo.Delete{{$Name}}(tx, id)
	})
}

func (s *service) Get{{$Name}}(ctx context.Context, id string) (*{{$name}}.V1{{$Name}}, error) {
	return s.repo.Get{{$Name}}(ctx, id)
}

func (s *service) List{{$Name}}(ctx context.Context, params SearchParams) ({{$name}}.V1{{$Name}}Collection, error) {
	return s.repo.List{{$Name}}(ctx, params)
}

func (s *service) GetTotalNum(ctx context.Context, params SearchParams) (uint64, error) {
	return s.repo.GetTotalNum(ctx, params)
}

```


# domain/{{.Params.name}}/repository.go

```golang
package {{.Params.name}}

import (
    "context"
    "errors"
    
    "github.com/jmoiron/sqlx"
    "github.com/igtm/example-go-api/server/gen/{{.Params.name}}"
)
{{ $Name := .Params.name | ToCamel}}
{{ $name := .Params.name }}

var (
	// ErrQueryExecution failed to execute query
	ErrQueryExecution = errors.New("failed to execute query")
)

// Repository
type Repository interface {
	WithDBTransaction(ctx context.Context, txFunc func(*sqlx.Tx) error) (err error)
	List{{$Name}}(ctx context.Context, params SearchParams) ({{$name}}.V1{{$Name}}Collection, error)
    Get{{$Name}}(ctx context.Context, id string) (*{{$name}}.V1{{$Name}}, error)
	GetTotalNum(ctx context.Context, params SearchParams) (uint64, error)
	Create{{$Name}}(tx *sqlx.Tx, body *{{$name}}.V1{{$Name}}) error
	Update{{$Name}}(tx *sqlx.Tx, id string, body *{{$name}}.V1{{$Name}}) error
	Delete{{$Name}}(tx *sqlx.Tx, id string) error
}

```

# domain/{{.Params.name}}/types.go

```golang
package {{.Params.name}}

import (
    "errors"
)

type SearchParams struct {
	Limit          *uint64
	Offset         *uint64
	Q              *string
}

var (
	ErrNotFound = errors.New("Record was not found.")
)

```

# infra/mysql/{{.Params.name}}.go

```golang
package mysql
{{ $Name := .Params.name | ToCamel}}
{{ $name := .Params.name }}

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	domain{{$Name}} "github.com/igtm/example-go-api/server/domain/{{$name}}"
	"github.com/orario-developer/haitatsu-web/server/gen/{{$name}}"
)

type {{$name}}Repository struct {
	{{$name}}DB *sqlx.DB
}

// New{{$Name}}Repository
func New{{$Name}}Repository({{$name}}DB *sqlx.DB) domain{{$Name}}.Repository {
	return &{{$name}}Repository{
		{{$name}}DB: {{$name}}DB,
	}
}

// WithTransaction
func (repo *{{$name}}Repository) WithDBTransaction(ctx context.Context, txFunc func(*sqlx.Tx) error) (err error) {
	tx, err := repo.{{$name}}DB.BeginTxx(ctx, nil)
	if err != nil {
		return
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p) // re-throw panic after Rollback
		} else if err != nil {
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()
	err = txFunc(tx)
	return err
}

// Get{{$Name}} ... 
func (repo *{{$name}}Repository) Get{{$Name}}(ctx context.Context, id uint64) (*{{$name}}.V1{{$Name}}, error) {
	query := sq.Select(`*`).
		From("{{$name}}").
		Where(sq.Eq{"{{$name}}.id": id})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}
	var list []{{$name}}.V1{{$Name}}
	err = repo.{{$name}}DB.SelectContext(ctx, &list, sql, args...)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, fmt.Errorf("{{$name}}Repository Get{{$Name}} error: %w", domain{{$Name}}.ErrNotFound)
	}
	return list[0], nil
}

// List{{$Name}}
func (repo *{{$name}}Repository) List{{$Name}}(ctx context.Context, params domain{{$Name}}.SearchParams) ([]*{{$name}}.V1{{$Name}}, error) {
	query := sq.Select(`*`).From("{{$name}}")

	var limit uint64 = 10
	if params.Limit != nil {
		limit = *params.Limit
	}
	query = query.Limit(limit)

	var offset uint64 = 0
	if params.Offset != nil {
		offset = *params.Offset
	}
	query = query.Offset(offset)

	// Q
	if params.Q != nil {
		query = query.Where("{{$name}}.name LIKE ?", fmt.Sprint("%", *params.Q, "%"))
	}
	sql, args, err := query.ToSql()
	if err != nil {
		return []*{{$name}}.V1{{$Name}}{}, err
	}

	var list []{{$name}}.V1{{$Name}}
	err = repo.{{$name}}DB.SelectContext(ctx, &list, sql, args...)
	return list, err
}

// GetTotalNum
func (repo *{{$name}}Repository) GetTotalNum(ctx context.Context, params domain{{$Name}}.SearchParams) (uint64, error) {
	query := sq.Select(`COUNT(1)`).
		From("{{$name}}")

	// Q
	if params.Q != nil {
		query = query.Where("{{$name}}.name LIKE ?", fmt.Sprint("%", *params.Q, "%"))
	}

	sql, args, err := query.ToSql()

	if err != nil {
		return 0, err
	}

	var cnt uint64
	err = repo.{{$name}}DB.GetContext(ctx, &cnt, sql, args...)
	return cnt, err
}

// Create{{$Name}}
func (repo *{{$name}}Repository) Create{{$Name}}(tx *sqlx.Tx, payload {{$name}}.{{$Name}}Payload) (uint64, error) {
	setMap := map[string]interface{}{
		"name":                                payload.Name,
	}
	
	result, err := sq.Insert("{{$name}}").SetMap(setMap).RunWith(tx).Exec()
	if err != nil {
		return 0, err
	}
	intID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	uint64ID := uint64(intID)
	return uint64ID, nil
}

// Update{{$Name}}
func (repo *{{$name}}Repository) Update{{$Name}}(tx *sqlx.Tx, payload {{$name}}.Update{{$Name}}Payload) error {
	setMap := map[string]interface{}{
		"name":                                payload.Name,
	}

	_, err := sq.Update("{{$name}}").
		SetMap(setMap).
		Where(sq.Eq{"id": payload.ID}).
		RunWith(tx).
		Exec()

	return err
}

// Delete{{$Name}}
func (repo *{{$name}}Repository) Delete{{$Name}}(tx *sqlx.Tx, id uint64) error {
	_, err := sq.Delete("{{$name}}").
		Where(sq.Eq{"id": id}).
		RunWith(tx).Exec()

	return err
}


```