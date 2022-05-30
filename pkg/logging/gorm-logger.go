package logging

import (
	"io"
	"path"
	"qnhd/pkg/setting"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

//定义自己的Writer
type GormWriter struct {
	mlog *logrus.Logger
}

//实现gorm/logger.Writer接口
func (m *GormWriter) Printf(format string, v ...interface{}) {
	m.mlog.Printf(format, v...)
}

func GormLogger() *GormWriter {
	logFilePath := setting.AppSetting.RuntimeRootPath + setting.AppSetting.GormLogSavePath
	logFileName := setting.AppSetting.LogSaveName
	// fileName拼接
	fileName := path.Join(logFilePath, logFileName)
	// 写入文件
	// f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, fs.ModeAppend)
	// if err != nil {
	// 	fmt.Println("err", err)
	// }
	// 实例化
	logger := logrus.New()
	// 设置输出
	logger.Out = io.Discard
	// 设置日志级别
	logger.SetLevel(logrus.DebugLevel)
	// 设置rotatelogs
	logWriter, _ := rotatelogs.New(
		// 分割后的文件名称
		fileName+".%Y%m%d.log",
		// 生成软链，指向最新的日志文件
		rotatelogs.WithLinkName(fileName),
		// 设置最长保存时间
		rotatelogs.WithMaxAge(7*24*time.Hour),
		// 设置日志切割间隔时间
		rotatelogs.WithRotationTime(24*time.Hour),
	)
	writeMap := lfshook.WriterMap{
		logrus.InfoLevel:  logWriter,
		logrus.FatalLevel: logWriter,
		logrus.DebugLevel: logWriter,
		logrus.WarnLevel:  logWriter,
		logrus.ErrorLevel: logWriter,
		logrus.PanicLevel: logWriter,
	}
	lfhook := lfshook.NewHook(writeMap, &logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	logger.AddHook(lfhook)
	return &GormWriter{mlog: logger}
}
