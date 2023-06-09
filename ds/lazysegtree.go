// interface for lazy seg tree
type lazySegTreeManager[T any, U any, L any] interface {
	neutralNode() T
	neutralLazy() L
	mergeNodes(node1, node2 T) T
	updateNode(node *T, lazy *L, update U, l, r int32)
	pushLazy(parent L, node1, node2 *T, lazy1, lazy2 *L, l, r int32)
}

// lazy segtree data structure
type lazySegTree[
	T any,
	U any,
	L any,
	M lazySegTreeManager[T, U, L],
] struct {
	nodes      []T
	lazy       []L
	first, len int // arbitrary indexing
	manager    M
}

// internal. pushes lazy tag from v to its children
func (t *lazySegTree[T, U, L, M]) _push(v, l, r int32) {
	if l != r {
		t.manager.pushLazy(t.lazy[v], &t.nodes[v<<1], &t.nodes[v<<1|1],
			&t.lazy[v<<1], &t.lazy[v<<1|1], l, r)
	}
	t.lazy[v] = t.manager.neutralLazy()
}

// internal. builds the segtree
func (t *lazySegTree[T, U, L, M]) _build(v, l, r int32, data []U) {
	if l == r {
		t.nodes[v] = t.manager.neutralNode()
		t.lazy[v] = t.manager.neutralLazy()
		if l < int32(len(data)) {
			t.manager.updateNode(&t.nodes[v], &t.lazy[v], data[l], l, r)
		}
		return
	}

	mid := (l + r) >> 1

	t._build(v<<1, l, mid, data)
	t._build(v<<1|1, mid+1, r, data)

	t.nodes[v] = t.manager.mergeNodes(t.nodes[v<<1], t.nodes[v<<1|1])
	t.lazy[v] = t.manager.neutralLazy()
}

// internal. updates the segtree
func (t *lazySegTree[T, U, L, M]) _update(v, l, r, l0, r0 int32, u U) {
	t._push(v, l, r)
	if l == l0 && r == r0 {
		t.manager.updateNode(&t.nodes[v], &t.lazy[v], u, l, r)
		return
	}

	mid := (l + r) >> 1

	if r0 <= mid {
		t._update(v<<1, l, mid, l0, r0, u)
		t._push(v<<1|1, mid+1, r)
	} else if mid < l0 {
		t._push(v<<1, l, mid)
		t._update(v<<1|1, mid+1, r, l0, r0, u)
	} else {
		t._update(v<<1, l, mid, l0, mid, u)
		t._update(v<<1|1, mid+1, r, mid+1, r0, u)
	}

	t.nodes[v] = t.manager.mergeNodes(t.nodes[v<<1], t.nodes[v<<1|1])
}

// internal. queries the segtree
func (t *lazySegTree[T, U, L, M]) _query(v, l, r, l0, r0 int32) T {
	t._push(v, l, r)
	if l == l0 && r == r0 {
		return t.nodes[v]
	}

	mid := (l + r) >> 1

	if r0 <= mid {
		return t._query(v<<1, l, mid, l0, r0)
	} else if l0 > mid {
		return t._query(v<<1|1, mid+1, r, l0, r0)
	}
	return t.manager.mergeNodes(
		t._query(v<<1, l, mid, l0, mid),
		t._query(v<<1|1, mid+1, r, mid+1, r0),
	)
}

// internal. finds the min pos where f() of range [l0, pos] is true, or returns -1
func (t *lazySegTree[T, U, L, M]) _findMinFromLeft(v, l, r, l0, r0 int32, f func(T) bool, left T) (
	int32, T) {
	if r < l0 || r0 < l {
		return -1, t.manager.neutralNode()
	}

	t._push(v, l, r)

	if l == r {
		if f(t.manager.mergeNodes(left, t.nodes[v])) {
			return l, t.nodes[v]
		}
		return -1, t.nodes[v]
	}

	mid := (l + r) >> 1

	if l0 <= l && r <= r0 {
		if x := t.manager.mergeNodes(left, t.nodes[v<<1]); f(x) {
			return t._findMinFromLeft(v<<1, l, mid, l0, r0, f, left)
		} else if f(t.manager.mergeNodes(x, t.nodes[v<<1|1])) {
			return t._findMinFromLeft(v<<1|1, mid+1, r, l0, r0, f, x)
		}
		return -1, t.nodes[v]
	}

	res, seg := t._findMinFromLeft(v<<1, l, mid, l0, r0, f, left)
	if res != -1 {
		return res, seg
	}

	return t._findMinFromLeft(v<<1|1, mid+1, r, l0, r0, f, t.manager.mergeNodes(left, seg))
}

