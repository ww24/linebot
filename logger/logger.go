package logger

import (
	"github.com/blendle/zapdriver"
	"go.uber.org/zap"
)

func New(name, version string) (*zap.Logger, error) {
	core := zapdriver.WrapCore(
		zapdriver.ReportAllErrors(true),
		zapdriver.ServiceName(name),
	)

	logger, err := zapdriver.NewProductionWithCore(core)
	if err != nil {
		return nil, err
	}

	logger = logger.With(zap.String("version", version))

	return logger, nil
}
