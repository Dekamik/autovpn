package logger

import (
	"fmt"
	"log"
	"strings"
)

const logFmtMessageFormat = "level=%s msg=\"%s\"\n"

var _ Logger = &LogFmtLogger{}

type LogFmtLogger struct {
	level  int64
	logger *log.Logger
}

func (l *LogFmtLogger) Log(level int64, message string) {
	if level >= l.level {
		l.logger.Printf(logFmtMessageFormat, levelNames[level], message)
	}
}

func (l *LogFmtLogger) LogTrace(message string) {
	l.Log(TRACE, message)
}

func (l *LogFmtLogger) LogDebug(message string) {
	l.Log(DEBUG, message)
}

func (l *LogFmtLogger) LogInformation(message string) {
	l.Log(INFORMATION, message)
}

func (l *LogFmtLogger) LogWarning(message string) {
	l.Log(WARNING, message)
}

func (l *LogFmtLogger) LogError(message string) {
	l.Log(ERROR, message)
}

func (l *LogFmtLogger) LogFatal(message string) {
	l.logger.Fatalf(logFmtMessageFormat, levelNames[FATAL], message)
}

func (l *LogFmtLogger) LogPanic(message string) {
	l.logger.Panicf(logFmtMessageFormat, levelNames[FATAL], message)
}

func NewLogFmtLogger(baseLogger *log.Logger, level int64, properties map[string]string) Logger {
	var prefix = ""
	for k, v := range properties {
		prefix += fmt.Sprintf("%s=%s ", replaceSpaces(k), addQuotes(v))
	}
	baseLogger.SetPrefix(prefix)

	return &LogFmtLogger{
		level:  level,
		logger: baseLogger,
	}
}

func addQuotes(s string) string {
	if strings.Contains(s, " ") {
		return fmt.Sprintf("\"%s\"", s)
	}
	return s
}

func replaceSpaces(s string) string {
	if strings.Contains(s, " ") {
		return strings.Replace(s, " ", "_", -1)
	}
	return s
}
