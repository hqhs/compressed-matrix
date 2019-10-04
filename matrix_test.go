package matrix

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestEqChecker(t *testing.T) {

}

func TestQueryConcat(t *testing.T) {

}

func TestCompleteMatrix(t *testing.T) {
	var N int
	if testing.Short() {
		N = 5000
	} else {
		N = categoriesAmount * locationsAmount
	}
	// TODO
	data := generateRandomData()
	quer, err := New(data)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	gen := newGen(data)
	subtest := func(tt *testing.T) {
		q, r := <-gen.q, <-gen.r

		res, ok := quer.Do(q)
		if !ok {
			tt.Errorf("result for provided query not found")
		} else if res.Type != r.Type || res.Price != res.Price {
			// tt.Errorf("locs: %v\nmcats: %v\n", quer.locations, quer.mcategories)
			tt.Errorf(
				"%s\nsingle query/result storage does not work:\n Original Result{%5d} - %16b\n   Stored Result{%5d} - %16b",
				// quer,
				q,
				r.Type, r.Price,
				res.Type, res.Price,
			)
		}
	}
	for i := 0; i < N; i++ {
		t.Run(fmt.Sprintf("test %d", i), subtest)
	}
}

func generateRandomData() map[Query]Result {
	locations := randUniqSlice(categoriesAmount)
	microcats := randUniqSlice(locationsAmount)

	data := map[Query]Result{}
	for _, l := range locations {
		for _, m := range microcats {
			q := Query{l, m}
			data[q] = genResult()
		}
	}
	return data
}

func randUniqSlice(n int) []uint32 {
	checker := map[uint32]bool{}

	slice := make([]uint32, n)

	var u, counter uint32
	for {
		if counter == uint32(n) {
			break
		}
		u = rand.Uint32()
		if ok := checker[u]; !ok {
			slice[counter] = u
			checker[u] = true
			counter++
		}
	}
	return slice
}

type genQR struct {
	q chan Query
	r chan Result
}

func newGen(data map[Query]Result) *genQR {
	g := &genQR{q: make(chan Query, 1), r: make(chan Result, 1)}
	go func() {
		for q, r := range data {
			g.q <- q
			g.r <- r
		}
	}()
	return g
}
