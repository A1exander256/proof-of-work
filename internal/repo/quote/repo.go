package quote

import "context"

type Repo interface {
	GetQuote(ctx context.Context) (string, error)
}
