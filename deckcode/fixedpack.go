package deckcode

import "errors"

// パック幅（変更しやすいように定数化）
const (
	idBits    = 20 // 1..1_000_000 < 2^20
	countBits = 4  // 1..10
)

// leader/tactics: 20bit×N で詰める
func packIDs20(ids []uint32) []byte {
	var bw bitWriter
	for _, id := range ids {
		bw.write(id, idBits)
	}
	return bw.finish()
}

func unpackIDs20(b []byte, n int) ([]uint32, error) {
	br := bitReader{src: b}
	out := make([]uint32, n)
	for i := 0; i < n; i++ {
		v, err := br.read(idBits)
		if err != nil {
			return nil, err
		}
		out[i] = v
	}
	// 余剰ビットはそのまま無視（上位は0詰め）
	return out, nil
}

// deck: 1ペア = 20bit(ID) + 4bit(count-1) = 24bit（=3バイト）
func packDeckPairs(ids []uint32, counts []uint8) ([]byte, error) {
	if len(ids) != len(counts) {
		return nil, errors.New("packDeckPairs: len mismatch")
	}
	var bw bitWriter
	for i := range ids {
		id := ids[i]
		c := counts[i]
		if c < 1 || c > 10 {
			return nil, errors.New("count out of range (1..10)")
		}
		bw.write(id, idBits)
		bw.write(uint32(c-1), countBits)
	}
	return bw.finish(), nil
}

func unpackDeckPairs(b []byte, n int) ([]uint32, []uint8, error) {
	br := bitReader{src: b}
	ids := make([]uint32, n)
	cs := make([]uint8, n)
	for i := 0; i < n; i++ {
		id, err := br.read(idBits)
		if err != nil {
			return nil, nil, err
		}
		cm1, err := br.read(countBits)
		if err != nil {
			return nil, nil, err
		}
		ids[i] = id
		cs[i] = uint8(cm1) + 1
	}
	return ids, cs, nil
}
