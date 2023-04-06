package logger

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var defaultLogger *zap.Logger

// Logger implement the Module interface
type Logger struct {
	AddCaller  bool
	Stacktrace zapcore.LevelEnabler
}

func (l *Logger) Register(*cobra.Command) {
	viper.SetDefault("settings.log.Level", "info")
	viper.SetDefault("settings.log.format", "console")
	viper.SetDefault("settings.log.output", "stdout")
}

func (l *Logger) MustCheck(*cobra.Command) {}

func (l *Logger) Initialize(*cobra.Command) error {
	logConfig := LogConfig{
		addCaller:  l.AddCaller,
		stacktrace: l.Stacktrace,
	}
	err := viper.Sub("settings.log").Unmarshal(&logConfig)
	if err != nil {
		return err
	}

	return logConfig.setDefault()
}

// LogConfig zap log config
type LogConfig struct {
	Level  string
	Format string
	Output string

	addCaller  bool
	stacktrace zapcore.LevelEnabler

	preOutputFiles []*os.File
	curOutputFiles []*os.File
}

func (l *LogConfig) setDefault() error {
	newLogger, err := l.newLogger()
	if err != nil {
		return err
	}

	if defaultLogger != nil {
		_ = defaultLogger.Sync()
		for _, file := range l.preOutputFiles {
			_ = file.Close()
		}
	}

	defaultLogger = newLogger
	_ = zap.ReplaceGlobals(defaultLogger)

	return nil
}

func (l *LogConfig) newLogger() (*zap.Logger, error) {
	level, err := l._newLevel()
	if err != nil {
		return nil, err
	}

	encoder, err := l._newEncoder()
	if err != nil {
		return nil, err
	}

	syncers, err := l._newOutput()
	if err != nil {
		return nil, err
	}

	logger := zap.New(zapcore.NewCore(encoder, syncers, level))

	if l.addCaller {
		logger = logger.WithOptions(zap.AddCaller(), zap.AddCallerSkip(1))
	}

	if l.stacktrace != nil {
		logger = logger.WithOptions(zap.AddStacktrace(l.stacktrace))
	}

	return logger, nil
}

func (l *LogConfig) _newLevel() (zap.AtomicLevel, error) {
	level, err := zapcore.ParseLevel(l.Level)
	if err != nil {
		unrecognized := "unrecognized Level: " + l.Level
		supported := "supported values: debug,info,warn,error,dpanic,panic,fatal"
		return zap.NewAtomicLevelAt(zapcore.InvalidLevel), fmt.Errorf(unrecognized + ", " + supported)
	}
	return zap.NewAtomicLevelAt(level), nil
}

func (l *LogConfig) _newEncoder() (zapcore.Encoder, error) {
	// encoderConfig
	zap.NewProductionEncoderConfig()
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:          "time",
		LevelKey:         "level",
		NameKey:          "logger",
		CallerKey:        "caller",
		FunctionKey:      zapcore.OmitKey,
		MessageKey:       "message",
		StacktraceKey:    "stacktrace",
		LineEnding:       zapcore.DefaultLineEnding,
		EncodeLevel:      zapcore.LowercaseLevelEncoder,
		EncodeTime:       zapcore.TimeEncoderOfLayout(time.DateTime),
		EncodeDuration:   zapcore.SecondsDurationEncoder,
		EncodeCaller:     zapcore.ShortCallerEncoder,
		EncodeName:       func(s string, encoder zapcore.PrimitiveArrayEncoder) { encoder.AppendString(s) },
		ConsoleSeparator: "",
	}

	// encoder
	switch l.Format {
	case "json":
		return zapcore.NewJSONEncoder(encoderConfig), nil
	case "console":
		encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
		return zapcore.NewConsoleEncoder(encoderConfig), nil
	default:
		unrecognized := "unrecognized format: " + l.Format
		supported := "supported values: json,console"
		return zapcore.NewConsoleEncoder(encoderConfig), fmt.Errorf(unrecognized + ", " + supported)
	}
}

func (l *LogConfig) _newOutput() (zapcore.WriteSyncer, error) {
	var (
		writeSyncers   []zapcore.WriteSyncer
		newOutputFiles []*os.File
	)

	for _, out := range strings.Split(l.Output, ",") {
		switch out {
		case "stdout":
			writeSyncers = append(writeSyncers, zapcore.AddSync(os.Stdout))
		case "stderr":
			writeSyncers = append(writeSyncers, zapcore.AddSync(os.Stderr))
		default:
			file, err := os.OpenFile(out, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
			if err != nil {
				return nil, err
			}
			writeSyncers = append(writeSyncers, zapcore.Lock(zapcore.AddSync(file)))
			newOutputFiles = append(newOutputFiles, file)
		}
	}

	l.preOutputFiles = l.curOutputFiles
	l.curOutputFiles = newOutputFiles

	return zapcore.NewMultiWriteSyncer(writeSyncers...), nil
}

func check() {
	if defaultLogger == nil {
		panic("logger not initialized, module lib needs to be added")
	}
}

func Debug(msg string, fields ...zap.Field) {
	check()
	defaultLogger.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	check()
	defaultLogger.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	check()
	defaultLogger.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	check()
	defaultLogger.Error(msg, fields...)
}

func DPanic(msg string, fields ...zap.Field) {
	check()
	defaultLogger.DPanic(msg, fields...)
}

func Panic(msg string, fields ...zap.Field) {
	check()
	defaultLogger.Panic(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	check()
	defaultLogger.Fatal(msg, fields...)
}
