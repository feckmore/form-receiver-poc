package model

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
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

type FormRecord struct {
	PK string `json:"pk"`
	SK string `json:"sk"`
	*A `json:",omitempty"`
	*B `json:",omitempty"`
}

func (w *FormWrapper) UnmarshalJSON(data []byte) error {
	log.Println("UnmarshalJSON()")

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

func (record *FormRecord) UnmarshalJSON(data []byte) error {
	log.Println("UnmarshalJSON()")

	var fields map[string]string
	err := json.Unmarshal(data, &fields)
	if err != nil {
		log.Println(err)
		return err
	}

	record.PK = "#TEMP#"
	record.SK = time.Now().Format(time.RFC3339)

	if pk, exists := fields["id"]; exists {
		record.PK = fmt.Sprintf("#ID#%s", pk)
		record.B = &B{}
		record.ID = pk
		record.Firstname = fields["firstname"]
	}

	if pk, exists := fields["name"]; exists {
		record.PK = fmt.Sprintf("#NAME#%s", strings.ReplaceAll(pk, " ", "_"))
		record.A = &A{}
		record.Name = pk
		record.Address = fields["address"]
	}

	return nil
}
