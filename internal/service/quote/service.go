package quote

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"net"
	"time"

	"github.com/proof-of-work/pkg/logger"
	powlib "github.com/proof-of-work/pkg/pow"
	tcputils "github.com/proof-of-work/pkg/tcp-utils"
)

type repo interface {
	GetQuote(ctx context.Context) (string, error)
}

type Service struct {
	l *logger.Logger

	deadline   time.Duration
	difficulty uint8

	r repo
}

func NewService(
	l *logger.Logger,
	deadline time.Duration,
	difficulty uint8,
	r repo,
) *Service {
	return &Service{
		l:          l,
		deadline:   deadline,
		difficulty: difficulty,
		r:          r,
	}
}

func (s *Service) Handle(ctx context.Context, conn net.Conn) error {
	defer conn.Close() //nolint:errcheck

	if err := conn.SetDeadline(time.Now().Add(s.deadline)); err != nil {
		return fmt.Errorf("setting deadline: %w", err)
	}

	s.l.Debug("New connection", zap.String("remote", conn.RemoteAddr().String()))

	pow := powlib.NewPOW(s.difficulty)

	challenge, err := pow.Challenge()
	if err != nil {
		return fmt.Errorf("cteating challenge: %w", err)
	}

	if err := tcputils.WriteMessage(conn, challenge); err != nil {
		return fmt.Errorf("writing challenge: %w", err)
	}

	solution, err := tcputils.ReadMessage(conn)
	if err != nil {
		return fmt.Errorf("reading solution: %w", err)
	}

	if !pow.Verify(solution) {
		msg := []byte("Invalid proof of work")

		if err := tcputils.WriteMessage(conn, msg); err != nil {
			return fmt.Errorf("writing error: %w", err)
		}
	}

	s.l.Debug("Solution verified", zap.String("remote", conn.RemoteAddr().String()))

	quote, err := s.r.GetQuote(ctx)
	if err != nil {
		return fmt.Errorf("getting quote from repo: %w", err)
	}

	if err := tcputils.WriteMessage(conn, []byte(quote)); err != nil {
		return fmt.Errorf("writing quote: %w", err)
	}

	return nil
}
