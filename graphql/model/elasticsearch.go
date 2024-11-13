package model

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
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

type QBSortOrder string

const (
	QBSortOrderAsc  QBSortOrder = "ASC"
	QBSortOrderDesc QBSortOrder = "DESC"
)

type SearchQueryBuilder struct {
	indexName    string
	termFilters  []map[string]interface{}
	matchFilters []map[string]interface{}
	sorts        []map[string]QBSortOrder
}

func init() {
	if !env.LighthouseConfig.Elasticsearch.Enable {
		return
	}
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

	if res.StatusCode == 404 {
		mapping := map[string]interface{}{
			"mappings": map[string]interface{}{
				"properties": propsMapping,
			},
		}
		mappingBytes, err := json.Marshal(mapping)
		if err != nil {
			return err
		}
		res, err := s.client.Indices.Create(
			model.IndexName(),
			s.client.Indices.Create.WithBody(bytes.NewReader(mappingBytes)),
		)
		if err != nil {
			return err
		}

		if res.StatusCode != 200 {
			bodyBytes, err := io.ReadAll(res.Body)
			if err != nil {
				return err
			}
			return fmt.Errorf("create index error: %s , reason: %s", res.Status(), string(bodyBytes))
		}

		defer res.Body.Close()
	} else if res.StatusCode == 200 {
		// 获取 index 当前 mapping
		res, err := s.client.Indices.GetMapping(
			s.client.Indices.GetMapping.WithIndex(model.IndexName()),
		)
		if err != nil {
			return err
		}
		defer res.Body.Close()
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}

		existingMapping := map[string]interface{}{}
		err = json.Unmarshal(bodyBytes, &existingMapping)
		if err != nil {
			return err
		}

		if v, ok := existingMapping[model.IndexName()].(map[string]interface{})["mappings"].(map[string]interface{})["properties"].(map[string]interface{}); ok {
			existingMapping = v
		} else {
			return fmt.Errorf("es index %s not found", model.IndexName())
		}

		// 与当前 mapping 比较并合并映射，将新的字段添加到 new mapping 中
		newMapping := map[string]interface{}{}
		for k, v := range propsMapping {
			if _, ok := existingMapping[k]; !ok {
				newMapping[k] = v
			}
		}

		if len(newMapping) == 0 {
			return nil
		}

		// 更新 mapping
		mappingBytes, err := json.Marshal(map[string]interface{}{
			"properties": newMapping,
		})
		if err != nil {
			return err
		}

		updateRes, err := s.client.Indices.PutMapping(
			[]string{model.IndexName()},
			bytes.NewReader(mappingBytes),
		)
		if err != nil {
			return err
		}
		if updateRes.StatusCode != 200 {
			bodyBytes, err := io.ReadAll(updateRes.Body)
			if err != nil {
				return err
			}
			return fmt.Errorf("update index error: %s , reason: %s", updateRes.Status(), string(bodyBytes))
		}
		defer updateRes.Body.Close()
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
func (s *Searcher) QuickSearch(model SearchModel, searchString string) (*[]string, error) {

	query := fmt.Sprintf(`{
        "query": {
            "match": {
                "_all": {
                    "query": "%s",
                    "analyzer": "ik_smart"
                }
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
	ids, err := ExtractDocIDs(res)
	if err != nil {
		return nil, err
	}

	return &ids, nil
}

// index data from database to elasticsearch
func (s *Searcher) IndexDocsByModel(model SearchModel, limit int, offset int) (int, error) {
	err := s.CreateOrUpdateIndex(model)
	if err != nil {
		return 0, err
	}
	var rs []map[string]interface{}
	if err := db.Table(model.TableName()).Offset(offset).Limit(limit).Find(&rs).Error; err != nil {
		return 0, err
	}
	for _, r := range rs {
		data := model.GetSearchData(r)
		err := s.IndexDocByMap(model.IndexName(), data)
		if err != nil {
			return 0, err
		}
	}
	return len(rs), nil
}

func NewSearchQueryBuilder(indexName string) *SearchQueryBuilder {
	return &SearchQueryBuilder{
		indexName:    indexName,
		termFilters:  []map[string]interface{}{},
		matchFilters: []map[string]interface{}{},
		sorts:        []map[string]QBSortOrder{},
	}
}

func (qb *SearchQueryBuilder) WhereTerm(field string, value interface{}) *SearchQueryBuilder {
	qb.termFilters = append(qb.termFilters, map[string]interface{}{
		"term": map[string]interface{}{
			field: value,
		},
	})
	return qb
}

func (qb *SearchQueryBuilder) Fuzzy(field string, value string) *SearchQueryBuilder {
	qb.matchFilters = append(qb.matchFilters, map[string]interface{}{
		"match": map[string]interface{}{
			field: value,
		},
	})
	return qb
}

// sort by field
func (qb *SearchQueryBuilder) OrderBy(field string, order QBSortOrder) *SearchQueryBuilder {
	qb.sorts = append(qb.sorts, map[string]QBSortOrder{field: order})
	return qb
}

// build elasticsearch query code
func (qb *SearchQueryBuilder) Build() (map[string]interface{}, error) {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{},
		},
	}

	// 构建 must 子句
	mustClauses := []map[string]interface{}{}
	mustClauses = append(mustClauses, qb.termFilters...)
	mustClauses = append(mustClauses, qb.matchFilters...)

	if len(mustClauses) > 0 {
		query["query"].(map[string]interface{})["bool"].(map[string]interface{})["must"] = mustClauses
	}

	// 正确构建排序数组
	if len(qb.sorts) > 0 {
		sortArray := make([]map[string]interface{}, 0)
		for _, sort := range qb.sorts {
			for field, order := range sort {
				sortArray = append(sortArray, map[string]interface{}{
					field: map[string]interface{}{
						"order": strings.ToLower(string(order)),
					},
				})
			}
		}
		query["sort"] = sortArray
	}

	return query, nil
}

func (qb *SearchQueryBuilder) Execute() (*esapi.Response, error) {
	query, err := qb.Build()
	if err != nil {
		return nil, err
	}
	queryBytes, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}

	res, err := searcher.client.Search(
		searcher.client.Search.WithIndex(qb.indexName),
		searcher.client.Search.WithBody(bytes.NewReader(queryBytes)),
	)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func ExtractDocIDs(res *esapi.Response) ([]string, error) {
	defer res.Body.Close()

	// 检查查询是否成功
	if res.IsError() {
		return nil, fmt.Errorf("search error: %s", res.Status())
	}

	// 解析 JSON 响应
	var response map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}

	// 提取文档 ID
	var docIDs []string
	hits, ok := response["hits"].(map[string]interface{})["hits"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected response format")
	}
	for _, hit := range hits {
		if hitMap, ok := hit.(map[string]interface{}); ok {
			if id, ok := hitMap["_id"].(string); ok {
				docIDs = append(docIDs, id)
			}
		}
	}

	return docIDs, nil
}
