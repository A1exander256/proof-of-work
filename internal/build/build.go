package build

import (
	"github.com/proof-of-work/internal/config"
	"github.com/proof-of-work/pkg/logger"
)

type Builder struct {
	l   *logger.Logger
	cfg config.Config

	shutdown shutdown
}

func New(l *logger.Logger, cfg config.Config) *Builder {
	//nolint:exhaustruct
	return &Builder{l: l, cfg: cfg}
}
