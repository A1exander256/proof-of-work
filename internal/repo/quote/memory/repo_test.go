package memory

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_repo_GetQuote(t *testing.T) {
	r := NewRepo()

	for range 10 {
		quoteText, err := r.GetQuote(t.Context())

		require.NoError(t, err)
		require.NotEmpty(t, quoteText)
		require.Contains(t, quotes, quoteText)
	}
}
