package matrix

import (
	"errors"
	"fmt"
)

var ErrorInvalidData = errors.New("Provided Category/Location/Type/Price is greater then maximum allowed by implementation")

// Query represents query for matrix to execute. Category takes 8 bits and Locations takes 9 bits.
type Query struct {
	Category uint32
	Location uint32
}

func NewQuery(category, location uint32) (Query, error) {
	q := Query{category, location}
	return q, querySanityCheck(q)
}

func querySanityCheck(q Query) error {
	if q.Category > 1<<8 || q.Location > 1<<9 {
		return ErrorInvalidData
	}
	return nil
}

func (q Query) String() string {
	return fmt.Sprintf("Query{Location: %d; Category: %d}", q.Location, q.Category)
}

type Result struct {
	Type  byte
	Price uint16
}

func NewResult(t byte, p uint16) (Result, error) {
	r := Result{t, p}
	return r, resultSanityCheck(r)
}

func resultSanityCheck(r Result) error {
	if r.Type > 1<<2 || r.Price > 1<<13 {
		return ErrorInvalidData
	}
	return nil
}

func UnpackResult(u uint32) Result {
	price := uint16((u & resultBitMask) >> 2)
	t := byte(u & (1<<2 - 1))
	return Result{t, price}
}

func (r Result) String() string {
	return fmt.Sprintf("Result{Type: %d; Price: %d}", r.Type, r.Price)
}

func (r Result) Pack() uint16 {
	return uint16(r.Price<<2) ^ uint16(r.Type)
}
