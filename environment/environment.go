package environment

import "maz-lang/object"

type Environment struct {
	values map[string]object.Object
	child  *Environment
}

func New() Environment {
	return Environment{
		values: make(map[string]object.Object),
	}
}

func (e *Environment) Set(name string, value object.Object) {
	e.values[name] = value
}

func (e *Environment) Get(name string) object.Object {
	if res, ok := e.values[name]; ok {
		return res
	}

	if e.child != nil {
		return e.child.Get(name)
	}

	return nil

}

func (e *Environment) Extend(env *Environment) {
	e.child = env
}
