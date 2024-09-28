package shadowtls

import (
	"context"
	"fmt"

	L "github.com/sagernet/sing/common/logger"
)

type singLogger struct{}

func (l singLogger) TraceContext(ctx context.Context, args ...any) {
	fmt.Println(args...)
}

func (l singLogger) DebugContext(ctx context.Context, args ...any) {
	fmt.Println(args...)
}

func (l singLogger) InfoContext(ctx context.Context, args ...any) {
	fmt.Println(args...)
}

func (l singLogger) WarnContext(ctx context.Context, args ...any) {
	fmt.Println(args...)
}

func (l singLogger) ErrorContext(ctx context.Context, args ...any) {
	fmt.Println(args...)
}

func (l singLogger) FatalContext(ctx context.Context, args ...any) {
	fmt.Println(args...)
}

func (l singLogger) PanicContext(ctx context.Context, args ...any) {
	fmt.Println(args...)
}

func (l singLogger) Trace(args ...any) {
	fmt.Println(args...)
}

func (l singLogger) Debug(args ...any) {
	fmt.Println(args...)
}

func (l singLogger) Info(args ...any) {
	fmt.Println(args...)
}

func (l singLogger) Warn(args ...any) {
	fmt.Println(args...)
}

func (l singLogger) Error(args ...any) {
	fmt.Println(args...)
}

func (l singLogger) Fatal(args ...any) {
	fmt.Println(args...)
}

func (l singLogger) Panic(args ...any) {
	fmt.Println(args...)
}

var SingLogger L.ContextLogger = singLogger{}
