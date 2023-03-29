package hash

import "hash/maphash"

type Hash struct {
	h *maphash.Hash
}

func New() *Hash {
	h := &Hash{h: &maphash.Hash{}}
	h.h.SetSeed(maphash.MakeSeed())
	return h
}

func (h *Hash) Sum64(b []byte) uint64 {
	h.h.Reset()
	_, err := h.h.Write(b)
	if err != nil {
		return 0
	}
	return h.h.Sum64()
}

func (h *Hash) Sum64String(s string) uint64 {
	return h.Sum64([]byte(s))
}
