package main

import (
	"encoding/json"
	"log"

	"github.com/preslavrachev/dig"
)

func main() {

	source := make(map[string]interface{})
	data := []byte(`{"a": "b", "c": {"a": "b"}, "d": [{"e": "f"}, {"g": "h"}]}`)
	err := json.Unmarshal(data, &source)
	if err != nil {
		log.Printf("%+v", err)
	}

	log.Printf("%+v", source)
	s := dig.NewMap(source)
	log.Println(s.PropertyPaths())

	log.Printf("%+v", source["a"])
	log.Printf("%+v", s.Source())
	err = s.SetValue("c.a", "42")
	err = s.SetValue("a", "1000")
	if err != nil {
		log.Println("Error")
	}
	log.Printf("%+v", s.Source())

	b, err := s.GetValue("c.a")
	if err != nil {
		log.Println("Error")
	}
	log.Println(b.(string))

	err = s.SetValue("d.0.e", "test")
	if err != nil {
		log.Println("Error")
	}
	log.Printf("%+v", s.Source())
}
