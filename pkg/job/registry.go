package job

import "reflect"

type Factory func() JobHandler

var registry = map[string]Factory{}

func RegisterJob(name string, factory Factory) {
	registry[name] = factory
}

func getFactory(name string) (Factory, bool) {
	f, ok := registry[name]
	return f, ok
}

func jobTypeName(j JobHandler) string {
	t := reflect.TypeOf(j)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t.Name()
}
