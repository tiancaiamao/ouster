package util

import (
    "fmt"
    "log"
    "os"
)

const (
    LOG_LEVEL_DEBUG = iota
    LOG_LEVEL_INFO
    LOG_LEVEL_WARN
    LOG_LEVEL_ERROR
)

var Log *DefaultLogger

type DefaultLogger struct {
    *log.Logger
    level int
}

func (l *DefaultLogger) LogLevel() int {
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

func init() {
    Log = &DefaultLogger{
        Logger: log.New(os.Stdout, "", log.LstdFlags|log.Llongfile),
        level:  LOG_LEVEL_DEBUG,
    }
}
