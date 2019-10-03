package matrix

import (
	"math/rand"
	"testing"
)

func genResult() Result {
	return UnpackResult(rand.Uint32() & resultBitMask)
}

func TestResultPackUnpack(t *testing.T) {
	for i := 0; i < resultBitMask; i++ {
		r := genResult()
		packed := r.Pack()
		unpacked := UnpackResult(uint32(packed))
		if r.Type != unpacked.Type || r.Price != unpacked.Price {
			t.Errorf("Result packing/unpacking does not work properly")
		}
	}
}
