Dig lets you access and modify property values in deeply nested, **unstructured** maps, using dot-separated paths:

```golang

source := make(map[string]interface{})
data := []byte(`{"a": "b", "c": {"a": "b"}}`)
json.Unmarshal(data, &source)


d := dig.NewMap(source)
log.Printf("%+v", s.Source())
// map[a:b c:map[a:b]]

s.SetValue("c.a", "42")
s.SetValue("a", "1000")
log.Printf("%+v", s.Source())
// map[a:1000 c:map[a:42]]

b, err := s.GetValue("c.a")
log.Println(b.(string))
// 42
```

**NOTE: Still in development and mainly for educational purposes. Use with caution!**

# Why unstructured?

Most programming languages use (some sort of) structs for structured and maps for unstructured data. While structs are a great choice when one knows the layout of the incoming data, this may often not be the case. This is where generic maps hit the stage.

# Motivation

One of the projects I worked on, was a data transformation pipeline that allowed users of the pipeline to send raw data to it, along with series of tiny rule mappings. Each rule mapping was a key-value pair, where both the key and the value would be tuples of the kind:

```
(nested.dotted.path.propertyToFilterUpon: filterValue), (nested.dotted.path.propertyToSetOrReplaceValueOf: newValue)
```

The first item of the tuple would be a dotted string, representing a nested path inside the data, e.g.:

```
location.address.postalCode
```

The value in each would be used to filter incoming data upon, or set the property under the given path, respectively. Key requirements for the pipeline were:

- agnostic to the data structure being passed to it
- resilient to change
- easy to configure (as in the example above), even by non-programmers

# Advantages of dig

Upon creating a dig map, an index gets created under the hood, traversing each key until it reaches value, which cannot be traversed further (no maps, or slices).
An index may look like this:

```
a -> val ptr,
a.b -> val.ptr
a.b.c -> val.ptr
a.b.d -> val.ptr
a.f -> val.ptr
// etc ...
```

Successing retrieval or value replacement with a nested path, is made using the index alone.

Creating the index is a slight overhead at the beginning, but it dramatically speeds up the look up of deeply nested paths, as it can be seen on the benchmark resutls below.

# Performance Benchmarks

Listed below is the (after-index) performance of dig, when compared to manually traversing the path in a regular Go map. Lastly, [GJSON](https://github.com/tidwall/gjson) - a popular JSON de/serialization library with similar capabilities gets thrown into the mix. As it can be seen, dig's reading speed is significantly higher.

```
Test                                            ns. per operation
---
BenchmarkPerf/dig-8                             29.1 ns/op
BenchmarkPerf/go-map-8                          64.6 ns/op
BenchmarkPerf/gjson-8                           201 ns/op
```
