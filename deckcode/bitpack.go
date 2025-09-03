package deckcode

import "errors"

// PackCounts4bit は count(1..10) を 4bit (count-1) に詰め、2 つで 1 バイトにします。
func PackCounts4bit(counts []uint) ([]byte, error) {
	n := len(counts)
	out := make([]byte, (n+1)/2)
	for i, c := range counts {
		if c < 1 || c > 10 {
			return nil, errors.New("count out of range (1..10)")
		}
		v := byte(c - 1) // 0..9
		if i%2 == 0 {
			out[i/2] = v & 0x0F
		} else {
			out[i/2] |= (v & 0x0F) << 4
		}
	}
	return out, nil
}

// UnpackCounts4bit は 4bit 詰めの配列を n 個の count(1..10) に戻します。
func UnpackCounts4bit(b []byte, n int) ([]uint, error) {
	if len(b) != (n+1)/2 {
		return nil, errors.New("invalid counts length")
	}
	out := make([]uint, n)
	for i := 0; i < n; i++ {
		var v byte
		if i%2 == 0 {
			v = b[i/2] & 0x0F
		} else {
			v = (b[i/2] >> 4) & 0x0F
		}
		c := uint(v) + 1 // 1..16 のはずだが 1..10 を想定
		if c < 1 || c > 10 {
			return nil, errors.New("decoded count out of range (1..10)")
		}
		out[i] = c
	}
	return out, nil
}
