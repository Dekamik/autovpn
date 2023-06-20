package logger

const (
	TRACE       int64 = -20
	DEBUG             = -10
	INFORMATION       = 0
	WARNING           = 10
	ERROR             = 20
	FATAL             = 30
)

var levelNames = map[int64]string{
	TRACE:       "trace",
	DEBUG:       "debug",
	INFORMATION: "info",
	WARNING:     "warn",
	ERROR:       "errors",
	FATAL:       "fatal",
}

type Logger interface {
	// Log is the base logging method
	Log(level int64, message string)

	// LogTrace should be used for tracing program execution
	LogTrace(message string)

	// LogDebug should be used for less important events that may provide important information when solving errors
	LogDebug(message string)

	// LogInformation should be called when important events occur in the application
	LogInformation(message string)

	// LogWarning should be called when errors that don't halt executions occur
	LogWarning(message string)

	// LogError should be called when execution-halting errors occur
	LogError(message string)

	// LogFatal logs to stdout and calls os.Exit(1)
	// Should be called when very serious errors occur
	LogFatal(message string)

	// LogPanic logs to stdout and calls panic()
	// Can be called when e.g. unhandled errors occur
	LogPanic(message string)
}
