package es

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/olivere/elastic/v7"
)

type Elastic struct {
	client  *elastic.Client
	index   string
	mapping string
	ctx     context.Context
}

func CreateElasticsearch() (*Elastic, error) {
	e := &Elastic{}
	e.SetIndex("analyze")
	err := e.Init()
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (e *Elastic) Init() error {
	var err error
	host := os.Getenv("ELASTICSEARCH_HOST")
	if host == "" {
		host = "localhost"
	}
	e.client, err = elastic.NewClient(elastic.SetURL("http://" + host + ":9200"))
	if err != nil {
		return err
	}
	e.ctx = context.Background()
	exists, err := e.client.IndexExists(e.index).Do(e.ctx)
	if err != nil {
		return err
	}

	if !exists {
		_, err := e.client.CreateIndex(e.index).BodyString(e.mapping).Do(e.ctx)
		if err != nil {
			log.Println("create es index", err)
		}
	}
	log.Println("Elasticsearch Initial Complete")
	return nil
}

func (e *Elastic) Ping() error {
	return nil
}

func (e *Elastic) Hits(key string) error {
	return nil
}

func (e *Elastic) SetIndex(index string) {
	e.index = index
}

func (e *Elastic) SetMapping(mapping string) {
	e.mapping = mapping
}

func (e *Elastic) GetAll() ([]map[string]string, error) {
	searchResult, err := e.client.
		Search().
		Index(e.index).
		Pretty(true).
		Size(30).
		Do(e.ctx)

	log.Println(searchResult.TotalHits())
	var notes []map[string]string
	if err != nil {
		log.Println("Get occured Error ", err)
		return notes, err
	}

	for _, hit := range searchResult.Hits.Hits {
		var n map[string]string
		err := json.Unmarshal(hit.Source, &n)
		if err != nil {
			log.Println("Failed Unmarshal", err)
		}

		notes = append(notes, n)
	}

	return notes, nil
}

func (e *Elastic) GetSynonym(key string) ([]map[string]string, error) {
	query := elastic.NewMultiMatchQuery(key, "Name", "Tags")
	query = query.Analyzer("korean_analyzer")
	searchResult, err := e.client.
		Search().
		Index(e.index).
		Query(query).
		Pretty(true).
		Do(e.ctx)

	log.Println(searchResult.TotalHits())
	var notes []map[string]string
	if err != nil {
		log.Println("Get Synonyms occured Error ", err)
		return notes, err
	}

	for _, hit := range searchResult.Hits.Hits {
		var n map[string]string
		err := json.Unmarshal(hit.Source, &n)
		if err != nil {
			log.Println("Failed Unmarshal", err)
		}

		notes = append(notes, n)
	}

	return notes, nil
}

func (e *Elastic) Get(key string) (string, error) {
	get, err := e.client.Get().Index(e.index).Id(key).Pretty(true).Do(e.ctx)
	if err != nil {
		log.Println("Get occured Error ", err)
		return "", err
	}

	result := string(get.Source)

	return result, nil
}

func (e *Elastic) Set(tag map[string]interface{}) error {
	now := time.Now().Format("2006-01-02")
	tag["UpdatedAt"] = now
	_, err := e.client.Index().
		Index(e.index).
		Id(tag["Name"].(string)).
		BodyJson(tag).
		Do(context.Background())

	if err != nil {
		log.Println("Set occured Error:", err)
		return err
	}
	return nil
}

func (e *Elastic) Update(key string, new_value string) error {
	_, err := e.client.Update().Index(e.index).Id(key).
		DocAsUpsert(true).
		Doc(map[string]interface{}{"Name": key, "Tags": new_value}).
		Do(e.ctx)
	if err != nil {
		log.Println("Update Error:", err)
		return err
	}
	return nil
}

func (e *Elastic) Delete(key string) error {
	_, err := e.client.Delete().Index(e.index).Id(key).Do(e.ctx)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
