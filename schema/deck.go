package schema

import (
	"encoding/json"
	"os"

	"github.com/konojunya/compression-algorithm/deckcode"
)

func LoadData() (deckcode.Data, error) {
	data, err := os.ReadFile("data.json")
	if err != nil {
		return deckcode.Data{}, err
	}

	var d deckcode.Data
	json.Unmarshal(data, &d)

	return d, nil
}
