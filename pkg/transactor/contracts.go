package transactor

import "context"

type Transactor interface {
	WithTx(ctx context.Context, function func(ctx context.Context) error) error
}
