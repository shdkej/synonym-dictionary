package dictionary

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/olivere/elastic/v7"
)

type Elastic struct {
	client *elastic.Client
	index  string
	ctx    context.Context
}

func (e *Elastic) Init() error {
	var err error
	host := os.Getenv("ELASTICSEARCH_HOST")
	if host == "" {
		host = "localhost"
	}
	e.client, err = elastic.NewClient(elastic.SetURL("http://" + host + ":9200"))
	if err != nil {
		log.Fatal(err)
	}
	e.index = "analyze"
	e.ctx = context.Background()
	exists, err := e.client.IndexExists(e.index).Do(e.ctx)
	if err != nil {
		log.Println("check es index exist ", err)
		return err
	}

	if !exists {
		mapping := `{
  "settings": {
    "analysis": {
      "tokenizer": {
        "nori_user_dict": {
          "type": "nori_tokenizer",
          "decompound_mode": "mixed",
          "user_dictionary": "userdict.txt"
        }
      },
      "analyzer": {
        "korean_analyzer": {
          "filter": [
            "pos_filter_speech", "nori_readingform",
            "lowercase", "synonym", "remove_duplicates"
          ],
          "tokenizer": "nori_user_dict"
        }
      },
      "filter": {
        "synonym" : {
          "type" : "synonym_graph",
          "synonyms_path" : "synonyms.txt"
        },
        "pos_filter_speech": {
          "type": "nori_part_of_speech",
          "stoptags": [
            "E", "J", "SC", "SE", "SF", "SP", "SSC", "SSO", "SY", "VCN", "VCP",
            "VSV", "VX", "XPN", "XSA", "XSN", "XSV"
          ]
        }
      }
    }
  }
}`
		_, err := e.client.CreateIndex(e.index).BodyString(mapping).Do(e.ctx)
		if err != nil {
			log.Fatal("create es index", err)
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

func (e *Elastic) SetIndex(index string) string {
	e.index = index
	return index
}

func (e *Elastic) GetAll() ([]Tag, error) {
	searchResult, err := e.client.
		Search().
		Index(e.index).
		Pretty(true).
		Size(30).
		Do(e.ctx)

	log.Println(searchResult.TotalHits())
	var notes []Tag
	if err != nil {
		log.Println("Get occured Error ", err)
		return notes, err
	}

	for _, hit := range searchResult.Hits.Hits {
		var n Tag
		err := json.Unmarshal(hit.Source, &n)
		if err != nil {
			log.Println("Failed Unmarshal", err)
		}

		notes = append(notes, n)
	}

	return notes, nil
}

func (e *Elastic) GetSynonym(key string) ([]Tag, error) {
	query := elastic.NewMultiMatchQuery(key, "Name", "Tags")
	query = query.Analyzer("korean_analyzer")
	searchResult, err := e.client.
		Search().
		Index(e.index).
		Query(query).
		Pretty(true).
		Do(e.ctx)

	log.Println(searchResult.TotalHits())
	var notes []Tag
	if err != nil {
		log.Println("Get Synonyms occured Error ", err)
		return notes, err
	}

	for _, hit := range searchResult.Hits.Hits {
		var n Tag
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

func (e *Elastic) Set(tag Tag) error {
	now := time.Now().Format("2006-01-02")
	tag.UpdatedAt = now
	_, err := e.client.Index().
		Index(e.index).
		Id(tag.Name).
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
		Doc(map[string]interface{}{"Tags": new_value}).
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
