package logger

import "go.uber.org/zap"

type ZapLogger struct {
	l *zap.Logger
}

func (z *ZapLogger) Debug(msg string, args ...Field) {
	z.l.Debug(msg, z.toZapFields(args)...)

}

func (z *ZapLogger) Info(msg string, args ...Field) {

	// func (log *Logger) Info(msg string, fields ...Field)
	z.l.Info(msg, z.toZapFields(args)...)
}

func (z *ZapLogger) Warn(msg string, args ...Field) {
	z.l.Warn(msg, z.toZapFields(args)...)
}

func (z *ZapLogger) Error(msg string, args ...Field) {
	z.l.Error(msg, z.toZapFields(args)...)

}

func NewZapLogger(l *zap.Logger) *ZapLogger {
	return &ZapLogger{
		l: l,
	}
}

func (z *ZapLogger) toZapFields(args []Field) []zap.Field {
	res := make([]zap.Field, 0, len(args))
	for _, arg := range args {
		res = append(res, zap.Any(arg.Key, arg.Value))
		// zap.Any Any(key string, value interface{}) Field 返回zap.Field
	}
	return res
}
