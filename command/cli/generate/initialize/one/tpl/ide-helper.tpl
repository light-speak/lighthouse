type Query
type Mutation
type Subscription

directive @skip(if: Boolean!) on FIELD_DEFINITION
directive @include(if: Boolean!) on FIELD_DEFINITION
directive @deprecated(reason: String) on FIELD_DEFINITION

# enum 枚举
directive @enum(value: Int!) on FIELD_DEFINITION


# model 
directive @index(name: String) on FIELD_DEFINITION
directive @tag(name: String!, value: String!) on FIELD_DEFINITION
directive @defaultString(value: String!) on FIELD_DEFINITION
directive @defaultInt(value: Int!) on FIELD_DEFINITION
directive @unique on FIELD_DEFINITION
directive @model(name: String) on OBJECT
directive @softDeleteModel(name: String) on OBJECT


# query returnType
directive @paginate(scopes: [String!]) on FIELD_DEFINITION
directive @find(scopes: [String!]) on FIELD_DEFINITION
directive @first(scopes: [String!]) on FIELD_DEFINITION


# argument filter
directive @in(field: String) on ARGUMENT_DEFINITION
directive @eq(field: String) on ARGUMENT_DEFINITION
directive @neq(field: String) on ARGUMENT_DEFINITION
directive @gt(field: String) on ARGUMENT_DEFINITION
directive @gte(field: String) on ARGUMENT_DEFINITION
directive @lt(field: String) on ARGUMENT_DEFINITION
directive @lte(field: String) on ARGUMENT_DEFINITION
directive @like(field: String) on ARGUMENT_DEFINITION
directive @notIn(field: String) on ARGUMENT_DEFINITION

# relation
directive @belongsTo(relation: String, foreignKey: String, reference: String) on FIELD_DEFINITION
directive @hasMany(relation: String, foreignKey: String, reference: String) on FIELD_DEFINITION
directive @hasOne(relation: String, foreignKey: String, reference: String) on FIELD_DEFINITION
directive @morphTo(morphType: String, morphKey: String, reference: String) on FIELD_DEFINITION
directive @morphToMany(relation: String!, morphType: String, morphKey: String, reference: String) on FIELD_DEFINITION
