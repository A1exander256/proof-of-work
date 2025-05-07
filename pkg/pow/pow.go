package pow

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"math"
	"math/bits"
	"strconv"
)

var errSolveChallenge = errors.New("failed to solve challenge")

type POW struct {
	Seed       string `json:"seed,omitempty"`
	Difficulty uint8  `json:"difficulty,omitempty"`
}

func NewPOW(difficulty uint8) POW {
	return POW{
		Seed:       uuid.NewString(),
		Difficulty: difficulty,
	}
}

func (p POW) Challenge() ([]byte, error) {
	chBytes, err := json.Marshal(p)
	if err != nil {
		return nil, fmt.Errorf("marshaling challenge: %w", err)
	}

	return chBytes, nil
}

func (p POW) Solve() ([]byte, error) {
	for i := range uint64(math.MaxUint64) {
		nonce := strconv.FormatUint(i, 10)
		hash := sha256.Sum256([]byte(p.Seed + nonce))

		if leadingZeroBits(hash) >= int(p.Difficulty) {
			return []byte(nonce), nil
		}
	}

	return nil, errSolveChallenge
}

func (p POW) Verify(nonce []byte) bool {
	hash := sha256.Sum256([]byte(p.Seed + string(nonce)))

	return leadingZeroBits(hash) >= int(p.Difficulty)
}

func leadingZeroBits(data [32]byte) int {
	count := 0

	for _, b := range data {
		if b != 0 {
			count += bits.LeadingZeros8(b)

			break
		}

		count += 8
	}

	return count
}
