package pow

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPOW(t *testing.T) {
	tests := []struct {
		name       string
		difficulty uint8
		expectErr  bool
	}{
		{"low difficulty", 10, false},
		{"medium difficulty", 16, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewPOW(tt.difficulty)

			challenge, err := p.Challenge()
			require.NoError(t, err)

			var parsed POW
			err = json.Unmarshal(challenge, &parsed)
			require.NoError(t, err)

			require.Equal(t, p.Seed, parsed.Seed)
			require.Equal(t, p.Difficulty, parsed.Difficulty)

			nonce, err := p.Solve()
			if tt.expectErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
			require.NotEmpty(t, nonce)
			require.True(t, p.Verify(nonce), "POW verification failed")
		})
	}
}
