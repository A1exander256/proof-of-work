package memory

import (
	"context"
	"crypto/rand"
	"math/big"

	"github.com/proof-of-work/internal/repo/quote"
)

var quotes = []string{
	"Fruits and wholesome herbs, including vegetables, which should be used with prudence and thanksgiving.",
	"The flesh “of beasts and of the fowls of the air, which is “to be used sparingly.",
	"Grains such as wheat, rice, and oats, which are “the staff of life.",
}

type repo struct{}

func NewRepo() quote.Repo {
	return &repo{}
}

func (*repo) GetQuote(context.Context) (string, error) {
	nBig, err := rand.Int(rand.Reader, big.NewInt(int64(len(quotes))))
	if err != nil {
		return "", err
	}

	return quotes[nBig.Int64()], nil
}
