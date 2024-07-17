package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"simpleMicroservice/internal/config"
	"time"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

var defaultLogger *slog.Logger
var logLevel slog.Level

type ctxTimeStampKey struct{}
type ctxEnvStringKey struct{}
type ctxAddressStringKey struct{}

func ContextWithTimeStamp(ctx context.Context, ts time.Time) context.Context {
	return context.WithValue(ctx, ctxTimeStampKey{}, ts)
}

func ContextWithEnvString(ctx context.Context, env string) context.Context {
	return context.WithValue(ctx, ctxEnvStringKey{}, env)
}

func ContextWithAddressString(ctx context.Context, address string) context.Context {
	return context.WithValue(ctx, ctxAddressStringKey{}, address)
}

func TimeStampFromContext(ctx context.Context) (time.Time, bool) {
	ts, ok := ctx.Value(ctxTimeStampKey{}).(time.Time)
	return ts, ok
}

func AddressFromContext(ctx context.Context) (string, bool) {
	addr, ok := ctx.Value(ctxAddressStringKey{}).(string)
	return addr, ok
}

func EnvFromContext(ctx context.Context) (string, bool) {
	env, ok := ctx.Value(ctxEnvStringKey{}).(string)
	return env, ok
}

type handler struct {
	slog.Handler
}

func (h *handler) Handle(ctx context.Context, r slog.Record) error {
	ts, ok := TimeStampFromContext(ctx)
	if ok {
		r.AddAttrs(slog.String("duration", time.Since(ts).String()))
	}

	env, ok := EnvFromContext(ctx)
	if ok {
		r.AddAttrs(slog.String("env", env))
	}

	addr, ok := AddressFromContext(ctx)
	if ok {
		r.AddAttrs(slog.String("address", addr))
	}

	return h.Handler.Handle(ctx, r)
}

func init() {
	h := handler{Handler: slog.NewTextHandler(os.Stdout, nil)}
	defaultLogger = slog.New(&h)
}

func setUpLogLevel(env string) slog.Level {
	switch env {
	case envLocal:
		fallthrough
	case envDev:
		return slog.LevelDebug
	case envProd:
		return slog.LevelInfo
	}

	return slog.LevelInfo
}

func SetUpLogger(cfg *config.Config) {
	env := cfg.Env
	level := setUpLogLevel(env)
	logLevel = level

	handlerOptions := &slog.HandlerOptions{Level: level}

	var h handler

	switch env {
	case envLocal:
		h = handler{Handler: slog.NewTextHandler(os.Stdout, handlerOptions)}
	case envDev:
		fallthrough
	case envProd:
		h = handler{Handler: slog.NewJSONHandler(os.Stdout, handlerOptions)}
	}

	defaultLogger = slog.New(&h)
}

func Info(ctx context.Context, msg string, args ...any) {
	defaultLogger.Log(ctx, logLevel, msg, args...)
}

func Infof(ctx context.Context, format string, args ...any) {
	defaultLogger.Log(ctx, logLevel, fmt.Sprintf(format, args...))
}

func Error(ctx context.Context, msg string, args ...any) {
	defaultLogger.Log(ctx, logLevel, msg, args...)
}

func Errorf(ctx context.Context, format string, args ...any) {
	defaultLogger.Log(ctx, logLevel, fmt.Sprintf(format, args...))
}

func Fatal(ctx context.Context, msg string, args ...any) {
	defaultLogger.Log(ctx, logLevel, msg, args...)
	os.Exit(1)
}

func Fatalf(ctx context.Context, format string, args ...any) {
	defaultLogger.Log(ctx, logLevel, fmt.Sprintf(format, args...))
	os.Exit(1)
}
