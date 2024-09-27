namespace go search

// 保留 Document 结构体，因为它仍然用于索引和更新操作
struct Document {
    1: string id
    2: string content  
}

// 新增一个简化的搜索结果结构体
struct SearchHit {
    1: string id
}

struct SearchRequest {
    1: string indexName
    2: string query
    3: i32 from
    4: i32 size
}

struct ScrollSearchRequest {
    1: string indexName
    2: string query
    3: string scrollTime
}

struct ContinueScrollSearchRequest {
    1: string scrollID
    2: string scrollTime
}

// 修改 SearchResponse 结构体
struct SearchResponse {
    1: list<SearchHit> hits
    2: i64 total
    3: optional string scrollID
}

service SearchService {
    void CreateIndex(1: string indexName, 2: string mapping) 
    void IndexDocument(1: string indexName, 2: Document doc)
    void UpdateDocument(1: string indexName, 2: Document doc)
    void DeleteIndex(1: string indexName)
    SearchResponse Search(1: SearchRequest request)
    SearchResponse ScrollSearch(1: ScrollSearchRequest request)
    SearchResponse ContinueScrollSearch(1: ContinueScrollSearchRequest request)
    void CreateOrUpdate(1: string indexName, 2: Document doc)
}