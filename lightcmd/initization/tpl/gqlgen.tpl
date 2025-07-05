schema:
  - graph/*.graphqls
  - graph/*.graphql
  - schema/*.graphqls
  - schema/*.graphql

exec:
  package: graph
  layout: follow-schema
  dir: graph
  filename_template: "{name}.generated.go"
  worker_limit: 1000


federation:
  filename: graph/federation.go
  package: graph
  version: 2
  options:
    computed_requires: true

model:
  filename: models/models_gen.go
  package: models

resolver:
  layout: follow-schema
  dir: resolver
  package: resolver
  filename_template: "{name}.resolvers.go"
  omit_template_comment: false

# Optional: turn on to skip generation of ComplexityRoot struct content and Complexity function
# omit_complexity: false

# Optional: set to speed up generation time by not performing a final validation pass.
# skip_validation: true

# Optional: set to skip running `go mod tidy` when generating server code
# skip_mod_tidy: true

# Optional: if this is set to true, argument directives that
# decorate a field with a null value will still be called.
#
# This enables argumment directives to not just mutate
# argument values but to set them even if they're null.
call_argument_directives_with_null: true

models:
  ID:
    model:
      - gitlab.staticoft.com/go-plugins/lightcmd/scalars.Uint
  Int:
    model:
      - github.com/99designs/gqlgen/graphql.Int
  DeletedAt:
    model: gitlab.staticoft.com/go-plugins/lightcmd/scalars.DeletedAt

directives:
  searchable:
    skip_runtime: true
  unique:
    skip_runtime: true
  index:
    skip_runtime: true
  varchar:
    skip_runtime: true
  text:
    skip_runtime: true
  longtext:
    skip_runtime: true
  default:
    skip_runtime: true
  loader:
    skip_runtime: true
  gorm:
    skip_runtime: true
  auth:
  hidden:
  own: