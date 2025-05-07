package build

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"net"

	"github.com/proof-of-work/internal/service/quote"
)

func (b *Builder) RunTCPServer(ctx context.Context) error {
	//nolint:exhaustruct
	lc := net.ListenConfig{
		KeepAlive: b.cfg.Server.KeepAlive,
	}

	addr := fmt.Sprintf("%s:%d", b.cfg.Server.Host, b.cfg.Server.Port)

	ln, err := lc.Listen(ctx, "tcp", addr)
	if err != nil {
		return fmt.Errorf("listening on tcp %s: %w", addr, err)
	}

	b.shutdown.addWithCtxE(ln.Close)

	b.l.Info("server started", zap.String("address", addr))

	quoteFn := b.newQuoteHandleFn()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				conn, err := ln.Accept()
				if err != nil {
					if errors.Is(err, net.ErrClosed) {
						return
					}

					b.l.Error("failed to accept connection", zap.Error(err))

					continue
				}

				go func() {
					if err := quoteFn(ctx, conn); err != nil {
						b.l.Error("failed to handle quote", zap.Error(err))
					}
				}()
			}
		}
	}()

	return nil
}

type handlerFn func(ctx context.Context, conn net.Conn) error

func (b *Builder) newQuoteHandleFn() handlerFn {
	s := quote.NewService(b.l, b.cfg.Server.Deadline, b.cfg.Pow.Difficulty, b.QuoteRepo())

	return s.Handle
}
