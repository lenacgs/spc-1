package logging

import (
	"io"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type EchoLogger = echo.Logger

type echolog struct {
	logger
	// cenas extra
}

func (l logger) MakeEchoLogger() EchoLogger {
	return echolog{
		logger: l,
	}
}

func (l echolog) Output() io.Writer {
	return l.logger.getOutput()
}

func (l echolog) SetOutput(w io.Writer) {
	l.logger.SetOutput(w)
}

func (l echolog) Prefix() string {
	return ""
}

func (l echolog) SetPrefix(p string) {

}

func (l echolog) Level() log.Lvl {
	return log.Lvl(l.GetLevel())
}

func (l echolog) SetLevel(v log.Lvl) {
	l.logger.SetLevel(Level(v))
}

func (l echolog) SetHeader(h string) {

}

func (l echolog) Print(i ...interface{}) {
	l.logger.Info(i...)
}

func (l echolog) Printf(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}

func (l echolog) Printj(j log.JSON) {
	l.logger.Infoj(j)
}

func (l echolog) Debug(i ...interface{}) {
	l.logger.Debug(i...)
}

func (l echolog) Debugf(format string, args ...interface{}) {
	l.logger.Debugf(format, args...)
}

func (l echolog) Debugj(j log.JSON) {
	l.logger.Debugj(j)
}

func (l echolog) Info(i ...interface{}) {
	l.logger.Info(i...)
}

func (l echolog) Infof(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}

func (l echolog) Infoj(j log.JSON) {
	l.logger.Infoj(j)
}

func (l echolog) Warn(i ...interface{}) {
	l.logger.Warn(i...)
}

func (l echolog) Warnf(format string, args ...interface{}) {
	l.logger.Warnf(format, args...)
}

func (l echolog) Warnj(j log.JSON) {
	l.logger.Warnj(j)
}

func (l echolog) Error(i ...interface{}) {
	l.logger.Error(i...)
}

func (l echolog) Errorf(format string, args ...interface{}) {
	l.logger.Errorf(format, args...)
}

func (l echolog) Errorj(j log.JSON) {
	l.logger.Errorj(j)
}

func (l echolog) Fatal(i ...interface{}) {
	l.logger.Fatal(i...)
}

func (l echolog) Fatalj(j log.JSON) {
	l.logger.Fatalj(j)
}

func (l echolog) Fatalf(format string, args ...interface{}) {
	l.logger.Fatalf(format, args...)
}

func (l echolog) Panic(i ...interface{}) {
	l.logger.Panic(i...)
}

func (l echolog) Panicj(j log.JSON) {
	l.logger.Panicj(j)
}

func (l echolog) Panicf(format string, args ...interface{}) {
	l.logger.Panicf(format, args...)
}
