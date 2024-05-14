package main

import (
	"bytes"
	"encoding/json"
	"log"

	"github.com/davecgh/go-spew/spew"
)

type Result struct {
	success Success
	error   Error
}

type Error struct {
	ErrMsg string `json:"error_msg"`
}

type Success struct {
	ID       string
	Username string
}

func (s *Success) UnmarshalJSON(d []byte) error {
	var r Result
	err := json.Unmarshal(d, &r)
	if err != nil {
		log.Println("error unmarshaling json")
	}
	return err
}

type IntOrString int

func (i *IntOrString) UnmarshalJSON(d []byte) error {
	var v int
	err := json.Unmarshal(bytes.Trim(d, `"`), &v)
	*i = IntOrString(v)
	return err
}

func main() {
	data := []byte(`[123,"321"]`)
	x := make([]IntOrString, 0)
	if err := json.Unmarshal(data, &x); err != nil {
		panic(err)
	}
	spew.Dump(x)
}
