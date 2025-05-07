package build

import (
	"github.com/proof-of-work/internal/repo/quote"
	quetememory "github.com/proof-of-work/internal/repo/quote/memory"
)

func (*Builder) QuoteRepo() quote.Repo {
	return quetememory.NewRepo()
}
