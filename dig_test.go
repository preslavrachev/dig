package dig

import (
	"encoding/json"
	"github.com/tidwall/gjson"
	"log"
	"testing"
)

func BenchmarkVsGJson(b *testing.B) {
	const jsonStr = `{"a": {"b": {"c": {"d": {"e": {"f": "g"}}}}}}`

	sourceMap := make(map[string]interface{})
	err := json.Unmarshal([]byte(jsonStr), &sourceMap)
	if err != nil {
		log.Println(err)
	}
	m := NewMap(sourceMap)

	b.Run("dig", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			m.GetValue("a.b.c.d.e.f")
		}
	})

	b.Run("gjson", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			gjson.Get(jsonStr, "a.b.c.d.e.f")
		}
	})
}
