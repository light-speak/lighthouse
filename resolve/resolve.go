package resolve

var R Resolve

type Resolve interface {
	IsResolver() bool
}
