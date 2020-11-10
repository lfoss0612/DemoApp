package request

type Value interface {
	Validate() error
}

type Factory interface {
	NewInstance() Value
}
