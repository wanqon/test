package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"pay.sc.weibo.com/accounts/tool"
	"strings"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
)

type LogField map[string]interface{}
var zapper *zap.Logger

func init()  {
	//zapcore.NewJSONEncoder
	encoder := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		TimeKey:        "ts",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		},
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendInt64(int64(d) / 1000000)
		},
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	})

	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl <= zapcore.InfoLevel
	})

	errorLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})

	//kafkaEncoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())

	basePath := tool.GetConfString("dir","log")
	infoWriter := getWriter(basePath + "/info.log")
	errorWriter := getWriter(basePath + "/error.log")

	core := zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.AddSync(infoWriter), infoLevel).With([]zap.Field{zap.Int("foo", 42), zap.String("bar", "baz")}),
		zapcore.NewCore(encoder, zapcore.AddSync(errorWriter), errorLevel),
	)

	//log := zap.New(core, zap.AddCaller(),zap.AddCallerSkip(1),zap.AddStacktrace(errorLevel))
	zapper = zap.New(core, zap.AddCaller(),zap.AddCallerSkip(1))
	//errorLogger = logger.Sugar()

}


func getWriter(filename string) io.Writer {
	hook, err := rotatelogs.New(
		strings.Replace(filename, ".log", "", -1) + "-%Y%m%d.log",
	)
	if err != nil {
		panic(err)
	}
	return hook
}

func Debug(msg string, fields ...zap.Field)  {
	zapper.Debug(msg, fields...)
}

func getFields(params LogField) []zap.Field {
	fields := make([]zap.Field,0)
	for k,v := range params {
		switch x := v.(type) {
		case int:
			fields = append(fields, zap.Int(k,x))
		case string:
			fields = append(fields, zap.String(k,x))
		}
	}
	return fields
}

func Info(msg string,logfields LogField)  {
	fields := getFields(logfields)
	zapper.Info(msg, fields...)
}

func Warn(msg string, logfields LogField)  {
	fields := getFields(logfields)
	zapper.Warn(msg, fields...)
}

func Error(msg string, logfields LogField)  {
	fields := getFields(logfields)
	zapper.Error(msg, fields...)
}

func DPanic(msg string, logfields LogField)  {
	fields := getFields(logfields)
	zapper.DPanic(msg, fields...)
}

func Panic(msg string, logfields LogField)  {
	fields := getFields(logfields)
	zapper.Panic(msg, fields...)
}

func Fatal(msg string, logfields LogField)  {
	fields := getFields(logfields)
	zapper.Fatal(msg, fields...)
}