directive @goModel(
    model: String
    models: [String!]
    forceGenerate: Boolean
) on OBJECT | INPUT_OBJECT | SCALAR | ENUM | INTERFACE | UNION

directive @goField(
    forceResolver: Boolean
    name: String
    omittable: Boolean
) on INPUT_FIELD_DEFINITION | FIELD_DEFINITION

directive @goTag(
    key: String!
    value: String
) on INPUT_FIELD_DEFINITION | FIELD_DEFINITION

directive @inject(field: String!, target: String!) on FIELD_DEFINITION

scalar Time

interface BaseModel {
    id: ID!
    createdAt: Time!
    updatedAt: Time!
}

interface BaseModelSoftDelete {
    id: ID!
    createdAt: Time!
    updatedAt: Time!
    deletedAt: Time
}

# 参数Eq，是否相等
directive @eq(key: String) on ARGUMENT_DEFINITION
# 参数Scope，会自动生成scope func，需要自己写逻辑
directive @scope(name: String!) on ARGUMENT_DEFINITION
# 第几页
directive @page on ARGUMENT_DEFINITION
# 每页数量
directive @size on ARGUMENT_DEFINITION

# 创建
directive @create on FIELD_DEFINITION
# 更新，要带上ID
directive @update on FIELD_DEFINITION
# 自定义Resolve，只要是准备自己写逻辑，都要附加这个注解
directive @resolve on FIELD_DEFINITION
# 获取第一条
directive @first on FIELD_DEFINITION
# 获取列表
directive @all(scopes: [String]) on FIELD_DEFINITION

# 统计数量
directive @count(model: String!, scopes: [String]) on FIELD_DEFINITION
# 求和
directive @sum(
    model: String!
    column: String!
    scopes: [String]
) on FIELD_DEFINITION

# 登录验证, 直接判断是否登录
# message: 自定义未登录提示信息
directive @auth(message: String) on FIELD_DEFINITION

# 排序
directive @orderBy(
    column: String!
    direction: SortDirection = ASC
) on FIELD_DEFINITION

# 缓存结果
directive @cache(ttl: Int!, key: String) on FIELD_DEFINITION

enum SortDirection {
    ASC
    DESC
}

enum SearchableType {
    # 用于全文搜索，适合文章内容、描述等
    TEXT
    # 用于精确匹配，适合用户ID、状态等
    KEYWORD
    # 用于长整型数值，如时间戳
    LONG
    # 用于整型数值，如年龄、数量
    INTEGER
    # 用于短整型数值，如小范围的枚举值
    SHORT
    # 用于字节型数值，如标志位
    BYTE
    # 用于双精度浮点数，如精确的金融计算
    DOUBLE
    # 用于单精度浮点数，如一般的科学计算
    FLOAT
    # 用于半精度浮点数，如简单的图形处理
    HALF_FLOAT
    # 用于可缩放的浮点数，如需要精确控制小数位的金额
    SCALED_FLOAT
    # 用于日期时间，支持多种日期格式，默认是 ISO 8601 格式
    DATE
    # 用于布尔值，表示 true 或 false
    BOOLEAN
    # 用于存储 IPv4 或 IPv6 地址
    IP
}

enum SearchableAnalyzer {
    # 最大化分词，尽可能多地分出所有可能的词汇
    IK_MAX_WORD
    # 智能分词，分出比较常用的词汇
    IK_SMART
}
# 可搜索指令
# 用于标记字段为可搜索，并指定搜索相关的参数
directive @searchable(
    # 搜索类型，指定字段的数据类型
    searchableType: SearchableType
    # 索引分析器，用于创建索引时的分词
    indexAnalyzer: SearchableAnalyzer = IK_MAX_WORD
    # 搜索分析器，用于搜索时的分词
    searchAnalyzer: SearchableAnalyzer = IK_SMART
) on FIELD_DEFINITION


