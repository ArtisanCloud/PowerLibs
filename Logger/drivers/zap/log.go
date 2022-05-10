package zap

import (
	"github.com/ArtisanCloud/PowerLibs/Logger/contract"
	"github.com/ArtisanCloud/PowerLibs/object"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

type Logger struct {
	Logger *zap.Logger
}

func NewLogger(config *object.HashMap) (logger contract.LoggerInterface, err error) {

	zapLogger, err := newZapLogger(config)
	if err != nil {
		return nil, err
	}
	logger = &Logger{
		Logger: zapLogger,
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
	loggerConfig.OutputPaths = []string{(*config)["outputPath"].(string)}
	loggerConfig.ErrorOutputPaths = []string{(*config)["errorPath"].(string)}
	loggerConfig.EncoderConfig.TimeKey = "timestamp"
	loggerConfig.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)

	logger, err = loggerConfig.Build()

	return logger, err
}

func (log *Logger) Debug(msg string, v ...interface{}) {
	log.Logger.Sugar().Debugw(msg, v)
}
func (log *Logger) Info(msg string, v ...interface{}) {
	log.Logger.Sugar().Infow(msg, v)
}
func (log *Logger) Warn(msg string, v ...interface{}) {
	log.Logger.Sugar().Warnw(msg, v)
}
func (log *Logger) Error(msg string, v ...interface{}) {
	log.Logger.Sugar().Errorw(msg, v)
}
func (log *Logger) Panic(msg string, v ...interface{}) {
	log.Logger.Sugar().Panicw(msg, v)
}
func (log *Logger) Fatal(msg string, v ...interface{}) {
	log.Logger.Sugar().Fatalw(msg, v)
}

func (log *Logger) DebugF(format string, args ...interface{}) {
	log.Logger.Sugar().Debugf(format, args)
}
func (log *Logger) InfoF(format string, args ...interface{}) {
	log.Logger.Sugar().Infof(format, args)
}
func (log *Logger) WarnF(format string, args ...interface{}) {
	log.Logger.Sugar().Warnf(format, args)
}
func (log *Logger) ErrorF(format string, args ...interface{}) {
	log.Logger.Sugar().Errorf(format, args)
}
func (log *Logger) PanicF(format string, args ...interface{}) {
	log.Logger.Sugar().Panicf(format, args)
}
func (log *Logger) FatalF(format string, args ...interface{}) {
	log.Logger.Sugar().Fatalf(format, args)
}
