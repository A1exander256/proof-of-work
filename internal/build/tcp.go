package build

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"net"
	"sync"
	"time"

	quoteclient "github.com/proof-of-work/internal/service/client/quote"
	quoteserver "github.com/proof-of-work/internal/service/server/quote"
)

func (b *Builder) RunTCPServer(ctx context.Context) error {
	//nolint:exhaustruct
	lc := net.ListenConfig{
		KeepAlive: b.cfg.Server.KeepAlive,
	}

	addr := b.cfg.Address()

	ln, err := lc.Listen(ctx, "tcp", addr)
	if err != nil {
		return fmt.Errorf("listening on tcp %s: %w", addr, err)
	}

	b.shutdown.addWithCtxE(ln.Close)

	b.l.Info("server started", zap.String("address", addr))

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
					defer conn.Close() //nolint:errcheck

					if err := conn.SetDeadline(time.Now().Add(b.cfg.Server.Deadline)); err != nil {
						b.l.Error("failed to set deadline", zap.Error(err))

						return
					}

					fn := b.newQuoteHandleFn(conn)
					if err := fn(ctx); err != nil {
						b.l.Error("failed to handle quote", zap.Error(err))
					}
				}()
			}
		}
	}()

	return nil
}

func (b *Builder) RunTCPClient(ctx context.Context) error {
	const maxConns = 50

	var (
		wp = make(chan struct{}, maxConns)
		wg = sync.WaitGroup{}
	)

	for range b.cfg.Client.RequestCount {
		wp <- struct{}{}

		wg.Add(1)

		go func() {
			defer func() {
				<-wp
				wg.Done()
			}()

			var d net.Dialer

			conn, err := d.DialContext(ctx, "tcp", b.cfg.Address())
			if err != nil {
				b.l.Error("failed to connect to client", zap.Error(err))

				return
			}

			defer conn.Close() //nolint:errcheck

			srv := quoteclient.NewService(b.l, conn)

			res, err := srv.GetQuote(ctx)
			if err != nil {
				b.l.Error("failed to get quote", zap.Error(err))

				return
			}

			b.l.Info("got quote", zap.String("quote", res))
		}()
	}

	wg.Wait()

	return nil
}

type handlerFn func(ctx context.Context) error

func (b *Builder) newQuoteHandleFn(conn net.Conn) handlerFn {
	s := quoteserver.NewService(b.l, conn, b.cfg.Pow.Difficulty, b.QuoteRepo())

	return s.Handle
}
