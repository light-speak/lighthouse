package service

import (
	"context"
	"fmt"
	"net"
	"search/kitex_gen/search"
	"search/kitex_gen/search/searchservice"
	"time"

	"github.com/cloudwego/kitex/server"
	"github.com/light-speak/lighthouse/log"
)

func RunServer(host string, esUrl string) {
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%s", host))
	if err != nil {
		log.Error("解析TCP地址失败: %v", err)
	}

	esService, err := NewElasticsearchService(esUrl)
	if err != nil {
		log.Error("创建Elasticsearch服务失败: %v", err)
	}

	impl := NewSearchServiceImpl(esService)

	svr := searchservice.NewServer(impl,
		server.WithServiceAddr(addr),
		server.WithReadWriteTimeout(time.Second*60),
	)

	err = svr.Run()
	if err != nil {
		log.Error("服务运行错误: %v", err)
	}
}

// NewSearchServiceImpl 创建 SearchService 实现
func NewSearchServiceImpl(esService *ElasticsearchService) *SearchServiceImpl { // 修改返回类型
	return &SearchServiceImpl{es: esService}
}

// SearchServiceImpl 实现 SearchService 接口
type SearchServiceImpl struct {
	es *ElasticsearchService
}

// 实现 searchservice.SearchService 接口的所有方法
func (s *SearchServiceImpl) CreateIndex(ctx context.Context, indexName string, mapping string) (err error) {
	return s.es.CreateIndex(ctx, indexName, mapping)
}

func (s *SearchServiceImpl) IndexDocument(ctx context.Context, indexName string, doc *search.Document) (err error) {
	log.Info("IndexDocument: %v", doc.Id)
	return s.es.IndexDocument(ctx, indexName, doc)
}

func (s *SearchServiceImpl) UpdateDocument(ctx context.Context, indexName string, doc *search.Document) (err error) {
	return s.es.UpdateDocument(ctx, indexName, doc)
}

func (s *SearchServiceImpl) DeleteIndex(ctx context.Context, indexName string) (err error) {
	return s.es.DeleteIndex(ctx, indexName)
}

func (s *SearchServiceImpl) Search(ctx context.Context, request *search.SearchRequest) (r *search.SearchResponse, err error) {
	return s.es.Search(ctx, request)
}

func (s *SearchServiceImpl) ScrollSearch(ctx context.Context, request *search.ScrollSearchRequest) (r *search.SearchResponse, err error) {
	return s.es.ScrollSearch(ctx, request)
}

func (s *SearchServiceImpl) ContinueScrollSearch(ctx context.Context, request *search.ContinueScrollSearchRequest) (r *search.SearchResponse, err error) {
	return s.es.ContinueScrollSearch(ctx, request)
}

func (s *SearchServiceImpl) CreateOrUpdate(ctx context.Context, indexName string, doc *search.Document) (err error) {
	return s.es.CreateOrUpdate(ctx, indexName, doc)
}
