package logging

import (
	"io"
	"os"
	"strings"
	"time"

	"github.com/mattn/go-colorable"
	"github.com/rs/zerolog"
)

// Logger wraps zerolog.Logger and provides convenience methods.
type Logger struct {
	logger zerolog.Logger
	output io.Writer
}

var (
	globalLogger    *Logger
	globalNoColor   bool
	currentLevel    zerolog.Level
	privacyMode     bool
	sensitiveFields = map[string]bool{
		"ip":         true,
		"isp":        true,
		"ip_address": true,
		"user_ip":    true,
		"user_isp":   true,
	}
)

// New creates and returns a new Logger instance with default settings.
func New() *Logger {
	var output io.Writer = colorable.NewColorableStderr()
	if globalNoColor {
		output = os.Stderr
	}
	opts := &zerolog.ConsoleWriter{Out: output, TimeFormat: time.RFC3339, NoColor: globalNoColor}
	logger := zerolog.New(opts).Level(zerolog.InfoLevel).With().Timestamp().Caller().Logger()

	globalLogger = &Logger{logger: logger, output: output}
	currentLevel = zerolog.InfoLevel
	return globalLogger
}

// Info logs a message at info level.
func (l *Logger) Info(msg string) {
	l.logger.Info().Msg(msg)
}

// InfoWithFields logs a message at info level with structured fields.
// Sensitive fields are redacted if privacy mode is enabled.
func (l *Logger) InfoWithFields(msg string, fields map[string]interface{}) {
	e := l.logger.Info()
	for k, v := range fields {
		if privacyMode && sensitiveFields[strings.ToLower(k)] {
			e = e.Interface(k, "[REDACTED]")
		} else {
			e = e.Interface(k, v)
		}
	}
	e.Msg(msg)
}

// Error logs a message at error level.
func (l *Logger) Error(msg string) {
	l.logger.Error().Msg(msg)
}

// Errorf logs a formatted message at error level.
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.logger.Error().Msgf(format, args...)
}

// ErrorWithFields logs a message at error level with structured fields.
// Sensitive fields are redacted if privacy mode is enabled.
func (l *Logger) ErrorWithFields(msg string, fields map[string]interface{}) {
	e := l.logger.Error()
	for k, v := range fields {
		if privacyMode && sensitiveFields[strings.ToLower(k)] {
			e = e.Interface(k, "[REDACTED]")
		} else {
			e = e.Interface(k, v)
		}
	}
	e.Msg(msg)
}

// Debug logs a message at debug level.
func (l *Logger) Debug(msg string) {
	l.logger.Debug().Msg(msg)
}

// Debugf logs a formatted message at debug level.
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.logger.Debug().Msgf(format, args...)
}

// Warn logs a message at warn level.
func (l *Logger) Warn(msg string) {
	l.logger.Warn().Msg(msg)
}

// Fatal logs a message at fatal level and exits the application.
func (l *Logger) Fatal(msg string) {
	l.logger.Fatal().Msg(msg)
}

// Fatalf logs a formatted message at fatal level and exits the application.
func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.logger.Fatal().Msgf(format, args...)
}

// Sync flushes any buffered log data to the underlying writer.
func (l *Logger) Sync() {
	if l.output == nil {
		return
	}
	if f, ok := l.output.(*os.File); ok {
		f.Sync()
	}
}

// Get returns the global logger instance, creating one if necessary.
func Get() *Logger {
	if globalLogger == nil {
		return New()
	}
	return globalLogger
}

// SetLevel sets the log level for the global logger.
// Valid levels: trace, debug, info, warn, error, fatal, panic.
func SetLevel(level string) error {
	lvl, err := zerolog.ParseLevel(level)
	if err != nil {
		return err
	}

	currentLevel = lvl

	if globalLogger == nil {
		New()
	}

	globalLogger.logger = globalLogger.logger.Level(lvl)
	return nil
}

// SetNoColor enables or disables colored output for the global logger.
func SetNoColor(noColor bool) {
	globalNoColor = noColor
	if globalLogger != nil {
		var output io.Writer = colorable.NewColorableStderr()
		if noColor {
			output = os.Stderr
		}
		opts := &zerolog.ConsoleWriter{Out: output, TimeFormat: time.RFC3339, NoColor: noColor}
		logger := zerolog.New(opts).Level(currentLevel).With().Timestamp().Caller().Logger()
		globalLogger.logger = logger
		globalLogger.output = output
	}
}

// SetPrivacyMode enables or disables privacy mode. When enabled, sensitive fields
// (ip, isp, etc.) are redacted as "[REDACTED]" in log output.
func SetPrivacyMode(enabled bool) {
	privacyMode = enabled
}

// SetOutput sets the output writer for the global logger.
func SetOutput(w io.Writer) {
	if globalLogger == nil {
		New()
	}
	var noColor bool
	if w == io.Discard {
		noColor = true
	}
	opts := &zerolog.ConsoleWriter{Out: w, TimeFormat: time.RFC3339, NoColor: noColor}
	logger := zerolog.New(opts).Level(currentLevel).With().Timestamp().Caller().Logger()
	globalLogger.logger = logger
	globalLogger.output = w
}

// Silence redirects log output to io.Discard, effectively disabling log output.
func Silence() {
	SetOutput(io.Discard)
}

// Restore resets log output to the default (stderr with color support).
func Restore() {
	var output io.Writer = colorable.NewColorableStderr()
	if globalNoColor {
		output = os.Stderr
	}
	opts := &zerolog.ConsoleWriter{Out: output, TimeFormat: time.RFC3339, NoColor: globalNoColor}
	logger := zerolog.New(opts).Level(currentLevel).With().Timestamp().Caller().Logger()
	globalLogger.logger = logger
	globalLogger.output = output
}
