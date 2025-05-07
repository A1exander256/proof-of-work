package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	t.Setenv("LOG_LEVEL", "info")
	t.Setenv("SERVER_HOST", "localhost")
	t.Setenv("SERVER_PORT", "12345")
	t.Setenv("SERVER_KEEP_ALIVE", "2s")
	t.Setenv("SERVER_DEADLINE", "5s")
	t.Setenv("POW_DIFFICULTY", "15")
	t.Setenv("CLIENT_REQUEST_COUNT", "3")

	cfg, err := Parse()
	require.NoError(t, err)

	require.Equal(t, "info", cfg.App.LogLevel)
	require.Equal(t, "localhost", cfg.Server.Host)
	require.Equal(t, 12345, cfg.Server.Port)
	require.Equal(t, 2*time.Second, cfg.Server.KeepAlive)
	require.Equal(t, 5*time.Second, cfg.Server.Deadline)
	require.Equal(t, uint8(15), cfg.Pow.Difficulty)
	require.Equal(t, uint8(3), cfg.Client.RequestCount)
}
