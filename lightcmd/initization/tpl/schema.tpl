scalar Time
scalar DeletedAt

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

directive @auth(msg: String) on FIELD_DEFINITION
directive @own on FIELD_DEFINITION
directive @hidden on FIELD_DEFINITION

directive @longtext on FIELD_DEFINITION
directive @text on FIELD_DEFINITION
directive @varchar(length: Int!) on FIELD_DEFINITION
directive @index(name: String) on FIELD_DEFINITION
directive @unique on FIELD_DEFINITION
directive @default(value: String!) on FIELD_DEFINITION
directive @gorm(value: String!) on FIELD_DEFINITION
directive @loader(keys: [String!], morphKey: String, unionTypes: [String!], extraKeys: [String!]) on OBJECT



input PaginationInput {
	current: Int
	pageSize: Int
}

enum SortOrder {
	ASC
	DESC
}

input SorterInput {
	field: String!
	order: SortOrder!
}

enum FilterOperator {
	EQ
	NE
	LT
	LTE
	GT
	GTE
	CONTAINS
	STARTS_WITH
	ENDS_WITH
	IN
	NIN
	NULL
	NNULL
	BETWEEN
}

input FilterInput {
	field: String!
	operator: FilterOperator!
	value: String
}