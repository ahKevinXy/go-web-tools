package hash

import "testing"

func TestString(t *testing.T) {
	h := New()
	a := h.Sum64String("ahKevinXy")
	b := h.Sum64String("ahKevinXy")
	t.Log(a)
	t.Log(b)
}

func TestInt64(t *testing.T) {
	h := New()
	a := h.Sum64([]byte("ahKevinXy"))
	b := h.Sum64([]byte("ahKevinXy"))
	t.Log(a)
	t.Log(b)
}
