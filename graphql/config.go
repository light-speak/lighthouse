package graphql

type contextKey struct {
}

var (
	ContextKey = contextKey{}
)

type Context struct {
	Column   *string
	Wheres   []*Where
	Data     *map[string]interface{}
	Paginate *Paginate
}

type Where struct {
	Path  string
	Query string
	Value interface{}
}

type Paginate struct {
	Page int64
	Size int64
}
