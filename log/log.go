package log

import (
    "fmt"
    "io"
    "log"
    "os"
)

type LogLevel int

const (
    LOG_LEVEL_DEBUG LogLevel = iota
    LOG_LEVEL_INFO
    LOG_LEVEL_WARN
    LOG_LEVEL_ERROR
)

var defaultLog *DefaultLogger

type DefaultLogger struct {
    *log.Logger
    level LogLevel
}

func (l *DefaultLogger) SetLevel(level LogLevel) {
    l.level = level
}

func (l *DefaultLogger) LogLevel() LogLevel {
    return l.level
}

func (l *DefaultLogger) Debug(args ...interface{}) {
    if l.level <= LOG_LEVEL_DEBUG {
        l.SetPrefix("[DEBUG]")
        l.Output(2, fmt.Sprint(args...))
    }
}

func (l *DefaultLogger) Debugln(args ...interface{}) {
    if l.level <= LOG_LEVEL_DEBUG {
        l.SetPrefix("[DEBUG]")
        l.Output(2, fmt.Sprintln(args...))
    }
}

func (l *DefaultLogger) Debugf(format string, args ...interface{}) {
    if l.level <= LOG_LEVEL_DEBUG {
        l.SetPrefix("[DEBUG]")
        l.Output(2, fmt.Sprintf(format, args...))
    }
}

func (l *DefaultLogger) Info(args ...interface{}) {
    if l.level <= LOG_LEVEL_INFO {
        l.SetPrefix("[INFO] ")
        l.Output(2, fmt.Sprint(args...))
    }
}

func (l *DefaultLogger) Infoln(args ...interface{}) {
    if l.level <= LOG_LEVEL_INFO {
        l.SetPrefix("[INFO] ")
        l.Output(2, fmt.Sprintln(args...))
    }
}

func (l *DefaultLogger) Infof(format string, args ...interface{}) {
    if l.level <= LOG_LEVEL_INFO {
        l.SetPrefix("[INFO] ")
        l.Output(2, fmt.Sprintf(format, args...))
    }
}

func (l *DefaultLogger) Warn(args ...interface{}) {
    if l.level <= LOG_LEVEL_WARN {
        l.SetPrefix("[WARN] ")
        l.Output(2, fmt.Sprint(args...))
    }
}

func (l *DefaultLogger) Warnln(args ...interface{}) {
    if l.level <= LOG_LEVEL_WARN {
        l.SetPrefix("[WARN] ")
        l.Output(2, fmt.Sprintln(args...))
    }
}

func (l *DefaultLogger) Warnf(format string, args ...interface{}) {
    if l.level <= LOG_LEVEL_WARN {
        l.SetPrefix("[WARN] ")
        l.Output(2, fmt.Sprintf(format, args...))
    }
}

func (l *DefaultLogger) Error(args ...interface{}) {
    if l.level <= LOG_LEVEL_ERROR {
        l.SetPrefix("[ERROR]")
        l.Output(2, fmt.Sprint(args...))
    }
}

func (l *DefaultLogger) Errorln(args ...interface{}) {
    if l.level <= LOG_LEVEL_ERROR {
        l.SetPrefix("[ERROR]")
        l.Output(2, fmt.Sprintln(args...))
    }
}

func (l *DefaultLogger) Errorf(format string, args ...interface{}) {
    if l.level <= LOG_LEVEL_ERROR {
        l.SetPrefix("[ERROR]")
        l.Output(2, fmt.Sprintf(format, args...))
    }
}

func New(out io.Writer, prefix string, flag int, lv LogLevel) *DefaultLogger {
    return &DefaultLogger{
        Logger: log.New(out, prefix, flag),
        level:  lv,
    }
}

func init() {
    defaultLog = New(os.Stderr, "", log.LstdFlags|log.Lshortfile, LOG_LEVEL_DEBUG)
}

func SetLevel(level LogLevel) {
    defaultLog.SetLevel(level)
}

func Level() LogLevel {
    return defaultLog.LogLevel()
}

func Debug(args ...interface{}) {
    defaultLog.Debug(args...)
}

func Debugln(args ...interface{}) {
    defaultLog.Debugln(args...)
}

func Debugf(format string, args ...interface{}) {
    defaultLog.Debugf(format, args...)
}

func Info(args ...interface{}) {
    defaultLog.Info(args...)
}

func Infoln(args ...interface{}) {
    defaultLog.Infoln(args...)
}

func Infof(format string, args ...interface{}) {
    defaultLog.Infof(format, args...)
}

func Warn(args ...interface{}) {
    defaultLog.Warn(args...)
}

func Warnln(args ...interface{}) {
    defaultLog.Warnln(args...)
}

func Warnf(format string, args ...interface{}) {
    defaultLog.Warnf(format, args...)
}

func Error(args ...interface{}) {
    defaultLog.Error(args...)
}

func Errorln(args ...interface{}) {
    defaultLog.Error(args...)
}

func Errorf(format string, args ...interface{}) {
    defaultLog.Errorf(format, args...)
}

func Println(v ...interface{}) {
    defaultLog.Println(v...)
}

func Printf(format string, v ...interface{}) {
    defaultLog.Printf(format, v...)
}

func Print(v ...interface{}) {
    defaultLog.Print(v...)
}
