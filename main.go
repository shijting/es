package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/shijting/es/funs"
	"github.com/shijting/es/importData"
	"github.com/shijting/es/inits"
	"log"
	"github.com/shijting/es/Middlewares"
	"github.com/olivere/elastic/v7"
)

const (
	BooksIndex = "books"
	BooksLogsIndex = "bookslogs"
)

func init()  {
	client := inits.GetEsClient()
	// 创建 books index
	createBooksIndex(client)
	// 创建 books log index
	createBooksLogIndex(client)
	
}
func createBooksIndex (client *elastic.Client) {
	exists, err :=client.IndexExists(BooksIndex).Do(context.Background())
	if err != nil {
		// Handle error
		panic(err)
	}
	// 创建索引
	if ! exists {
		log.Println("creating books mapping...")
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
		createIndex, err :=client.CreateIndex(BooksLogsIndex).Body(mapping).Do(context.Background())
		if err != nil {
			// Handle error
			panic(err)
		}
		if !createIndex.Acknowledged {
			// Not acknowledged
		}
	}
}
// 创建 books log index
func createBooksLogIndex (client *elastic.Client) {
	exists, err :=client.IndexExists(BooksLogsIndex).Do(context.Background())
	if err != nil {
		// Handle error
		panic(err)
	}
	if !exists {
		log.Println("creating bookslogs mapping...")
		mapping := `{
			"mappings": {
			  "properties": {
				"ip":    { "type": "text" },
				"status":    { "type": "integer"},  
				"duration":   { "type": "integer"},  
				"method":   { "type": "keyword"},  
				"url":   { "type": "text"},
				"time":   { "type": "date"},
				"level":   { "type": "keyword"},
				"msg":  { "type": "text"}
			  }
			}
		  }`
		  createIndex, err :=client.CreateIndex(BooksLogsIndex).Body(mapping).Do(context.Background())
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
	
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(Middlewares.LogMiddleware())
	r.GET("/import", importData.ImportMany2ES)
	// 精确查找
	r.GET("/search-one/:author", funs.BooksByField)
	// 多条件查找
	r.GET("/search", funs.GetByConditions)
	r.Run(":8080")
}
