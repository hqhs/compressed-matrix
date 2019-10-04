package matrix

import (
	"errors"
	"sort"
)

const (
	categoriesAmount = 200
	locationsAmount  = 400
)

var (
	ErrorInvalidConstant = errors.New("Amount of categories/locations is fixed and limited by implementation")
	ErrorQueryNotFound   = errors.New("Failed to found provided Query in compressed data")
)

type Matrix interface {
	Do(Query) (Result, bool)
	Bytes() []byte
	Map() map[Query]Result
}

type matrix struct {
	compressed []uint32
	locations  []uint32
	categories []uint32
}

func New(data map[Query]Result) (Matrix, error) {
	if categoriesAmount > 1<<8 || locationsAmount > 1<<9 {
		return nil, ErrorInvalidConstant
	}
	m := &matrix{
		compressed: make([]uint32, 0, categoriesAmount*locationsAmount),
		locations:  make([]uint32, 0, locationsAmount),
		categories: make([]uint32, 0, categoriesAmount),
	}

	uniqChecker := make(map[uint32]bool)
	for q := range data {
		if ok := uniqChecker[q.Location]; !ok {
			uniqChecker[q.Location] = true
			m.locations = append(m.locations, q.Location)
		}
	}

	uniqChecker = make(map[uint32]bool)
	for q := range data {
		if ok := uniqChecker[q.Category]; !ok {
			uniqChecker[q.Category] = true
			m.categories = append(m.categories, q.Category)
		}
	}

	sortU(m.locations)
	sortU(m.categories)

	var u uint32
	for q, r := range data {
		locN, _ := findU(q.Location, m.locations)
		catN, _ := findU(q.Category, m.categories)

		packR := r.Pack()

		u = pack(catN, locN, uint32(packR))
		m.compressed = append(m.compressed, u)
	}

	// NOTE check compr uniq maybe?..

	sortU(m.compressed)

	return m, nil
}

func NewFromBytes(b []byte) (Matrix, error) {
	// first two numbers are length of categories and locations, then locations and categories, then compressed data
	return nil, nil
}

func (m *matrix) Bytes() []byte {
	return nil
}

func (m *matrix) Do(q Query) (Result, bool) {
	return m.do(q)
}

func (m *matrix) do(q Query) (Result, bool) {
	catN, ok := findU(q.Category, m.categories)
	if !ok {
		return Result{}, false
	}
	locN, ok := findU(q.Location, m.locations)
	if !ok {
		return Result{}, false
	}
	N := pack(catN, locN, 0)

	for _, entry := range m.compressed {
		if cmpHead17bits(entry, N) {
			_, res := m.unpackQR(entry)
			return res, true
		}
	}

	return Result{}, false
}

func (m *matrix) unpackQR(u uint32) (Query, Result) {
	catN, locN, res := unpack(u)
	return Query{Location: m.locations[locN-1], Category: m.categories[catN-1]}, UnpackResult(res)
}

func cmpHead17bits(a uint32, b uint32) bool {
	mask := uint32(categoryBitMask | locationBitMask)
	return (mask & a) == (mask & b)
}

func (m *matrix) Map() map[Query]Result {
	r := make(map[Query]Result)
	for _, entry := range m.compressed {
		q, res := m.unpackQR(entry)
		r[q] = res
	}
	return r
}

func sortU(u []uint32) {
	sort.Slice(u, func(i, j int) bool {
		return u[i] < u[j]
	})
}

func findU(u uint32, us []uint32) (uint32, bool) {
	// TODO use binary search
	for i := 0; i < len(us); i++ {
		if us[i] == u {
			return uint32(i + 1), true
		}
	}
	return uint32(0), false
}
