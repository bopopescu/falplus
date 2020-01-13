package logrus

import (
	"io"
	"path"
	"runtime"
	"strconv"
	"strings"
)

var (
	// std is the name of the standard logger in stdlib `log`
	std = New()
)

func StandardLogger() *Logger {
	return std
}

// SetOutput sets the standard logger output.
func SetOutput(out io.Writer) {
	std.mu.Lock()
	defer std.mu.Unlock()
	std.Out = out
}

// SetFormatter sets the standard logger formatter.
func SetFormatter(formatter Formatter) {
	std.mu.Lock()
	defer std.mu.Unlock()
	std.Formatter = formatter
}

// SetLevel sets the standard logger level.
func SetLevel(level Level) {
	std.mu.Lock()
	defer std.mu.Unlock()
	std.Level = level
}

// GetLevel returns the standard logger level.
func GetLevel() Level {
	std.mu.Lock()
	defer std.mu.Unlock()
	return std.Level
}

// AddHook adds a hook to the standard logger hooks.
func AddHook(hook Hook) {
	std.mu.Lock()
	defer std.mu.Unlock()
	std.Hooks.Add(hook)
}

// WithError creates an entry from the standard logger and adds an error to it, using the value defined in ErrorKey as key.
func WithError(err error) *Entry {
	return std.WithField(ErrorKey, err)
}

// WithField creates an entry from the standard logger and adds a field to
// it. If you want multiple fields, use `WithFields`.
//
// Note that it doesn't log until you call Debug, Print, Info, Warn, Fatal
// or Panic on the Entry it returns.
func WithField(key string, value interface{}) *Entry {
	return std.WithField(key, value)
}

// WithFields creates an entry from the standard logger and adds multiple
// fields to it. This is simply a helper for `WithField`, invoking it
// once for each field.
//
// Note that it doesn't log until you call Debug, Print, Info, Warn, Fatal
// or Panic on the Entry it returns.
func WithFields(fields Fields) *Entry {
	return std.WithFields(fields)
}

// Debug logs a message at level Debug on the standard logger.
func Debug(args ...interface{}) {
	if std.level() < DebugLevel {
		return
	}
	newargs := Caller1(args...)
	std.Debug(newargs...)
}

// Print logs a message at level Info on the standard logger.
func Print(args ...interface{}) {
	if std.level() < InfoLevel {
		return
	}
	newargs := Caller1(args...)
	std.Print(newargs...)
}

// Info logs a message at level Info on the standard logger.
func Info(args ...interface{}) {
	if std.level() < InfoLevel {
		return
	}
	newargs := Caller1(args...)
	std.Info(newargs...)
}

// Warn logs a message at level Warn on the standard logger.
func Warn(args ...interface{}) {
	if std.level() < WarnLevel {
		return
	}
	newargs := Caller1(args...)
	std.Warn(newargs...)
}

// Warning logs a message at level Warn on the standard logger.
func Warning(args ...interface{}) {
	if std.level() < WarnLevel {
		return
	}
	newargs := Caller1(args...)
	std.Warning(newargs...)
}

// Error logs a message at level Error on the standard logger.
func Error(args ...interface{}) {
	if std.level() < ErrorLevel {
		return
	}
	newargs := Caller1(args...)
	std.Error(newargs...)
}

// Panic logs a message at level Panic on the standard logger.
func Panic(args ...interface{}) {
	if std.level() < PanicLevel {
		return
	}
	newargs := Caller1(args...)
	std.Panic(newargs...)
}

// Fatal logs a message at level Fatal on the standard logger.
func Fatal(args ...interface{}) {
	if std.level() < FatalLevel {
		return
	}
	newargs := Caller1(args...)
	std.Fatal(newargs...)
}

// Debugf logs a message at level Debug on the standard logger.
func Debugf(format string, args ...interface{}) {
	if std.level() < DebugLevel {
		return
	}
	format = Caller2(format)
	std.Debugf(format, args...)
}

// Printf logs a message at level Info on the standard logger.
func Printf(format string, args ...interface{}) {
	if std.level() < InfoLevel {
		return
	}
	format = Caller2(format)
	std.Printf(format, args...)
}

