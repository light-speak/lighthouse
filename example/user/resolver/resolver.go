package resolver

type Resolver struct {
	Www string
}

func (r *Resolver) IsResolver() bool {
	return true
}
