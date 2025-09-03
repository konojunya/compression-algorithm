package main

import (
	"encoding/json"
	"fmt"

	"github.com/konojunya/compression-algorithm/deckcode"
	"github.com/konojunya/compression-algorithm/schema"
)

func printJson(d interface{}) {
	json, err := json.Marshal(d)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(json))
}

func main() {
	println("test some compression algorithms")

	d, err := schema.LoadData()
	if err != nil {
		panic(err)
	}

	printJson(d)

	s, err := deckcode.Encode(d)
	if err != nil {
		panic(err)
	}

	fmt.Println(s)
	fmt.Println(len(s))

	d2, err := deckcode.Decode(s)
	if err != nil {
		panic(err)
	}

	printJson(d2)
}
