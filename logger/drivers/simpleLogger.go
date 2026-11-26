package drivers

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/ArtisanCloud/PowerLibs/v3/logger/contract"
)

const (
	DebugLevel   int8 = 0
	InfoLevel    int8 = 1
	WarningLevel int8 = 2
	ErrorLevel   int8 = 3
	PanicLevel   int8 = 4
	FatalLevel   int8 = 5
)

type SimpleLogger struct {
	ctx context.Context
}

func (l *SimpleLogger) WithContext(ctx context.Context) contract.LoggerInterface {
	l.ctx = ctx
	return l
}

func simpleMessage(msg string, args []interface{}) string {
	if len(args) == 0 {
		return msg
	}
	// SimpleLogger also accepts legacy printf templates on methods without an F suffix.
	if strings.Contains(msg, "%") {
		return fmt.Sprintf(msg, args...)
	}

	var line strings.Builder
	line.WriteString(msg)
	for len(args) >= 2 {
		fmt.Fprintf(&line, " %v=%v", args[0], args[1])
		args = args[2:]
	}
	if len(args) == 1 {
		fmt.Fprintf(&line, " %v", args[0])
	}
	return line.String()
}

func (l *SimpleLogger) Debug(msg string, v ...interface{}) {
	log.Print("[DEBUG] ", simpleMessage(msg, v))
}

func (l *SimpleLogger) Info(msg string, v ...interface{}) {
	log.Print("[INFO] ", simpleMessage(msg, v))
}

func (l *SimpleLogger) Warn(msg string, v ...interface{}) {
	log.Print("[WARNING] ", simpleMessage(msg, v))
}

func (l *SimpleLogger) Error(msg string, v ...interface{}) {
	log.Print("[ERROR] ", simpleMessage(msg, v))
}

func (l *SimpleLogger) Panic(msg string, v ...interface{}) {
	log.Panic("[PANIC] ", simpleMessage(msg, v))
}

func (l *SimpleLogger) Fatal(msg string, v ...interface{}) {
	log.Fatal("[FATAL] ", simpleMessage(msg, v))
}

func (l *SimpleLogger) DebugF(format string, args ...interface{}) {
	log.Printf("[DEBUG] "+format, args...)
}

func (l *SimpleLogger) InfoF(format string, args ...interface{}) {
	log.Printf("[INFO] "+format, args...)
}

func (l *SimpleLogger) WarnF(format string, args ...interface{}) {
	log.Printf("[WARNING] "+format, args...)
}

func (l *SimpleLogger) ErrorF(format string, args ...interface{}) {
	log.Printf("[ERROR] "+format, args...)
}

func (l *SimpleLogger) PanicF(format string, args ...interface{}) {
	log.Panicf("[PANIC] "+format, args...)
}

func (l *SimpleLogger) FatalF(format string, args ...interface{}) {
	log.Fatalf("[FATAL] "+format, args...)
}
