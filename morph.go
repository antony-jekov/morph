package morph

type Morph interface {
	Struct(structPtr interface{}) error
}

type transformer struct {
}

func (t transformer) Struct(structPtr interface{}) error {
	return nil
}

func New() Morph {
	return transformer{}
}
