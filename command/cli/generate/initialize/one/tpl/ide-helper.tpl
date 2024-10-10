type Query
type Mutation
type Subscription

directive @external repeatable on FIELD_DEFINITION
directive @requires(fields: [String!]!) repeatable on FIELD_DEFINITION
directive @provides(fields: [String!]!) repeatable on FIELD_DEFINITION
directive @key(fields: [String!]!) repeatable on OBJECT | INTERFACE
directive @extends repeatable on OBJECT

directive @paginate(scopes: [String!]) on FIELD_DEFINITION
directive @skip(if: Boolean!) on FIELD_DEFINITION
directive @include(if: Boolean!) on FIELD_DEFINITION
directive @enum(value: Int!) on FIELD_DEFINITION