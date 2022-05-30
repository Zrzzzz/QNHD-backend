package logging

import (
	"io"
	"path"
	"qnhd/pkg/setting"
	"time"

	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

// logrus: 日志到文件
func GinLogger() gin.HandlerFunc {
	logFilePath := setting.AppSetting.RuntimeRootPath + setting.AppSetting.GinLogSavePath
	logFileName := setting.AppSetting.LogSaveName
	// fileName拼接
	fileName := path.Join(logFilePath, logFileName)
	// // 写入文件
	// f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
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
	// 新增hook
	logger.AddHook(lfhook)

	return func(c *gin.Context) {
		startTime := time.Now()
		// 处理请求
		c.Next()
		endTime := time.Now()
		// 执行时间
		latencyTime := endTime.Sub(startTime)
		// 请求方式
		reqMethod := c.Request.Method
		// 请求路由
		reqUrl := c.Request.URL
		// 状态码
		statuCode := c.Writer.Status()
		// 请求IP
		clientIP := c.ClientIP()
		// 日志格式
		logger.Infof("| %3d | %13v | %15s | %s | %s",
			statuCode,
			latencyTime,
			clientIP,
			reqMethod,
			reqUrl,
		)
	}
}
