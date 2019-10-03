package matrix

const (
	categoryBitMask = 0xFF000000 // 11111111000000000000000000000000
	locationBitMask = 0xFF8000   // 00000000111111111000000000000000

	resultBitMask = 0x7FFF // 00000000000000000111111111111111

	tail8bits = (1<<9 - 1)
	tail9bits = (1<<10 - 1)
)

func pack(cat, loc, res uint32) uint32 {
	return uint32(((cat & tail8bits) << 24) + ((loc & tail9bits) << 15) + (res & resultBitMask)) // FIXME XOR would be faster
}

func unpack(u uint32) (cat, loc, res uint32) {
	cat = (u & categoryBitMask) >> 24
	loc = (u & locationBitMask) >> 15
	res = u & resultBitMask
	return
}
