package quote

import (
	"context"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"net"

	"github.com/proof-of-work/pkg/logger"
	powlib "github.com/proof-of-work/pkg/pow"
	tcputils "github.com/proof-of-work/pkg/tcp-utils"
)

type Service struct {
	l *logger.Logger

	conn net.Conn
}

func NewService(l *logger.Logger, conn net.Conn) *Service {
	return &Service{l: l, conn: conn}
}

func (s *Service) GetQuote(context.Context) (string, error) {
	s.l.Debug("Connected to server")

	line, err := tcputils.ReadMessage(s.conn)
	if err != nil {
		return "", fmt.Errorf("reading challenge: %w", err)
	}

	s.l.Debug("Received challenge", zap.ByteString("challenge", line))

	var pow powlib.POW
	if err := json.Unmarshal(line, &pow); err != nil {
		return "", fmt.Errorf("parsing challenge: %w", err)
	}

	nonce, err := pow.Solve()
	if err != nil {
		return "", fmt.Errorf("solving challenge: %w", err)
	}

	s.l.Debug("Solved challenge", zap.ByteString("nonce", nonce))

	if err := tcputils.WriteMessage(s.conn, nonce); err != nil {
		return "", fmt.Errorf("writing solve: %w", err)
	}

	res, err := tcputils.ReadMessage(s.conn)
	if err != nil {
		return "", fmt.Errorf("reading result: %w", err)
	}

	return string(res), nil
}
