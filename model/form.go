package model

import (
	"encoding/json"
	"log"
)

type FormWrapper struct {
	*A `json:",omitempty"`
	*B `json:",omitempty"`
}

type A struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

type B struct {
	ID        string `json:"id"`
	Firstname string `json:"firstname"`
}

func (w *FormWrapper) UnmarshalJSON(data []byte) error {
	log.Println("UnmarshalJSON()")
	log.Println(string(data))

	type Source struct {
		SourceID string `json:"source"`
	}

	var source Source
	err := json.Unmarshal(data, &source)
	if err != nil {
		log.Println(err)
		return err
	}
	out, _ := json.Marshal(source)
	log.Println(string(out))

	switch source.SourceID {
	case "first":
		log.Println("do first thing")
		var a A
		w.A = &a
		return json.Unmarshal(data, &a)
	case "second":
		log.Println("do second thing")
		var b B
		w.B = &b
		return json.Unmarshal(data, &b)
	default:
		log.Println("do default thing")
	}

	return nil
}