// Infof logs a message at level Info on the standard logger.
func Infof(format string, args ...interface{}) {
	if std.level() < InfoLevel {
		return
	}
	format = Caller2(format)
	std.Infof(format, args...)
}

// Warnf logs a message at level Warn on the standard logger.
func Warnf(format string, args ...interface{}) {
	if std.level() < WarnLevel {
		return
	}
	format = Caller2(format)
	std.Warnf(format, args...)
}

// Warningf logs a message at level Warn on the standard logger.
func Warningf(format string, args ...interface{}) {
	if std.level() < WarnLevel {
		return
	}
	format = Caller2(format)
	std.Warningf(format, args...)
}

// Errorf logs a message at level Error on the standard logger.
func Errorf(format string, args ...interface{}) {
	if std.level() < ErrorLevel {
		return
	}
	format = Caller2(format)
	std.Errorf(format, args...)
}

// Panicf logs a message at level Panic on the standard logger.
func Panicf(format string, args ...interface{}) {
	if std.level() < PanicLevel {
		return
	}
	format = Caller2(format)
	std.Panicf(format, args...)
}

// Fatalf logs a message at level Fatal on the standard logger.
func Fatalf(format string, args ...interface{}) {
	if std.level() < FatalLevel {
		return
	}
	format = Caller2(format)
	std.Fatalf(format, args...)
}

// Debugln logs a message at level Debug on the standard logger.
func Debugln(args ...interface{}) {
	if std.level() < DebugLevel {
		return
	}
	newargs := Caller1(args...)
	std.Debugln(newargs...)
}

// Println logs a message at level Info on the standard logger.
func Println(args ...interface{}) {
	if std.level() < InfoLevel {
		return
	}
	newargs := Caller1(args...)
	std.Println(newargs...)
}

// Infoln logs a message at level Info on the standard logger.
func Infoln(args ...interface{}) {
	if std.level() < InfoLevel {
		return
	}
	newargs := Caller1(args...)
	std.Infoln(newargs...)
}

// Warnln logs a message at level Warn on the standard logger.
func Warnln(args ...interface{}) {
	if std.level() < WarnLevel {
		return
	}
	newargs := Caller1(args...)
	std.Warnln(newargs...)
}

// Warningln logs a message at level Warn on the standard logger.
func Warningln(args ...interface{}) {
	if std.level() < WarnLevel {
		return
	}
	newargs := Caller1(args...)
	std.Warningln(newargs...)
}

// Errorln logs a message at level Error on the standard logger.
func Errorln(args ...interface{}) {
	if std.level() < ErrorLevel {
		return
	}
	newargs := Caller1(args...)
	std.Errorln(newargs...)
}

// Panicln logs a message at level Panic on the standard logger.
func Panicln(args ...interface{}) {
	if std.level() < PanicLevel {
		return
	}
	newargs := Caller1(args...)
	std.Panicln(newargs...)
}

// Fatalln logs a message at level Fatal on the standard logger.
func Fatalln(args ...interface{}) {
	if std.level() < FatalLevel {
		return
	}
	newargs := Caller1(args...)
	std.Fatalln(newargs...)
}

func Caller1(args ...interface{}) []interface{} {

	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "???"
		line = 0
	}
	funcName := "???"
	f := runtime.FuncForPC(pc)
	if f != nil {
		funcName = f.Name()
	}
	_, filename := path.Split(file)
	flist := strings.Split(funcName, ".")
	funcName = flist[len(flist)-1]
	format := "[" + filename + ":" + funcName + ":" + strconv.FormatInt(int64(line), 10) + "] "
	var newargs []interface{}
	newargs = append(newargs, format)
	newargs = append(newargs, args...)
	return newargs
}

func Caller2(format string) string {
	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "???"
		line = 0
	}
	funcName := "???"
	f := runtime.FuncForPC(pc)
	if f != nil {
		funcName = f.Name()
	}
	_, filename := path.Split(file)
	flist := strings.Split(funcName, ".")
	funcName = flist[len(flist)-1]
	format = "[" + filename + ":" + funcName + ":" + strconv.FormatInt(int64(line), 10) + "] " + format
	return format
}
