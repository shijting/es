package funs

import (
	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic/v7"
	"github.com/shijting/es/inits"
	"log"
	"net/http"
)

// 将查找结果映射成struct

// 精确查找，
func BooksByField(c *gin.Context) {
	author := c.Param("author")
	log.Println(author)
	termQuery := elastic.NewTermQuery("book_author", author)
	resp, err := inits.GetEsClient().Search().Query(termQuery).Index("books").Do(c)
	if err != nil {
		log.Fatal(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "nil",
		})
	}
	c.JSON(http.StatusOK, resp)
}
// 多条件查询
func GetByConditions(ctx *gin.Context)  {
	queryList := make([]elastic.Query, 0)
	// 模糊匹配book_name
	mustQuery := elastic.NewMatchQuery("book_name", "java")
	queryList = append(queryList, mustQuery)
	// 按价格范围查询(p >= 10 and p <=99)
	rangeQuery := elastic.NewRangeQuery("book_price1").Gte(10).Lte(99)
	queryList = append(queryList, rangeQuery)
	// 按价格升序
	sortList := make([]elastic.Sorter, 0)
	sortQuery := elastic.NewFieldSort("book_price1").Asc()
	sortList = append(sortList, sortQuery)
	bl := elastic.NewBoolQuery().Must(queryList...)
	resp ,err := inits.GetEsClient().Search().Query(bl).SortBy(sortList...).Index("books").Do(ctx)
	if err != nil {
		log.Fatal(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "nil",
		})
	}
	ctx.JSON(http.StatusOK, resp)

}

