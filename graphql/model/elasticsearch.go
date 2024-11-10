package model

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/light-speak/lighthouse/env"
	"github.com/light-speak/lighthouse/log"
)

var searcher *Searcher

type Searcher struct {
	client *elasticsearch.Client
}

type SearchModel interface {
	ModelInterface
	FieldMapping() map[string]interface{}
	SearchId() int64
	IndexName() string
	GetSearchData(mapData ...map[string]interface{}) map[string]interface{}
}

func InitSearch() {
	c, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{fmt.Sprintf("http://%s:%s", env.LighthouseConfig.Elasticsearch.Host, env.LighthouseConfig.Elasticsearch.Port)},
		Username:  env.LighthouseConfig.Elasticsearch.User,
		Password:  env.LighthouseConfig.Elasticsearch.Password,
	})
	if err != nil {
		log.Error().Err(err).Msg("failed to connect elasticsearch")
	}
	searcher = &Searcher{
		client: c,
	}
}

func GetSearcher() *Searcher {
	return searcher
}

// create index if not exists, and update mapping if exists
func (s *Searcher) CreateOrUpdateIndex(model SearchModel) error {
	res, err := s.client.Indices.Exists([]string{model.IndexName()})
	if err != nil {
		return err
	}

	propsMapping := model.FieldMapping()
	mapping := map[string]interface{}{
		"mappings": map[string]interface{}{
			"properties": propsMapping,
		},
	}
	mappingBytes, err := json.Marshal(mapping)
	if err != nil {
		return err
	}

	if res.StatusCode == 404 {
		res, err := s.client.Indices.Create(
			model.IndexName(),
			s.client.Indices.Create.WithBody(bytes.NewReader(mappingBytes)),
		)
		if err != nil {
			return err
		}
		defer res.Body.Close()
	} else if res.StatusCode == 200 {
		res, err := s.client.Indices.PutMapping(
			[]string{model.IndexName()},
			bytes.NewReader(mappingBytes),
		)
		if err != nil {
			return err
		}
		defer res.Body.Close()
	}
	return nil
}

// index data to elasticsearch
func (s *Searcher) IndexDoc(model SearchModel) error {
	data := model.GetSearchData()
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	res, err := s.client.Index(
		model.IndexName(),
		bytes.NewReader(dataBytes),
		s.client.Index.WithDocumentID(fmt.Sprintf("%d", model.SearchId())),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}

func (s *Searcher) IndexDocByMap(index string, mapData map[string]interface{}) error {
	dataBytes, err := json.Marshal(mapData)
	if err != nil {
		return err
	}

	res, err := s.client.Index(
		index,
		bytes.NewReader(dataBytes),
		s.client.Index.WithDocumentID(fmt.Sprintf("%d", mapData["id"])),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}

// update data to elasticsearch
func (s *Searcher) UpdateDoc(model SearchModel) error {
	data := model.GetSearchData()
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	res, err := s.client.Update(
		model.IndexName(),
		fmt.Sprintf("%d", model.SearchId()),
		bytes.NewReader(dataBytes),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}

// delete data from elasticsearch
func (s *Searcher) DeleteDoc(model SearchModel) error {
	res, err := s.client.Delete(
		model.IndexName(),
		fmt.Sprintf("%d", model.SearchId()),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}

// search data from elasticsearch, response doc ids
func (s *Searcher) Search(model SearchModel, searchString string) (*[]string, error) {

	query := fmt.Sprintf(`{
        "query": {
            "query_string": {
                "query": "%s"
            }
        },
        "_source": false
    }`, searchString)

	req := esapi.SearchRequest{
		Index: []string{model.IndexName()},
		Body:  strings.NewReader(query),
	}

	res, err := req.Do(context.Background(), s.client)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("search error: %s", res.Status())
	}

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, err
	}

	ids := []string{}
	hits, ok := r["hits"].(map[string]interface{})["hits"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected response format")
	}
	for _, hit := range hits {
		id, ok := hit.(map[string]interface{})["_id"].(string)
		if !ok {
			return nil, fmt.Errorf("unexpected response format")
		}
		ids = append(ids, id)
	}

	return &ids, nil
}

// index data from database to elasticsearch
func (s *Searcher) IndexDocsByModel(model SearchModel, limit int, offset int) (int, error) {
	s.CreateOrUpdateIndex(model)
	var rs []map[string]interface{}
	if err := db.Table(model.TableName()).Offset(offset).Limit(limit).Find(&rs).Error; err != nil {
		return 0, err
	}
	for _, r := range rs {
		s.IndexDocByMap(model.IndexName(), r)
	}
	return len(rs), nil
}
