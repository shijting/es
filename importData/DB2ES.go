package importData

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic/v7"
	"github.com/shijting/es/Models"
	"github.com/shijting/es/inits"
	"log"
	"math"
	"net/http"
	"strconv"
)

// 单条记录导入到es中
func Import2ES(c *gin.Context) {
	bookList := Models.BookList{}
	page := 1
	pagesize := 10
	err := inits.GetDB().Select("book_id,book_name,book_intr,book_price1,book_price2,book_author,book_press,book_kind " +
		",if(book_date='','1970-01-01',ltrim(book_date)) as book_date").Order("book_id desc").Limit(pagesize).Offset((page - 1) * pagesize).Find(&bookList).Error
	if err != nil || len(bookList) == 0 {
		log.Fatal(err)
	}
	for _, book := range bookList {
		//b, err := json.Marshal(book)
		//if err != nil {
		//	log.Fatal(err)
		//}
		_, err = inits.GetEsClient().Index().Index("books").Id(strconv.Itoa(book.BookID)).BodyJson(book).Do(c)
		if err != nil {
			log.Fatal(err)
		}
	}

}

// 批量导入到es中(可以使用协程)
func ImportMany2ES(c *gin.Context) {
	bookList := Models.BookList{}
	page := 1
	pagesize := 500
	total := 0
	model := inits.GetDB().Model(&Models.Books{})
	model.Count(&total)
	totalPage := math.Ceil(float64(total/ pagesize))
	fmt.Println(total)
	for {
		err := model.Select("book_id,book_name,book_intr,book_price1,book_price2,book_author,book_press,book_kind " +
			",if(book_date='','1970-01-01',ltrim(book_date)) as book_date").Order("book_id desc").Limit(pagesize).Offset((page - 1) * pagesize).Find(&bookList).Error
		if err != nil || len(bookList) == 0 {
			log.Fatal(err)
		}
		client := inits.GetEsClient()
		bulk := client.Bulk()
		for _, book := range bookList {
			req := elastic.NewBulkIndexRequest()
			_ = req.Index("books").Id(strconv.Itoa(book.BookID)).Doc(book)
			bulk.Add(req)
		}
		ctx := context.Background()
		_, err = bulk.Do(ctx)
		if err != nil {
			log.Fatal(err)
		}
		page = page +1

		fmt.Println("total page:", totalPage , "--- current page", page)
		if totalPage < float64(page) {
			break
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "ok",
	})
}
func ImportMany2ESByGoroutine(c *gin.Context) {
	bookList := Models.BookList{}
	page := 1
	pagesize := 500
	total := 0
	model := inits.GetDB().Model(&Models.Books{})
	model.Count(&total)
	totalPage := math.Ceil(float64(total/ pagesize))
	fmt.Println(total)
	for {
		go func() {
			err := model.Select("book_id,book_name,book_intr,book_price1,book_price2,book_author,book_press,book_kind " +
				",if(book_date='','1970-01-01',ltrim(book_date)) as book_date").Order("book_id desc").Limit(pagesize).Offset((page - 1) * pagesize).Find(&bookList).Error
			if err != nil || len(bookList) == 0 {
				log.Fatal(err)
			}
			client := inits.GetEsClient()
			bulk := client.Bulk()
			for _, book := range bookList {
				req := elastic.NewBulkIndexRequest()
				_ = req.Index("books").Id(strconv.Itoa(book.BookID)).Doc(book)
				bulk.Add(req)
			}
			_, err = bulk.Do(c)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("total page:", totalPage , "--- current page", page)
			page = page +1
		}()
		if totalPage < float64(page) {
			break
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "ok",
	})
}
