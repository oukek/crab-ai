package common

import (
	"context"
	"errors"
	"io"
	"os"
	"oukek/crab-ai/app/service/env"
	"path"
	"strings"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm/logger"
)

var LevelMap = map[string]logrus.Level{
	"DEBUG": logrus.DebugLevel,
	"ERROR": logrus.ErrorLevel,
	"WARN":  logrus.WarnLevel,
	"INFO":  logrus.InfoLevel,
}

var BaseLogger *logrus.Logger

func init() {
	workDir, err := os.Getwd()
	SimpleCheck(err)
	baseLogLevel := env.Get("log_base_level")
	if baseLogLevel == "" {
		baseLogLevel = "DEBUG"
	}
	var std io.Writer
	if env.Get("ENV") != "prod" {
		std = os.Stdout
	}
	BaseLogger, err = New(workDir+"/runtime", "base.log", baseLogLevel, std, 0)
	SimpleCheck(err)
}

// New 创建 @filePth: 如果路径不存在会创建 @fileName: 如果存在会被覆盖  @std: os.stdout/stderr 标准输出和错误输出
func New(filePath string, fileName string, level string, std io.Writer, count uint) (*logrus.Logger, error) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		if err := os.MkdirAll(filePath, 755); err != nil {
			return nil, err
		}
	}
	fn := path.Join(filePath, fileName)

	logger := logrus.New()

	//timeFormatter := &logrus.TextFormatter{
	//	FullTimestamp:   true,
	//	TimestampFormat: "2006-01-02 15:04:05.999999999",
	//}
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05.999999999",
	}) // 设置日志格式为json格式

	level = strings.ToUpper(level)
	if logLevel, ok := LevelMap[level]; !ok {
		return nil, errors.New("log level not found")
	} else {
		logger.SetLevel(logLevel)
	}

	//logger.SetFormatter(timeFormatter)

	if 0 == count {
		count = 90 // 0的话则是默认保留90天
	}
	logFd, err := rotatelogs.New(
		fn+".%Y-%m-%d",
		// rotatelogs.WithLinkName(fn),
		//rotatelogs.WithMaxAge(time.Duration(24*count)*time.Hour),
		rotatelogs.WithRotationTime(time.Duration(24)*time.Hour),
		rotatelogs.WithRotationCount(count),
	)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = logFd.Close() // don't need handle error
	}()

	if nil != std {
		logger.SetOutput(io.MultiWriter(logFd, std)) // 设置日志输出
	} else {
		logger.SetOutput(logFd) // 设置日志输出
	}
	// logger.SetReportCaller(true)   // 测试环境可以开启，生产环境不能开，会增加很大开销
	return logger, nil
}

type GormLogger struct {
	Log *logrus.Logger
}

func (l *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	return l // 简单实现，按需扩展
}
func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	l.Log.Infof(msg, data...)
}
func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	l.Log.Warnf(msg, data...)
}
func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	l.Log.Errorf(msg, data...)
}
func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()
	if err != nil {
		l.Log.WithFields(logrus.Fields{
			"elapsed": elapsed,
			"rows":    rows,
			"sql":     sql,
		}).Error(err)
	} else {
		l.Log.WithFields(logrus.Fields{
			"elapsed": elapsed,
			"rows":    rows,
			"sql":     sql,
		}).Info("SQL executed")
	}
}
