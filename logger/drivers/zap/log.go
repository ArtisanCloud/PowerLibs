package zap

import (
	"github.com/ArtisanCloud/PowerLibs/v2/logger/contract"
	"github.com/ArtisanCloud/PowerLibs/v2/object"
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
	//loggerConfig.OutputPaths = []string{(*config)["outputPath"].(string)}
	loggerConfig.ErrorOutputPaths = []string{(*config)["errorPath"].(string)}
	loggerConfig.EncoderConfig.TimeKey = "timestamp"
	loggerConfig.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)

	var file *os.File
	if _, err = os.Stat(outputFile); os.IsNotExist(err) {
		file, err = os.Create(outputFile)
	} else {
		if file, err = os.OpenFile(outputFile, os.O_APPEND|os.O_WRONLY, os.ModeAppend); err != nil {
			return nil, err
		}

	}
	if err != nil {
		return nil, err
	}
	writeSyncer := zapcore.AddSync(file)

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(loggerConfig.EncoderConfig),
		zapcore.Lock(zapcore.NewMultiWriteSyncer(os.Stdout, writeSyncer)),
		zapcore.DebugLevel,
	)

	logger = zap.New(core)

	return logger, err
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
