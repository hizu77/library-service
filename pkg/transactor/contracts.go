package transactor

import "context"

//go:generate mockgen -source=contracts.go -destination=../../internal/usecase/mock/transactor.go -package=mock

type Transactor interface {
	WithTx(ctx context.Context, function func(ctx context.Context) error) error
}
