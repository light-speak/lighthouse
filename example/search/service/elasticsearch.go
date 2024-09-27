package service

import (
	"context"
	"encoding/json"
	"errors"
	"search/kitex_gen/search"

	"github.com/light-speak/lighthouse/log"
	"github.com/olivere/elastic/v7"
)

type ElasticsearchService struct {
	client *elastic.Client
}

func NewElasticsearchService(url string) (*ElasticsearchService, error) {
	client, err := elastic.NewClient(elastic.SetURL(url))
	if err != nil {
		return nil, err
	}
	return &ElasticsearchService{client: client}, nil
}

// CreateIndex 创建索引
func (s *ElasticsearchService) CreateIndex(ctx context.Context, indexName string, mapping string) error {
	// 检查索引是否存在
	exists, err := s.client.IndexExists(indexName).Do(ctx)
	if err != nil {
		return err
	}

	if exists {
		// 如果索引存在，更新映射
		var mappingMap map[string]interface{}
		if err := json.Unmarshal([]byte(mapping), &mappingMap); err != nil {
			return err
		}
		mappings, ok := mappingMap["mappings"].(map[string]interface{})
		if !ok {
			log.Error("invalid type for mappings")
			return errors.New("invalid type for mappings")
		}
		_, err = s.client.PutMapping().Index(indexName).BodyJson(mappings).Do(ctx)
	} else {
		// 如果索引不存在，创建索引
		_, err = s.client.CreateIndex(indexName).BodyJson(mapping).Do(ctx)
	}
	return err
}

// IndexDocument 录入数据
func (s *ElasticsearchService) IndexDocument(ctx context.Context, indexName string, doc *search.Document) error {

	// 索引文档
	_, err := s.client.Index().
		Index(indexName).
		Id(doc.Id).
		BodyJson(doc.Content).
		Do(ctx)
	if err == nil {
		log.Info("数据录入成功: %v", doc.Id)
	} else {
		log.Error("数据录入失败: %v", err)
	}
	return err
}

// UpdateDocument 更新数据
func (s *ElasticsearchService) UpdateDocument(ctx context.Context, indexName string, doc *search.Document) error {

	// 更新文档
	_, err := s.client.Update().
		Index(indexName).
		Id(doc.Id).
		Doc(doc.Content).
		Do(ctx)
	return err
}

// DeleteIndex 删除索引
func (s *ElasticsearchService) DeleteIndex(ctx context.Context, indexName string) error {
	_, err := s.client.DeleteIndex(indexName).Do(ctx)
	return err
}

// Search 传统翻页搜索
func (s *ElasticsearchService) Search(ctx context.Context, request *search.SearchRequest) (*search.SearchResponse, error) {

	searchResult, err := s.client.Search().
		Index(request.IndexName).
		Query(elastic.NewQueryStringQuery(request.Query)).
		From(int(request.From)).
		Size(int(request.Size)).
		Do(ctx)

	if err != nil {
		return nil, err
	}

	hits := make([]*search.SearchHit, len(searchResult.Hits.Hits))
	for i, hit := range searchResult.Hits.Hits {
		hits[i] = &search.SearchHit{Id: hit.Id}
	}

	return &search.SearchResponse{
		Hits:  hits,
		Total: searchResult.TotalHits(),
	}, nil
}

// ScrollSearch 基于 scroll id 的查询翻页
func (s *ElasticsearchService) ScrollSearch(ctx context.Context, request *search.ScrollSearchRequest) (*search.SearchResponse, error) {
	searchResult, err := s.client.Scroll().
		Index(request.IndexName).
		Query(elastic.NewQueryStringQuery(request.Query)).
		Size(100).
		Scroll(request.ScrollTime).
		Do(ctx)

	if err != nil {
		return nil, err
	}

	return s.convertSearchResult(searchResult)
}

// ContinueScrollSearch 继续 scroll 查询
func (s *ElasticsearchService) ContinueScrollSearch(ctx context.Context, request *search.ContinueScrollSearchRequest) (*search.SearchResponse, error) {
	searchResult, err := s.client.Scroll().
		ScrollId(request.ScrollID).
		Scroll(request.ScrollTime).
		Do(ctx)

	if err != nil {
		return nil, err
	}

	return s.convertSearchResult(searchResult)
}

// CreateOrUpdate 创建或更新文档
func (s *ElasticsearchService) CreateOrUpdate(ctx context.Context, indexName string, doc *search.Document) error {
	// 解析 JSON 字符串
	var docJSON map[string]interface{}
	err := json.Unmarshal([]byte(doc.Content), &docJSON)
	if err != nil {
		return err
	}

	// 创建或更新文档
	_, err = s.client.Index().
		Index(indexName).
		Id(doc.Id).
		BodyJson(docJSON).
		Do(ctx)
	return err
}

func (s *ElasticsearchService) convertSearchResult(searchResult *elastic.SearchResult) (*search.SearchResponse, error) {
	hits := make([]*search.SearchHit, len(searchResult.Hits.Hits))
	for i, hit := range searchResult.Hits.Hits {
		hits[i] = &search.SearchHit{Id: hit.Id}
	}

	return &search.SearchResponse{
		Hits:     hits,
		Total:    searchResult.TotalHits(),
		ScrollID: &searchResult.ScrollId,
	}, nil
}
