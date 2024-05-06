package provider

import (
	"fmt"
	"io"
	"os"
	"path"

	"login-meta-jatis/provider/dailylogger"
	"login-meta-jatis/util"

	"github.com/sirupsen/logrus"
)

type LogType int

const (
	AppLog = iota
	MongoLog
)

type ILogger interface {
	Infof(logType LogType, format string, args ...interface{})
	Errorf(logType LogType, format string, args ...interface{})
	Debugf(logType LogType, format string, args ...interface{})
	WithFields(logType LogType, fields logrus.Fields) *logrus.Entry
}

type logrusLogger struct {
	appLog   *logrus.Logger
	mongoLog *logrus.Logger
}

func NewLogger() *logrusLogger {
	// workingDirectory, _ := os.Getwd()
	// logDirectory := path.Join(workingDirectory, "log")
	appInfoLogFile := path.Join(util.Configuration.Logger.Dir, "info", fmt.Sprintf("%s.app.info.log", util.Configuration.Logger.FileName))
	appErrorLogFile := path.Join(util.Configuration.Logger.Dir, "error", fmt.Sprintf("%s.app.error.log", util.Configuration.Logger.FileName))
	mongoInfoLogFile := path.Join(util.Configuration.Logger.Dir, "mongo", fmt.Sprintf("%s.mongo.info.log", util.Configuration.Logger.FileName))
	mongoErrorLogFile := path.Join(util.Configuration.Logger.Dir, "mongoerror", fmt.Sprintf("%s.mongo.error.log", util.Configuration.Logger.FileName))

	appLog := logrus.New()
	mongoLog := logrus.New()

	maxAge := util.Configuration.Logger.MaxAge
	maxBackups := util.Configuration.Logger.MaxBackups
	maxSize := util.Configuration.Logger.MaxSize
	compress := util.Configuration.Logger.Compress
	localTime := util.Configuration.Logger.LocalTime

	formatter := &logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05.000",
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime: "@time",
		},
	}

	appLog.SetFormatter(formatter)
	mongoLog.SetFormatter(formatter)

	// Send info and debug logs to stdout
	appLog.AddHook(&WriterHook{
		Writer: dailylogger.NewDailyRotateLogger(appInfoLogFile, maxSize, maxBackups, maxAge, localTime, compress),
		LogLevels: []logrus.Level{
			logrus.InfoLevel,
			logrus.DebugLevel,
		},
	})

	// Send logs with level higher than warning to stderr
	appLog.AddHook(&WriterHook{
		Writer: dailylogger.NewDailyRotateLogger(appErrorLogFile, maxSize, maxBackups, maxAge, localTime, compress),
		LogLevels: []logrus.Level{
			logrus.PanicLevel,
			logrus.FatalLevel,
			logrus.ErrorLevel,
			logrus.WarnLevel,
		},
	})

	// Send info and debug logs to stdout
	mongoLog.AddHook(&WriterHook{
		Writer: dailylogger.NewDailyRotateLogger(mongoInfoLogFile, maxSize, maxBackups, maxAge, localTime, compress),
		LogLevels: []logrus.Level{
			logrus.InfoLevel,
			logrus.DebugLevel,
		},
	})

	// Send logs with level higher than warning to stderr
	mongoLog.AddHook(&WriterHook{
		Writer: dailylogger.NewDailyRotateLogger(mongoErrorLogFile, maxSize, maxBackups, maxAge, localTime, compress),
		LogLevels: []logrus.Level{
			logrus.PanicLevel,
			logrus.FatalLevel,
			logrus.ErrorLevel,
			logrus.WarnLevel,
		},
	})

	return &logrusLogger{appLog: appLog, mongoLog: mongoLog}
}

func (l *logrusLogger) Infof(logType LogType, format string, args ...interface{}) {
	logger := l.checkType(logType)
	logger.Infof(format, args...)
}

func (l *logrusLogger) Errorf(logType LogType, format string, args ...interface{}) {
	logger := l.checkType(logType)
	logger.Errorf(format, args...)
}

func (l *logrusLogger) Debugf(logType LogType, format string, args ...interface{}) {
	logger := l.checkType(logType)
	logger.Debugf(format, args...)
}

func (l *logrusLogger) WithFields(logType LogType, fields logrus.Fields) *logrus.Entry {
	logger := l.checkType(logType)
	return logger.WithFields(fields)
}

func (l *logrusLogger) checkType(logType LogType) *logrus.Logger {
	var logger *logrus.Logger

	if logType == AppLog {
		logger = l.appLog
	} else {
		logger = l.mongoLog
	}

	return logger
}

// WriterHook is a hook that writes logs of specified LogLevels to specified Writer
type WriterHook struct {
	Writer    io.Writer
	LogLevels []logrus.Level
}

// Fire will be called when some logging function is called with current hook
// It will format log entry to string and write it to appropriate writer
func (hook *WriterHook) Fire(entry *logrus.Entry) error {
	line, err := entry.String()
	if err != nil {
		return err
	}
	_, err = hook.Writer.Write([]byte(line))
	return err
}

// Levels define on which log levels this hook would trigger
func (hook *WriterHook) Levels() []logrus.Level {
	return hook.LogLevels
}

func InitLogDir() {
	// workingDirectory, err := os.Getwd()
	// if err != nil {
	// 	panic(err)
	// }

	workingDirectory := util.Configuration.Logger.Dir

	logDirectory := path.Join(workingDirectory)
	if _, err := os.Stat(logDirectory); os.IsNotExist(err) {
		if err := util.CreateDirectory(logDirectory); err != nil {
			panic(err)
		}
	}

	infoLogDirectory := path.Join(logDirectory, "info")
	if _, err := os.Stat(infoLogDirectory); os.IsNotExist(err) {
		if err := util.CreateDirectory(infoLogDirectory); err != nil {
			panic(err)
		}
	}

	errorLogDirectory := path.Join(logDirectory, "error")
	if _, err := os.Stat(errorLogDirectory); os.IsNotExist(err) {
		if err := util.CreateDirectory(errorLogDirectory); err != nil {
			panic(err)
		}
	}

	mongoLogDirectory := path.Join(logDirectory, "mongo")
	if _, err := os.Stat(mongoLogDirectory); os.IsNotExist(err) {
		if err := util.CreateDirectory(mongoLogDirectory); err != nil {
			panic(err)
		}
	}

	mongoErrorLogDirectory := path.Join(logDirectory, "mongoerror")
	if _, err := os.Stat(mongoErrorLogDirectory); os.IsNotExist(err) {
		if err := util.CreateDirectory(mongoErrorLogDirectory); err != nil {
			panic(err)
		}
	}
}
