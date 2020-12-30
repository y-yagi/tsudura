// +build linux

package log

import "log/syslog"

type Logger struct {
	syslog *syslog.Writer
}

func NewLogger(app string) (*Logger, error) {
	syslog, err := syslog.New(syslog.LOG_NOTICE|syslog.LOG_USER, app)
	if err != nil {
		return nil, err
	}

	l := &Logger{syslog: syslog}
	return l, nil
}

func (logger *Logger) Close() {
	logger.syslog.Close()
}

func (logger *Logger) Info(msg string) {
	logger.syslog.Info(msg)
}

func (logger *Logger) Warning(msg string) {
	logger.syslog.Warning(msg)
}

func (logger *Logger) Error(msg string) {
	logger.syslog.Err(msg)
}
