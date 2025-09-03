package deckcode

import (
	"encoding/base64"
	"errors"
)

// Encode は固定幅ビットパックでエンコードし、Base64URL（ノーパディング）を返します。
func Encode(d Data) (string, error) {
	raw, err := encodeBytes(d)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(raw), nil
}

func Decode(s string) (Data, error) {
	raw, err := base64.RawURLEncoding.DecodeString(s)
	if err != nil {
		return Data{}, err
	}
	return decodeBytesV2(raw)
}

func decodeBytesV2(b []byte) (Data, error) {
	cur := 0
	if cur >= len(b) || b[cur] != versionByte {
		return Data{}, errors.New("unsupported version or empty")
	}
	cur++

	// ---- L (Leader) ----
	nL, err := readUvarint(b, &cur)
	if err != nil {
		return Data{}, err
	}
	lb, err := readUvarint(b, &cur) // ← 常に読む（N==0 でも 0 が来る想定）
	if err != nil {
		return Data{}, err
	}
	var L []uint32
	if nL > 0 {
		if cur+int(lb) > len(b) {
			return Data{}, errors.New("leader bytes: eof")
		}
		L, err = unpackIDs20(b[cur:cur+int(lb)], int(nL))
		if err != nil {
			return Data{}, err
		}
		cur += int(lb)
	} else {
		if lb != 0 {
			return Data{}, errors.New("leader length must be 0 when N==0")
		}
	}

	// ---- D (Deck) ----
	nD, err := readUvarint(b, &cur)
	if err != nil {
		return Data{}, err
	}
	db, err := readUvarint(b, &cur) // Deck も常に length を読む（既存通り）
	if err != nil {
		return Data{}, err
	}
	var deck map[uint32]uint8
	if nD > 0 {
		if cur+int(db) > len(b) {
			return Data{}, errors.New("deck bytes: eof")
		}
		ids, cs, err := unpackDeckPairs(b[cur:cur+int(db)], int(nD))
		if err != nil {
			return Data{}, err
		}
		cur += int(db)
		deck = make(map[uint32]uint8, nD)
		for i := range ids {
			deck[ids[i]] = cs[i]
		}
	} else {
		if db != 0 {
			return Data{}, errors.New("deck length must be 0 when N==0")
		}
		deck = map[uint32]uint8{}
	}

	// ---- T (Tactics) ----
	nT, err := readUvarint(b, &cur)
	if err != nil {
		return Data{}, err
	}
	tb, err := readUvarint(b, &cur) // ← 常に読む
	if err != nil {
		return Data{}, err
	}
	var T []uint32
	if nT > 0 {
		if cur+int(tb) > len(b) {
			return Data{}, errors.New("tactics bytes: eof")
		}
		T, err = unpackIDs20(b[cur:cur+int(tb)], int(nT))
		if err != nil {
			return Data{}, err
		}
		cur += int(tb)
	} else {
		if tb != 0 {
			return Data{}, errors.New("tactics length must be 0 when N==0")
		}
	}

	return Data{Leader: L, Tactics: T, Deck: deck}, nil
}
