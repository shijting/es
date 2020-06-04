package Middlewares

import(
	"time"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"os"
)

func LogMiddleware()gin.HandlerFunc  {
	// 保存日志到文件中
	var logger = logrus.New()
	logFile,err := os.OpenFile("./mylog", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		logger.Fatal(err)
	}
	logger.AddHook(NewEsHook())
	logger.Out = logFile
	return func (ctx *gin.Context)  {
		startTime := time.Now()
		ctx.Next()
		endTime := time.Now()
		delta := endTime.Sub(startTime)
		method := ctx.Request.Method
		uri := ctx.Request.RequestURI
		status := ctx.Writer.Status()
		clientIp := ctx.ClientIP()
		logger.WithField("ip",clientIp).
			WithField("status",status).
			WithField("duration",delta.Milliseconds()).
			WithField("method",method).
			WithField("url",uri).Info()
	}
}