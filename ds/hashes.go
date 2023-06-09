// manager interface for the hashSequence
type hashingManager[T any] interface {
	hash(x T, base, mod int) int
}

// a data structure represetiong a hash sequence
type hashSequence[T any, M hashingManager[T]] struct {
	manager    M
	len        int
	base       int
	mod        int
	basePowers []int
	hashes     []int
}

// inits the hash sequence with base, mod and data being added
func (h *hashSequence[T, M]) init(n, base, mod int, data ...T) {
	h.hashes = make([]int, 0, n)
	h.len = n

	h.base = base
	h.mod = mod

	h.basePowers = make([]int, n)
	h.basePowers[0] = 1
	for i := 1; i < n; i++ {
		h.basePowers[i] = h.basePowers[i-1] * base % h.mod
	}

	if len(data) > n {
		panic("data is too large")
	}

	for _, x := range data {
		h.add(x)
	}

}

// adds x to the end of the sequence
func (h *hashSequence[T, M]) add(x T) {
	if len(h.hashes) == h.len {
		panic("can't add more then len elements")
	}

	h.hashes = append(h.hashes, h.manager.hash(x, h.base, h.mod)*
		h.basePowers[h.len-len(h.hashes)-1]%h.mod)

	if len(h.hashes) >= 2 {
		h.hashes[len(h.hashes)-1] += h.hashes[len(h.hashes)-2]
		h.hashes[len(h.hashes)-1] %= h.mod
	}
}

// returns hash of the sequence [l, r]
func (h hashSequence[T, M]) query(l, r int) int {
	if l > r || l < 0 || r >= h.len {
		panic("hash query out of bounds")
	}
	res := h.hashes[r]
	if l != 0 {
		res -= h.hashes[l-1]
		res *= h.basePowers[l]
		res %= h.mod
		if res < 0 {
			res += h.mod
		}
	}
	return res
}

// a data structure represetiong a double hash sequence
type doubleHashSequence[T any, M hashingManager[T]] struct {
	a, b hashSequence[T, M]
}

// inits the hash sequence with base, mod and data being added
func (h *doubleHashSequence[T, M]) init(n, base1, mod1, base2, mod2 int, data ...T) {
	h.a.init(n, base1, mod1, data...)
	h.b.init(n, base2, mod2, data...)
}

// adds x to the end of the sequence
func (h *doubleHashSequence[T, M]) add(x T) {
	h.a.add(x)
	h.b.add(x)
}

// returns hash of the sequence [l, r]
func (h *doubleHashSequence[T, M]) query(l, r int) pair[int] {
	return pair[int]{h.a.query(l, r), h.b.query(l, r)}
}
