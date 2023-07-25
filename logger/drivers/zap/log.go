package zap

import (
	"github.com/ArtisanCloud/PowerLibs/v3/logger/contract"
	"github.com/ArtisanCloud/PowerLibs/v3/object"
	os2 "github.com/ArtisanCloud/PowerLibs/v3/os"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

type Logger struct {
	Logger *zap.Logger
	sugar  *zap.SugaredLogger
}

func NewLogger(config *object.HashMap) (logger contract.LoggerInterface, err error) {

	zapLogger, err := newZapLogger(config)
	if err != nil {
		return nil, err
	}

	defer zapLogger.Sync() // flushes buffer, if any

	logger = &Logger{
		Logger: zapLogger,
		sugar:  zapLogger.Sugar(),
	}

	return logger, err
}

func newZapLogger(config *object.HashMap) (logger *zap.Logger, err error) {
	env := (*config)["env"].(string)
	var loggerConfig zap.Config
	if env == "production" {
		loggerConfig = zap.NewProductionConfig()
	} else {
		loggerConfig = zap.NewDevelopmentConfig()
	}

	outputFile := (*config)["outputPath"].(string)
	errorFile := (*config)["errorPath"].(string)

	err = os2.CreateDirectoriesForFiles(outputFile)
	if err != nil {
		return nil, err
	}
	err = os2.CreateDirectoriesForFiles(errorFile)
	if err != nil {
		return nil, err
	}

	loggerConfig.OutputPaths = []string{outputFile}
	loggerConfig.ErrorOutputPaths = []string{errorFile}
	loggerConfig.EncoderConfig.TimeKey = "timestamp"
	loggerConfig.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)

	outputSyncer, err := newFileWriteSyncer(outputFile)
	if err != nil {
		return nil, err
	}

	errorSyncer, err := newFileWriteSyncer(errorFile)
	if err != nil {
		return nil, err
	}

	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.ErrorLevel
	})

	errorLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})

	core := zapcore.NewTee(
		zapcore.NewCore(
			zapcore.NewJSONEncoder(loggerConfig.EncoderConfig),
			zapcore.Lock(outputSyncer),
			infoLevel,
		),
		zapcore.NewCore(
			zapcore.NewJSONEncoder(loggerConfig.EncoderConfig),
			zapcore.Lock(errorSyncer),
			errorLevel,
		),
	)

	logger = zap.New(core)

	return logger, err
}

func newFileWriteSyncer(filename string) (zapcore.WriteSyncer, error) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	return zapcore.AddSync(file), nil
}
func (log *Logger) Debug(msg string, v ...interface{}) {
	log.sugar.Debugw(msg, v...)
}
func (log *Logger) Info(msg string, v ...interface{}) {
	log.sugar.Infow(msg, v...)
}
func (log *Logger) Warn(msg string, v ...interface{}) {
	log.sugar.Warnw(msg, v...)
}
func (log *Logger) Error(msg string, v ...interface{}) {
	log.sugar.Errorw(msg, v...)
}
func (log *Logger) Panic(msg string, v ...interface{}) {
	log.sugar.Panicw(msg, v...)
}
func (log *Logger) Fatal(msg string, v ...interface{}) {
	log.sugar.Fatalw(msg, v...)
}

func (log *Logger) DebugF(format string, args ...interface{}) {
	log.sugar.Debugf(format, args...)
}
func (log *Logger) InfoF(format string, args ...interface{}) {
	log.sugar.Infof(format, args...)
}
func (log *Logger) WarnF(format string, args ...interface{}) {
	log.sugar.Warnf(format, args...)
}
func (log *Logger) ErrorF(format string, args ...interface{}) {
	log.sugar.Errorf(format, args...)
}
func (log *Logger) PanicF(format string, args ...interface{}) {
	log.sugar.Panicf(format, args...)
}
func (log *Logger) FatalF(format string, args ...interface{}) {
	log.sugar.Fatalf(format, args...)
}
