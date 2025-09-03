package deckcode

const versionByte = 1

type Data struct {
	Leader  []uint32         `json:"leader"`
	Deck    map[uint32]uint8 `json:"deck"`
	Tactics []uint32         `json:"tactics"`
}
