package deckcode

import "errors"

type bitWriter struct {
	buf   []byte
	bits  uint64 // バッファ（下位側から詰める）
	nbits int    // buf内の有効ビット数
}

func (w *bitWriter) write(val uint32, width int) {
	// 下位widthビットを追加（LSB-firstで詰め、最後にバイトに落とす）
	w.bits |= uint64(val&((1<<width)-1)) << w.nbits
	w.nbits += width
	for w.nbits >= 8 {
		w.buf = append(w.buf, byte(w.bits&0xFF))
		w.bits >>= 8
		w.nbits -= 8
	}
}

func (w *bitWriter) finish() []byte {
	if w.nbits > 0 {
		w.buf = append(w.buf, byte(w.bits&0xFF))
		w.bits = 0
		w.nbits = 0
	}
	return w.buf
}

type bitReader struct {
	src   []byte
	acc   uint64
	nbits int
	cur   int
}

func (r *bitReader) read(width int) (uint32, error) {
	for r.nbits < width {
		if r.cur >= len(r.src) {
			return 0, errors.New("bitReader: eof")
		}
		r.acc |= uint64(r.src[r.cur]) << r.nbits
		r.cur++
		r.nbits += 8
	}
	mask := uint64((1 << width) - 1)
	out := uint32(r.acc & mask)
	r.acc >>= width
	r.nbits -= width
	return out, nil
}
