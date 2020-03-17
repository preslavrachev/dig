package dig

import (
	"errors"
)

type dict map[string]interface{}
type cache map[string]func(...interface{}) interface{}
type getterSetterFunc func(i ...interface{}) interface{}

type mmap struct {
	source dict
	cache  cache
}

func NewMap(source dict) *mmap {
	cache := make(cache)

	recurse(source, cache, "")

	return &mmap{
		source: source,
		cache:  cache,
	}
}

func recurse(m dict, cache cache, pathPrefix string) {
	for key, value := range m {

		path := key
		if pathPrefix != "" {
			path = pathPrefix + "." + key
		}

		switch value.(type) {
		case map[string]interface{}:
			recurse(value.(map[string]interface{}), cache, path)
		case string, int, int32, int64, float32, float64:
			cache[path] = func(m dict, key string) getterSetterFunc {
				return func(i ...interface{}) interface{} {
					if i != nil {
						m[key] = i[0]
					}

					return m[key]
				}
			}(m, key)
		}
	}
}

func (m *mmap) GetValue(path string) (interface{}, error) {
	getterFunc, ok := m.cache[path]
	if !ok {
		return nil, errors.New("the path does not exist")
	}

	return getterFunc(), nil
}

func (m *mmap) SetValue(path string, value interface{}) error {
	setterFunc, ok := m.cache[path]
	if !ok {
		return errors.New("the path does not exist")
	}

	setterFunc(value)
	return nil
}

func (m *mmap) Source() dict {
	return m.source
}
