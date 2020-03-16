Dig lets you access and modify property values in deeply nested maps, using dot-separated paths:

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