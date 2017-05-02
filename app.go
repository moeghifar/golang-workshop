package main

import (
	"context"
	"log"

	"gopkg.in/olivere/elastic.v5"

	"fmt"

	"reflect"

	"github.com/garyburd/redigo/redis"
	"github.com/moeghifar/golang-workshop/src/util"
)

// ProductsData ...
type ProductsData struct {
	HotID     string `json:"hot_id"`
	ProductID string `json:"product_id"`
	Status    string `json:"status"`
}

// ESClient ...
var ESClient *elastic.Client

func init() {
	util.NewRedis("localhost:6379")
	// create client
	ESClient, _ = elastic.NewClient(elastic.SetURL("http://localhost:9200"))
}

func main() {
	log.Println("hella world!")
	log.Println("============")
	_ = elasticSearchDeleteIndex("products-sample")
	err := elasticSearchInsertIndex("9292", 35)
	if err == nil {
		elasticSearchExample("9292")
	}
}

func redisPingRoutine() {
	res, err := redisPing()
	if err != nil {
		log.Println(err)
	} else {
		log.Println(res)
	}
}

func redisPing() (string, error) {
	var err error

	RedisPool := util.Pool.Get()
	defer RedisPool.Close()

	resVal, err := RedisPool.Do("PING")
	returnValue, err := redis.String(resVal, err)
	return returnValue, err
}
func elasticSearchDeleteIndex(indexName string) error {
	exists, err := ESClient.IndexExists(indexName).Do(context.Background())
	if exists {
		_, err = ESClient.DeleteIndex(indexName).Do(context.Background())
	}
	if err != nil {
		return err
	}
	log.Printf("Index Deleted [%s]\n", indexName)
	return nil
}
func elasticSearchInsertIndex(hotID string, loop int) error {
	log.Println("exec [elasticSearchInsertIndex]")
	// create index
	exists, err := ESClient.IndexExists("products-sample").Do(context.Background())
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	_, err = ESClient.CreateIndex("products-sample").Do(context.Background())
	if err != nil {
		log.Println("[ERROR] while creating index ->", err)
		return err
	}
	for i := 0; i <= loop; i++ {
		log.Printf("test loop [%d]\n", i)
		// add document to index
		pids := fmt.Sprintf("10%d", i)
		prodHot := ProductsData{HotID: hotID, ProductID: pids, Status: "1"}
		_, err = ESClient.Index().
			Index("products-sample").
			Type("default").
			Id(prodHot.HotID).
			BodyJson(prodHot).
			Refresh("true").
			Do(context.Background())
		if err != nil {
			log.Println(err)
		} else {
			log.Printf("indexing [%s]\n", pids)
		}
	}
	return nil
}

func elasticSearchExample(hotID string) {
	// search with term
	termQuery := elastic.NewTermQuery("hot_id", hotID)
	searchResult, err := ESClient.Search().
		Index("products-sample").
		Query(termQuery).
		Sort("hot_id", true).
		Do(context.Background())
	if err != nil {
		log.Println(err)
	}

	fmt.Printf("Query took %d millisecnods\n", searchResult.TookInMillis)

	var ptyp ProductsData
	for _, item := range searchResult.Each(reflect.TypeOf(ptyp)) {
		if t, ok := item.(ProductsData); ok {
			fmt.Printf("Hot curated for id [%s] product id [%s] status [%s]\n", t.HotID, t.ProductID, t.Status)
		}
	}
	fmt.Println("Found total of product id", searchResult.TotalHits())
}
