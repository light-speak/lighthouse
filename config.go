package lighthouse

type contextKey struct {
}

var (
	ContextKey = contextKey{}
)

type Context struct {
	Wheres   []*Where
	Paginate *Paginate
}

type Where struct {
	Query string
	Value interface{}
}

type Paginate struct {
	Page int
	Size int
}
