package dig

import (
	"errors"
	"strconv"
)

type enumerable interface {
	getAt(key string) interface{}
	setAt(key string, value interface{})
}

type dict map[string]interface{}

func (d dict) getAt(key string) interface{} {
	return d[key]
}

func (d dict) setAt(key string, value interface{}) {
	d[key] = value
}

type slice []interface{}
type cache map[string]func(...interface{}) interface{}
type getterSetterFunc func(i ...interface{}) interface{}

type mmap struct {
	source dict
	cache  cache
}

func NewMap(source dict) *mmap {
	cache := make(cache)

	recurseDict(source, cache, "")

	return &mmap{
		source: source,
		cache:  cache,
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

func (m *mmap) PropertyPaths() []string {
	var pp []string
	for key := range m.cache {
		pp = append(pp, key)
	}

	return pp
}

func recurseDict(m dict, cache cache, pathPrefix string) {
	for key, value := range m {

		path := key
		if pathPrefix != "" {
			path = pathPrefix + "." + key
		}

		process(value, cache, path, m, key)
	}
}

func process(value interface{}, cache cache, path string, m enumerable, key string) {
	switch value.(type) {
	case map[string]interface{}:
		recurseDict(value.(map[string]interface{}), cache, path)
	case []interface{}:
		recurseSlice(value.([]interface{}), cache, path)
	case string, int, int32, int64, float32, float64:
		cache[path] = func(m enumerable, key string) getterSetterFunc {
			return func(i ...interface{}) interface{} {
				if i != nil {
					m.setAt(key, i[0])
				}

				return m.getAt(key)
			}
		}(m, key)
	}
}

func recurseSlice(m slice, cache cache, pathPrefix string) {
	for idx, value := range m {
		strIdx := strconv.Itoa(idx)
		path := strIdx
		if pathPrefix != "" {
			path = pathPrefix + "." + strIdx
		}

		process(value, cache, path, nil, strIdx)
	}
}

func (m *mmap) Source() dict {
	return m.source
}
