package Middlewares

import (
	"context"
	"github.com/shijting/es/inits"
	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
	"log"
	"strings"
)

type EsHook struct {

}
func NewEsHook() *EsHook  {
	return &EsHook{}
}
func(this *EsHook) Fire (entry *logrus.Entry )  error {
	data:=entry.Data
	data["time"]=entry.Time
	data["level"]=entry.Level
	data["msg"]=entry.Message
	if strings.Index(data["url"].(string),"/favicon.ico")>=0{
		return nil
	}
	client:=inits.GetEsClient()
	bulk:=client.Bulk()
	 {
		 req:=elastic.NewBulkIndexRequest()
		 req.Index("bookslogs").Doc(data)//直接插入
		 bulk.Add(req)
	 }
	_,err:=bulk.Do(context.Background())
	if err!=nil {
		log.Println(err)
	}
	return nil
}
func(this *EsHook) Levels () []logrus.Level {
		return logrus.AllLevels
}