// internal. finds the max pos where f() of range [pos, r0] is true, or returns -1
func (t *lazySegTree[T, U, L, M]) _findMaxFromRight(v, l, r, l0, r0 int32, f func(T) bool, right T) (
	int32, T) {
	if r < l0 || r0 < l {
		return -1, t.manager.neutralNode()
	}

	t._push(v, l, r)

	if l == r {
		if f(t.manager.mergeNodes(right, t.nodes[v])) {
			return l, t.nodes[v]
		}
		return -1, t.nodes[v]
	}

	mid := (l + r) >> 1

	if l0 <= l && r <= r0 {
		if x := t.manager.mergeNodes(right, t.nodes[v<<1|1]); f(x) {
			return t._findMaxFromRight(v<<1|1, mid+1, r, l0, r0, f, right)
		} else if f(t.manager.mergeNodes(x, t.nodes[v<<1])) {
			return t._findMaxFromRight(v<<1, l, mid, l0, r0, f, x)
		}
		return -1, t.nodes[v]
	}

	res, seg := t._findMaxFromRight(v<<1|1, mid+1, r, l0, r0, f, right)
	if res != -1 {
		return res, seg
	}

	return t._findMaxFromRight(v<<1, l, mid, l0, r0, f, t.manager.mergeNodes(right, seg))
}

// inits the segtree on range [l, r]
func (t *lazySegTree[T, U, L, M]) init(l, r int, data []U) {
	t.first = l
	t.len = r - l + 1

	t.nodes = make([]T, (t.len<<2)+5)
	t.lazy = make([]L, (t.len<<2)+5)

	t._build(1, 0, int32(t.len-1), data)
}

// updates segtree by applying update on range [l, r]
func (t *lazySegTree[T, U, L, M]) update(l, r int, update U) {
	l -= t.first
	r -= t.first
	l = max(l, 0)
	r = min(r, t.len-1)
	if l > r {
		panic("l > r in segtree update")
	}

	t._update(1, 0, int32(t.len-1), int32(l), int32(r), update)
}

// queries segtree on range [l, r]
func (t *lazySegTree[T, U, L, M]) query(l, r int) T {
	l -= t.first
	r -= t.first
	l = max(l, 0)
	r = min(r, t.len-1)
	if l > r {
		panic("l > r in segtree query")
		// return t.manager.neutral()
	}
	return t._query(1, 0, int32(t.len-1), int32(l), int32(r))
}

// finds the min pos where f() of range [l, pos] is true, or returns r+1
func (t *lazySegTree[T, U, L, M]) findMinFromLeft(l, r int, f func(T) bool) int {
	l -= t.first
	r -= t.first
	l = max(l, 0)
	r = min(r, t.len-1)
	if l > r {
		panic("l > r in segtree findFirst")
	}

	pos, _ := t._findMinFromLeft(1, 0, int32(t.len-1), int32(l), int32(r), f, t.manager.neutralNode())
	if pos == -1 {
		return t.first + r + 1
	}
	return t.first + int(pos)
}

// finds the max pos where f() of range [l, pos] is true, or returns l-1
func (t *lazySegTree[T, U, L, M]) findMaxFromLeft(l, r int, f func(T) bool) int {
	pos := t.findMinFromLeft(l, r, func(node T) bool {
		return !f(node)
	})
	return pos - 1
}

// finds the max pos where f() of range [pos, r] is true, or returns r+1
func (t *lazySegTree[T, U, L, M]) findMaxFromRight(l, r int, f func(T) bool) int {
	l -= t.first
	r -= t.first
	l = max(l, 0)
	r = min(r, t.len-1)
	if l > r {
		panic("l > r in segtree findFirst")
	}

	pos, _ := t._findMaxFromRight(1, 0, int32(t.len-1), int32(l), int32(r), f, t.manager.neutralNode())
	if pos == -1 {
		return t.first - 1
	}
	return t.first + int(pos)
}

// finds the min pos where f() of range [pos, r] is true, or returns l-1
func (t *lazySegTree[T, U, L, M]) findMinFromRight(l, r int, f func(T) bool) int {
	pos := t.findMaxFromRight(l, r, func(node T) bool {
		return !f(node)
	})
	return pos + 1
}

// returns a slice with all the data of the tree
func (t *lazySegTree[T, U, L, M]) slice() []T {
	var res []T
	for i := 0; i < t.len; i++ {
		res = append(res, t.query(i, i))
	}	
	return res
}

// returns a string representation of the tree
func (t *lazySegTree[T, U, L, M]) str() string {
	return fmt.Sprint(t.slice())
}
