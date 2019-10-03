package matrix

import "testing"

func TestPackUnpack(t *testing.T) {
	tA, tB, tC := uint32((1<<7 + 1)), uint32((1<<8 + 1)), uint32((1<<13 + 1))

	packed := pack(tA, tB, tC)
	unA, unB, unC := unpack(packed)
	if tA != unA || tB != unB || tC != unC {
		t.Errorf("pack/unpack does not work propertly: pack(%d,%d,%d) = %b", tA, tB, tC, packed)
	}
}

// 8 second test on mac pro 2018
func TestCompletePackUnpack(t *testing.T) {

	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}

	maxCategory := categoryBitMask >> 24
	maxLocation := locationBitMask >> 15
	maxResult := resultBitMask

	for c := 0; c < maxCategory; c++ {
		for l := 0; l < maxLocation; l++ {
			for r := 0; r < maxResult; r++ {
				packed := pack(uint32(c), uint32(l), uint32(r))
				unpackedC, unpackedL, unpackedR := unpack(packed)
				if uint32(c) != unpackedC || uint32(l) != unpackedL || uint32(r) != unpackedR {
					t.Errorf("pack/unpack does not work properly")
				}
			}
		}
	}
}
