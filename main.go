package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/shijting/es/funs"
	"github.com/shijting/es/importData"
	"github.com/shijting/es/inits"
	"log"
)

const (
	Index = "books"
)

func init()  {
	exists, err := inits.GetEsClient().IndexExists(Index).Do(context.Background())
	if err != nil {
		// Handle error
		panic(err)
	}
	// 创建索引
	if ! exists {
		log.Println("creating mapping...")
		mapping := `{
	"settings":{
		"number_of_shards":1,
		"number_of_replicas":0
	},
  "mappings": {
    "properties": {
      "book_id": {"type": "integer"},
      "book_name":    { 
        "type": "text",
        "analyzer": "ik_max_word",
        "search_analyzer": "ik_smart"
      },  
      "book_intr":  { 
        "type": "text",
        "analyzer": "ik_max_word",
        "search_analyzer": "ik_smart"
      }, 
      "book_price1":   { "type": "float"},  
      "book_price2":   { "type": "float"},  
      "book_author":   { "type": "keyword"},
      "book_date":   { "type": "date"},
      "book_kind":   { "type": "integer"}
    }
  }
}`
		createIndex, err :=inits.GetEsClient().CreateIndex("books").Body(mapping).Do(context.Background())
		if err != nil {
			// Handle error
			panic(err)
		}
		if !createIndex.Acknowledged {
			// Not acknowledged
		}
	}
}

func main()  {

	r := gin.Default()
	r.GET("/import", importData.ImportMany2ES)
	// 精确查找
	r.GET("/search-one/:author", funs.BooksByField)

	// 多条件查找
	r.GET("/search", funs.GetByConditions)
	r.Run(":8080")
}
