type Query
type Mutation
type Subscription


directive @paginate(scopes: [String!]) on FIELD_DEFINITION
directive @skip(if: Boolean!) on FIELD_DEFINITION
directive @include(if: Boolean!) on FIELD_DEFINITION
directive @enum(value: Int!) on FIELD_DEFINITION