package inits

import "github.com/olivere/elastic/v7"

func GetEsClient() *elastic.Client  {
	client,err:=elastic.NewClient(
		elastic.SetURL("http://106.53.5.146:9200/"),
		elastic.SetSniff(false),
	)

	if err!=nil{
		return nil
	}
	return client

}
