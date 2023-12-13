package lighthouse

type contextKey struct {
}

var (
	ContextKey = contextKey{}
)

type Context struct {
	Wheres   []*Where
	Data     *map[string]interface{}
	Paginate *Paginate
}

type Where struct {
	Query string
	Value interface{}
}

type Paginate struct {
	Page int64
	Size int64
}
