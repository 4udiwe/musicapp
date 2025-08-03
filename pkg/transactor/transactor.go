package transactor

import "context"

type Transactor interface {
	WithinTransaction(ctx context.Context, fn func(context.Context) error) error
}
