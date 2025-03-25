package util

import (
	"context"

	"github.com/sirupsen/logrus"
)

func GetLogger(ctx context.Context) *logrus.Entry {
	if ctx == nil {
		return logrus.NewEntry(logrus.StandardLogger())
	}
	if ctx.Value("logger") == nil {
		return logrus.NewEntry(logrus.StandardLogger())
	}
	return ctx.Value("logger").(*logrus.Entry)
}

func WithLogger(ctx context.Context, logger *logrus.Entry) context.Context {
	return context.WithValue(ctx, "logger", logger)
}
