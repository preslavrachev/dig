package dig

import (
	"encoding/json"
	"github.com/tidwall/gjson"
	"log"
	"testing"
)

func TestReadingFromANestedNode(t *testing.T) {
	const jsonStr = `{"a": {"b": {"c": {"d": {"e": {"f": "g"}}}}}}`
	sourceMap := make(map[string]interface{})
	_ = json.Unmarshal([]byte(jsonStr), &sourceMap)
	m := NewMap(sourceMap)

	result, err := m.GetValue("a.b.c.d.e.f")
	if err != nil {
		t.Errorf("Error %v", err)
	}

	expected := "g"
	if result.(string) != expected {
		t.Errorf("The result does not match the expected value: %s. Expected: %s", result.(string), expected)
	}
}

func BenchmarkVsGJson(b *testing.B) {
	const jsonStr = `{"a": {"b": {"c": {"d": {"e": {"f": "g"}}}}}}`

	sourceMap := make(map[string]interface{})
	err := json.Unmarshal([]byte(jsonStr), &sourceMap)
	if err != nil {
		log.Println(err)
	}
	m := NewMap(sourceMap)

	b.Run("map", func(b *testing.B) {
		// Simulate pure map fetching performance
		for i := 0; i < b.N; i++ {
			_ = sourceMap["a"].
			(map[string]interface{})["b"].
			(map[string]interface{})["c"].
			(map[string]interface{})["d"].
			(map[string]interface{})["e"].
			(map[string]interface{})["f"]
		}
	})

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
