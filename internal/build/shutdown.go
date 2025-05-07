package build

import (
	"context"
	"os"
	"os/signal"
	"slices"
	"syscall"

	"github.com/proof-of-work/pkg/logger"
)

func (b *Builder) WaitShutdown(ctx context.Context) {
	stopSignals := []os.Signal{syscall.SIGKILL, syscall.SIGINT, syscall.SIGTERM}
	s := make(chan os.Signal, len(stopSignals))

	signal.Notify(s, stopSignals...)

	b.l.Sugar().Infof("got %s os signal. application will be stopped", <-s)

	b.shutdown.do(ctx)
}

func (b *Builder) ShutdownChannel(ctx context.Context) chan struct{} {
	stop := make(chan struct{})

	go func() {
		stopSignals := []os.Signal{syscall.SIGTERM, syscall.SIGINT}
		s := make(chan os.Signal, len(stopSignals))

		signal.Notify(s, stopSignals...)

		b.l.Sugar().Infof("got %s os signal. application will be stopped", <-s)

		b.shutdown.do(ctx)

		close(stop)
	}()

	return stop
}

type shutdownFn func(context.Context) error

type shutdown struct {
	fn []shutdownFn
}

func (s *shutdown) addWithCtxE(f func() error) {
	s.fn = append(s.fn, withCtxE(f))
}

func withCtxE(f func() error) func(context.Context) error {
	return func(context.Context) error {
		return f()
	}
}

func (s *shutdown) do(ctx context.Context) {
	slices.Reverse(s.fn)

	for _, fn := range s.fn {
		if err := fn(ctx); err != nil {
			logger.Ctx(ctx).Error(err.Error())
		}
	}
}
